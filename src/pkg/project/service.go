package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/config"
	"github.com/kplcloud/kplcloud/src/jenkins"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/discovery"
	"github.com/kplcloud/kplcloud/src/pkg/hooks"
	"github.com/kplcloud/kplcloud/src/redis"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/convert"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"github.com/kplcloud/kplcloud/src/util/helper"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	"github.com/kplcloud/kplcloud/src/util/pods"
	"gopkg.in/guregu/null.v3"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"sync"
)

var (
	ErrProjectExist          = errors.New("项目名称已存在，请重新填写")
	ErrProjectNotExist       = errors.New("项目不存在")
	ErrProjectName           = errors.New("项目英文名不合法")
	ErrProjectAuditSubmit    = errors.New("项目已提交审核，不允许修改")
	ErrProjectAuditPass      = errors.New("项目已审核通过，不允许修改")
	ErrProjectList           = errors.New("项目列表获取错误")
	ErrTemplateGet           = errors.New("模板获取错误")
	ErrTemplateEncodeGet     = errors.New("模板生成错误")
	ErrProjectTemplateGet    = errors.New("项目生成的模版获取错误")
	ErrProjectTemplateUpdate = errors.New("项目生成的模版更新错误")
	ErrProjectLanguageJava   = errors.New("该项目不是Java语言不需要调整pom.xml文件路径")
	ErrProjectGet            = errors.New("项目读取错误,请重试")
	ErrProjectDeleteName     = errors.New("项目删除失败，项目名不正确")
	ErrIsInGroupFailed       = errors.New("判断是否在组里失败")
	ErrGroupNoPermission     = errors.New("没有相关组的权限")
	ErrPodDeploymentGet      = errors.New("项目获取错误")
	ErrPodDeploymentPodList  = errors.New("项目pods获取错误")
	ErrParamsMetrics         = errors.New("metrics不能为空")
)

type Service interface {
	// 创建项目step1
	Post(ctx context.Context, ns, name, displayName, desc string) (err error)

	// 创建项目step2
	BasicPost(ctx context.Context, name string, req basicRequest) (err error)

	// 项目列表 todo: 可以考虑定时同步的方案...... 能省很多事
	List(ctx context.Context, page, limit int, name string, groupId int64) (res map[string]interface{}, err error)

	// 根据业务线名称获取业务应用列表
	ListByNs(ctx context.Context) (res []map[string]interface{}, err error)

	// 修改pomfile文件地址
	PomFile(ctx context.Context, pomFile string) error

	// 修改git路径
	GitAddr(ctx context.Context, gitAddr string) error

	// 项目详情页信息
	Detail(ctx context.Context) (res map[string]interface{}, err error)

	// 更新项目信息
	Update(ctx context.Context, displayName, desc string) (err error)

	// Workspace
	Workspace(ctx context.Context) (res []map[string]interface{}, err error) //工作台展示
	//Metrics(ctx context.Context) (res map[string]interface{}, err error) // 这个挪到pod模块里了

	// 同步项目信息
	Sync(ctx context.Context) (err error)

	// 删除项目
	Delete(ctx context.Context, ns, name, code string) (err error)

	// 获取配置信息
	Config(ctx context.Context) (res map[string]interface{}, err error)

	// 获取项目监控指标 /project/{ns}/monitor/{projectName}?podName=xxxxx&metrics=memory/request&container
	Monitor(ctx context.Context, metrics, podName, container string) (res map[string]map[string]map[string][]pods.XYRes, err error)

	// 告警统计
	Alerts(ctx context.Context) (res alertsResponse, err error)
}

type service struct {
	logger         log.Logger
	config         *config.Config
	redisInterface redis.RedisInterface
	k8sClient      kubernetes.K8sClient
	amqp           amqp.AmqpClient
	jenkins        jenkins.Jenkins
	repository     repository.Repository
	hookQueueSvc   hooks.ServiceHookQueue
	//istioClient               istio.IstioClient
}

func NewService(logger log.Logger, config *config.Config,
	redisInterface redis.RedisInterface,
	k8sClient kubernetes.K8sClient,
	amqp amqp.AmqpClient,
	jenkins jenkins.Jenkins,
	repository repository.Repository,
	hookQueueSvc hooks.ServiceHookQueue) Service {
	return &service{
		logger,
		config,
		redisInterface,
		k8sClient,
		amqp,
		jenkins,
		repository,
		hookQueueSvc}
}

func (c *service) Alerts(ctx context.Context) (res alertsResponse, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		alertTotal, err := c.repository.Notice().CountByAction(ns, project.Name, types.NoticeActionAlarm)
		if err != nil {
			_ = level.Warn(c.logger).Log("Notice", "CountByAction", "err", err.Error())
		}
		res.AlertTotal = alertTotal
		wg.Done()
	}()

	go func() {
		buildTotal, err := c.repository.Build().CountByStatus(ns, project.Name, "")
		if err != nil {
			_ = level.Warn(c.logger).Log("Build", "CountByStatus", "err", err.Error())
		}
		res.BuildTotal = buildTotal
		wg.Done()
	}()

	go func() {
		rollbackNum, err := c.repository.Build().CountByStatus(ns, project.Name, repository.BuildRoolback)
		if err != nil {
			_ = level.Warn(c.logger).Log("Build", "CountByStatus", "err", err.Error())
		}
		res.RollbackTotal = rollbackNum
		wg.Done()
	}()
	go func() {
		buildFailure, err := c.repository.Build().CountByStatus(ns, project.Name, repository.BuildFailure)
		if err != nil {
			_ = level.Warn(c.logger).Log("Build", "CountByStatus", "err", err.Error())
		}
		res.BuildFailureTotal = buildFailure
		wg.Done()
	}()

	go func() {
		buildSuccess, err := c.repository.Build().CountByStatus(ns, project.Name, repository.BuildSuccess)
		if err != nil {
			_ = level.Warn(c.logger).Log("Build", "CountByStatus", "err", err.Error())
		}
		res.BuildSuccessTotal = buildSuccess
		wg.Done()
	}()

	wg.Wait()

	return
}

func (c *service) Monitor(ctx context.Context, metrics, podName, container string) (res map[string]map[string]map[string][]pods.XYRes, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	dep, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(project.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return res, ErrPodDeploymentGet
	}

	var selectorKey, selectorVal string
	for key, val := range dep.Spec.Selector.MatchLabels {
		selectorKey = key
		selectorVal = val
	}

	podList, err := c.k8sClient.Do().CoreV1().Pods(ns).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", selectorKey, selectorVal),
	})

	if err != nil {
		_ = level.Error(c.logger).Log("Pods", "List", "err", err.Error())
		return res, ErrPodDeploymentPodList
	}

	resp := map[string]map[string]map[string][]pods.XYRes{}
	for _, pod := range podList.Items {
		containersList := map[string]map[string][]pods.XYRes{}
		for _, v := range pod.Spec.Containers {

			metricsData := pods.GetPodContainerMetrics(ns, pod.Name, c.config.GetString("server", "heapster_url"), v.Name, []string{
				"memory/usage",
				"cpu/usage",
				"network/tx_rate", // 每秒通过网络发送的字节数。
				"network/rx_rate", // 每秒通过网络接收的字节数。
			})

			containersList[v.Name] = metricsData.Metrics[v.Name]
		}
		resp[pod.Name] = containersList
	}

	// /api/v1/model/namespaces/kube-public/pods/kplcloud-6f5987d5df-lhsh5/metrics/ 容器指示
	// /api/v1/model/namespaces/kube-public/pods/kplcloud-6f5987d5df-lhsh5/containers/kplcloud/metrics 容器指示

	// 各个容器的内存指标
	// 各个容器的网络指标
	// 各个容器的CPU指标

	return resp, nil
}

func (c *service) Sync(ctx context.Context) (err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	userId := ctx.Value(middleware.UserIdContext).(int64)

	list, err := c.k8sClient.Do().AppsV1().Deployments(ns).List(metav1.ListOptions{})
	if err != nil {
		return
	}

	for _, v := range list.Items {
		if project, err := c.repository.Project().FindByNsName(ns, v.Name); err == nil && project != nil {
			_ = level.Error(c.logger).Log("projectRepository", "FindByNsName", "err", "project is exists.")
			continue
		}
		project := &types.Project{
			AuditState:   int64(repository.AuditPass),
			Desc:         v.ResourceVersion,
			Language:     repository.Golang.String(),
			MemberID:     userId,
			DisplayName:  v.Name,
			Namespace:    v.Namespace,
			Name:         v.Name,
			PublishState: 1,
			Step:         2,
		}
		if err = c.repository.Project().Create(project); err != nil {
			_ = level.Error(c.logger).Log("projectRepository", "Create", "err", err.Error())
			continue
		}
		//tpl, err := c.templateRepository.FindByKindType(repository.DeploymentKind)
		//if err != nil {
		//	_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
		//	continue
		//}
		b, _ := yaml.Marshal(v)
		var ports []types.Port
		for _, port := range v.Spec.Template.Spec.Containers[0].Ports {
			ports = append(ports, types.Port{
				Name:     port.Name,
				Port:     int(port.ContainerPort),
				Protocol: string(port.Protocol),
			})
		}
		field := types.TemplateField{
			Args:      v.Spec.Template.Spec.Containers[0].Args,
			Command:   v.Spec.Template.Spec.Containers[0].Command,
			GitAddr:   "git@github.com/kubernetes/" + v.Name + ".git",
			GitType:   "branch",
			Image:     v.Spec.Template.Spec.Containers[0].Image,
			Branch:    "master",
			Language:  "Golang",
			Name:      v.Name,
			Namespace: v.Namespace,
			Ports:     ports,
		}
		bb, _ := json.Marshal(field)
		_, err = c.repository.ProjectTemplate().FirstOrCreate(project.ID, repository.Deployment, string(bb), string(b), 1)
		if err != nil {
			_ = level.Error(c.logger).Log("projectTemplateRepository", "FirstOrCreate", "err", err.Error())
			continue
		}

		if svc, err := c.k8sClient.Do().CoreV1().Services(ns).Get(v.Name, metav1.GetOptions{}); err == nil {
			{
				var ports []types.Port
				for _, port := range svc.Spec.Ports {
					ports = append(ports, types.Port{
						Name:     port.Name,
						Port:     int(port.Port),
						Protocol: string(port.Protocol),
					})
				}
				field := types.ServiceField{
					Name:      svc.Name,
					Namespace: svc.Namespace,
					Ports:     ports,
				}
				b, _ := yaml.Marshal(svc)
				bb, _ := json.Marshal(field)
				_, err = c.repository.ProjectTemplate().FirstOrCreate(project.ID, repository.Service, string(bb), string(b), 1)
				if err != nil {
					_ = level.Error(c.logger).Log("projectTemplateRepository", "FirstOrCreate", "err", err.Error())
					continue
				}
			}
		}

		if ing, err := c.k8sClient.Do().ExtensionsV1beta1().Ingresses(ns).Get(v.Name, metav1.GetOptions{}); err == nil {
			{
				var rules []*types.RuleStruct
				var paths []*types.Paths
				for _, rule := range ing.Spec.Rules {
					for _, path := range rule.HTTP.Paths {
						paths = append(paths, &types.Paths{
							Path:        path.Path,
							ServiceName: path.Backend.ServiceName,
							PortName:    int(path.Backend.ServicePort.IntVal),
						})
					}
					rules = append(rules, &types.RuleStruct{
						Domain: rule.Host,
						Paths:  paths,
					})
				}
				field := types.IngressField{
					Namespace: ing.Namespace,
					Name:      ing.Name,
					Rules:     rules,
				}
				b, _ := yaml.Marshal(ing)
				bb, _ := json.Marshal(field)
				_, err = c.repository.ProjectTemplate().FirstOrCreate(project.ID, repository.Ingress, string(bb), string(b), 1)
				if err != nil {
					_ = level.Error(c.logger).Log("projectTemplateRepository", "FirstOrCreate", "err", err.Error())
					continue
				}
			}
		}

		if configmap, err := c.k8sClient.Do().CoreV1().ConfigMaps(ns).Get(v.Name, metav1.GetOptions{}); err == nil {
			{
				b, _ := yaml.Marshal(configmap)
				_, err = c.repository.ProjectTemplate().FirstOrCreate(project.ID, repository.ConfigMap, "", string(b), 1)
				if err != nil {
					_ = level.Error(c.logger).Log("projectTemplateRepository", "FirstOrCreate", "err", err.Error())
					continue
				}
			}
		}
	}

	return
}

func (c *service) Update(ctx context.Context, displayName, desc string) (err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	project.DisplayName = displayName
	project.Desc = desc
	err = c.repository.Project().Update(project)
	return
}

func (c *service) Post(ctx context.Context, ns, name, displayName, desc string) (err error) {
	if !convert.IsEnNameString(name) {
		_ = level.Error(c.logger).Log("Project", "Post", "NameRefused", name)
		return ErrProjectName
	}

	memberId := ctx.Value(middleware.UserIdContext).(int64)
	_, notExist := c.repository.Project().FindByNsNameExist(ns, name)
	if notExist != true {
		_ = level.Error(c.logger).Log("Project", "Post", "FindByNsNameExist", notExist)
		return ErrProjectExist
	}

	return c.repository.Project().Create(&types.Project{
		Namespace:   ns,
		Name:        name,
		DisplayName: displayName,
		Desc:        desc,
		MemberID:    memberId,
		Language:    repository.Golang.String(), //默认Golang,第二步可自行选择修改
	})
}

func (c *service) BasicPost(ctx context.Context, name string, req basicRequest) (err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	project, notExist := c.repository.Project().FindByNsNameExist(ns, name)
	if notExist == true {
		return ErrProjectNotExist
	}
	//check AuditState.
	if project.AuditState == int64(repository.AuditSubmit) {
		return ErrProjectAuditSubmit
	}
	if project.AuditState == int64(repository.AuditPass) {
		return ErrProjectAuditPass
	}

	req.GitAddr = helper.GitUrl(c.config.GetString("git", "git_addr")) + req.GitAddr

	// check Service
	var ports []map[string]interface{}
	if len(req.Ports) > 0 {
		for _, v := range req.Ports {
			if res, err := convert.Struct2Json2Map(v); err == nil {
				ports = append(ports, res)
			}
		}
		data := map[string]interface{}{
			"ports":        ports,
			"name":         project.Name,
			"namespace":    project.Namespace,
			"resourceType": discovery.ServiceResourceType,
		}
		if err = c.rewriteTemplate(project.ID, repository.ServiceKind, data); err != nil {
			return err
		}
	}

	// Deployment
	deployment, err := convert.Struct2Json2Map(req)
	if err != nil {
		return err
	}
	resourceInfo := pods.CreateCpuData(req.Resources)
	deployment["memory"] = resourceInfo.Memory
	deployment["maxMemory"] = resourceInfo.MaxMemory
	//@todo 记得修改模板中的Args，Command为[]string格式

	if err = c.rewriteTemplate(project.ID, repository.DeploymentKind, deployment); err != nil {
		_ = level.Error(c.logger).Log("BasicPost", "Rewrite Deployment", "err", err.Error())
		return err
	}

	// create projectJenkins
	if err = c.saveJenkins(project, req); err != nil {
		_ = level.Error(c.logger).Log("BasicPost", "Save ProjectJenkins", "err", err.Error())
		return err
	}

	// update projectState
	project.AuditState = int64(repository.AuditSubmit)
	project.Step = 2
	project.Language = req.Language
	if err = c.repository.Project().UpdateProjectById(project); err != nil {
		_ = level.Error(c.logger).Log("BasicPost", "Update Project", "err", err.Error())
		return err
	}

	// 发送审核通知
	ctx = context.WithValue(ctx, middleware.ProjectContext, project)
	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.ApplyEvent,
			name, ns,
			fmt.Sprintf("项目提交审核: %v.%v", name, ns)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()
	return nil
}

func (c *service) List(ctx context.Context, page, limit int, name string, groupId int64) (res map[string]interface{}, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	offset := (page - 1) * limit

	// 如果传入组ID
	// 判断当前用户是否为超管,如果超管,直接查,按照组来查
	// 如果当前用户不是超管,判断是否有该组的权限
	// 如果没有传入组,直接查数据
	var projects []*types.Project
	var count int64
	if groupId != 0 {
		isAdmin := ctx.Value(middleware.IsAdmin).(bool)
		memberId := ctx.Value(middleware.UserIdContext).(int64)
		// 如果不是超管,看当前用户是否有该组的权限
		if !isAdmin {
			res, err := c.repository.Groups().IsInGroup(groupId, memberId)
			if err != nil {
				_ = level.Error(c.logger).Log("groupRepository", "IsInGroup", "err", err.Error())
				return nil, ErrIsInGroupFailed
			}
			if !res {
				return nil, ErrGroupNoPermission
			}
		}
		projects, count, err = c.repository.Project().GetProjectByGroupAndPNsAndPName(name, ns, groupId, offset, limit)
	} else {
		projects, count, err = c.repository.Project().GetProjectAndTemplateByNs(ns, name, offset, limit)
	}

	if err != nil {
		_ = level.Error(c.logger).Log("projectRepository", "GetProjectAndTemplateByNs", "err", err.Error())
		return nil, ErrProjectGet
	}

	p := paginator.NewPaginator(page, limit, int(count))

	projectLen := len(projects)

	deploymentData := c.getDeployment(projects, projectLen)

	var resp []map[string]interface{}
	var wg sync.WaitGroup
	wg.Add(projectLen)
	var podData []map[string]interface{}

	var i int
	for _, project := range projects {
		var imageVersion string
		var deployment *v1.Deployment
		for _, pt := range project.ProjectTemplates {
			if repository.Kind(pt.Kind) == repository.Deployment {
				var fields types.TemplateField
				if err = json.Unmarshal([]byte(pt.Fields), &fields); err == nil {
					imageVersion = fields.Image
				}
				_ = yaml.Unmarshal([]byte(pt.FinalTemplate), &deployment)
				break
			}
		}

		if repository.AuditState(project.AuditState) == repository.AuditPass {
			go c.getPods(project.Namespace, project.Name, i, &podData, &wg)
		} else {
			wg.Done()
		}

		resp = append(resp, map[string]interface{}{
			"audit_state":   project.AuditState,
			"created_at":    project.CreatedAt,
			"member_name":   project.Member.Username,
			"name":          project.Name,
			"namespace":     project.Namespace,
			"display_name":  project.DisplayName,
			"step":          project.Step,
			"image_version": imageVersion,
			"id":            project.ID,
			"deployment":    deploymentData[project.Name],
		})
		i++
	}
	wg.Wait()

	for _, val := range podData {
		resp[val["index"].(int)]["pods"] = val["pods"]
	}

	return map[string]interface{}{
		"list": resp,
		"page": p.Result(),
	}, nil
}

func (c *service) getDeployment(projects []*types.Project, projectsLen int) map[string]*v1.Deployment {

	var wg sync.WaitGroup
	wg.Add(projectsLen)
	sm := new(SafeMap)
	sm.Map = make(map[string]*v1.Deployment, projectsLen)

	data := map[string]*v1.Deployment{}
	for _, project := range projects {
		if project.AuditState == int64(repository.AuditPass) && project.PublishState == int64(repository.PublishPass) {
			go func(ns, name string) {
				if deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{}); err == nil {
					sm.writeMap(name, deployment)
					//data[name] = deployment
				}
				wg.Done()
			}(project.Namespace, project.Name)
		} else {
			//data[project.Name] = nil
			sm.writeMap(project.Name, nil)
			wg.Done()
		}
	}

	wg.Wait()

	return data
}

type SafeMap struct {
	sync.RWMutex
	Map map[string]*v1.Deployment
}

func (sm *SafeMap) writeMap(key string, value *v1.Deployment) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

func (c *service) getPods(ns, name string, index int, podsData *[]map[string]interface{}, wg *sync.WaitGroup) {

	defer wg.Done()
	podList, err := c.k8sClient.Do().CoreV1().Pods(ns).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", name),
	})
	if err != nil {
		_ = level.Error(c.logger).Log("Pods", "List", "err", err.Error())
		return
	}

	var message, lastMessage, podName string

	var list []map[string]interface{}
	for _, podInfo := range podList.Items {
		podName = podInfo.Name
		var isReady = true
		for _, v := range podInfo.Status.Conditions {
			if v.Type != corev1.PodReady && v.Status == corev1.ConditionFalse {
				isReady = false
				message = v.Message
				break
			}
		}

		if !isReady {
			for _, v := range podInfo.Status.ContainerStatuses {
				if v.Name == podInfo.Name {
					if v.State.Waiting != nil {
						message = v.State.Waiting.Message
					}
					if v.LastTerminationState.Terminated != nil {
						lastMessage = v.LastTerminationState.Terminated.Message
					}
					break
				}
			}
		}
		list = append(list, map[string]interface{}{
			"name":         podName,
			"message":      message,
			"last_message": lastMessage,
			"project_name": name,
			"namespace":    ns,
		})
	}

	*podsData = append(*podsData, map[string]interface{}{
		"index": index,
		"pods":  list,
	})
	return
}

func (c *service) ListByNs(ctx context.Context) (res []map[string]interface{}, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	list, err := c.repository.Project().GetProjectByNs(ns)
	if err != nil {
		_ = level.Error(c.logger).Log("Project", "GetProjectByNs", "err", err.Error())
		return res, ErrProjectList
	}
	var dat map[string]interface{}
	for _, v := range list {
		dat = map[string]interface{}{
			"name":         v.Name,
			"namespace":    v.Namespace,
			"display_name": v.DisplayName,
		}
		res = append(res, dat)
	}

	return
}

func (c *service) PomFile(ctx context.Context, pomFile string) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	if project.Language != repository.Java.String() {
		return ErrProjectLanguageJava
	}

	projectTemplate, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindByProjectId", "err", err.Error())
		return ErrProjectTemplateGet
	}

	projectTemplate.FieldStruct.PomFile = pomFile

	if err = c.repository.ProjectTemplate().UpdateTemplate(projectTemplate); err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "UpdateTemplate", "err", err.Error())
		return ErrProjectTemplateUpdate
	}

	return nil
}

func (c *service) GitAddr(ctx context.Context, gitAddr string) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	projectTemplate, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindByProjectId", "err", err.Error())
		return ErrProjectTemplateGet
	}

	projectTemplate.FieldStruct.GitAddr = gitAddr

	if err = c.repository.ProjectTemplate().UpdateTemplate(projectTemplate); err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "UpdateTemplate", "err", err.Error())
		return ErrProjectTemplateUpdate
	}

	return nil
}

func (c *service) Detail(ctx context.Context) (res map[string]interface{}, err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	project.Member.Password = null.StringFrom("")
	project.Member.Openid = ""
	body := map[string]interface{}{
		"project":      project,
		"kibana_url":   c.config.GetString("server", "kibana_url"),
		"transfer_url": c.config.GetString("server", "transfer_url"),
	}

	// 待审核情况下直接获取项目信息
	_ = level.Info(c.logger).Log("auditState", project.AuditState)
	if project.AuditState == int64(repository.AuditSubmit) {
		if projectTemplate, err := c.repository.ProjectTemplate().FindProjectTemplateByProjectId(project.ID); err == nil {
			body["templateProject"] = projectTemplate
		}
		return body, nil
	}
	if project.AuditState != int64(repository.AuditPass) {
		return body, nil
	}

	if repository.AuditState(project.AuditState) == repository.AuditSubmit {
		body["member"] = project.Member
		return body, nil
	}

	tplRes, err := c.repository.ProjectTemplate().FindProjectTemplateByProjectId(project.ID)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindProjectTemplateByProjectId", "err", err.Error())
		return nil, ErrProjectTemplateGet
	}

	project.ProjectTemplates = tplRes

	getPod := make(chan []map[string]interface{})
	defer close(getPod)

	go func() {
		if podsList, err := c.k8sClient.Do().CoreV1().Pods(project.Namespace).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", project.Name),
		}); err == nil {
			var podList []map[string]interface{}
			for _, pod := range podsList.Items {
				podList = append(podList, c.parsePods(&pod))
			}
			getPod <- podList
		} else {
			getPod <- nil
		}
	}()

	var wg sync.WaitGroup
	wg.Add(len(tplRes))
	for _, tpl := range tplRes {
		go func(tpl *types.ProjectTemplate) {
			body[strings.ToLower(tpl.Kind)] = c.parseTemplate(tpl, &wg)
		}(tpl)
	}
	wg.Wait()

	body["pods"] = <-getPod
	body["templateProject"] = tplRes

	// 这里好像并没有什么用
	//eps, err := c.k8sClient.Do().CoreV1().Endpoints(project.Namespace).Get(project.Name, metav1.GetOptions{})
	//if err != nil {
	//	_ = level.Error(c.logger).Log("Endpoints", "Get", "err", err.Error())
	//}
	//
	//fmt.Println(eps)

	return body, nil
}

func (c *service) Delete(ctx context.Context, ns, name, code string) (err error) {
	if code != name {
		_ = level.Error(c.logger).Log("Delete", "Check ProjectName")
		return ErrProjectDeleteName
	}
	//step1 delete deployments
	if err = c.k8sClient.Do().AppsV1().Deployments(ns).Delete(name, nil); err != nil {
		_ = level.Error(c.logger).Log("DeleteStep1", "Delete Deployments", "err", err.Error())
	} else {
		_ = level.Info(c.logger).Log("DeleteStep1", "Delete Deployments", "Result", "Success")
	}

	//step2 delete service
	if err = c.k8sClient.Do().CoreV1().Services(ns).Delete(name, nil); err != nil {
		_ = level.Error(c.logger).Log("DeleteStep2", "Delete Services", "err", err.Error())
	} else {
		_ = level.Info(c.logger).Log("DeleteStep2", "Delete Services", "Result", "Success")
	}

	//step3 delete virtualservice

	//step4 delete ingress
	if err = c.k8sClient.Do().ExtensionsV1beta1().Ingresses(ns).Delete(name, &metav1.DeleteOptions{}); err != nil {
		_ = level.Error(c.logger).Log("DeleteStep4", "Delete Ingress", "err", err.Error())
	} else {
		_ = level.Info(c.logger).Log("DeleteStep4", "Delete Ingress", "Result", "Success")
	}

	//step5 delete configmap
	if err = c.k8sClient.Do().CoreV1().ConfigMaps(ns).Delete(name, &metav1.DeleteOptions{}); err != nil {
		_ = level.Error(c.logger).Log("DeleteStep5", "Delete ConfigMaps", "err", err.Error())
	} else {
		_ = level.Info(c.logger).Log("DeleteStep5", "Delete ConfigMaps", "Result", "Success")
	}

	//step6 delete DB config_map && config_data
	if err = c.repository.ConfigMap().DeleteByNsName(ns, name); err != nil {
		_ = level.Error(c.logger).Log("DeleteStep6", "Delete DB ConfigMap And ConfigData", "err", err.Error())
	} else {
		_ = level.Info(c.logger).Log("DeleteStep6", "Delete DB ConfigMap And ConfigData", "Result", "Success")
	}

	//step7 delete DB virtual_service
	//step8 delete DB virtual_host
	//step9 delete DB webhooks

	//step10 delete DB builds
	if err = c.repository.Build().Delete(ns, name); err != nil {
		_ = level.Error(c.logger).Log("DeleteStep10", "Delete DB Builds", "err", err.Error())
	} else {
		_ = level.Info(c.logger).Log("DeleteStep10", "Delete DB Builds", "Result", "Success")
	}

	//step11 delete DB project_jenkins
	if err = c.repository.ProjectJenkins().Delete(ns, name); err != nil {
		_ = level.Error(c.logger).Log("DeleteStep11", "Delete DB ProjectJenkins", "err", err.Error())
	} else {
		_ = level.Info(c.logger).Log("DeleteStep11", "Delete DB ProjectJenkins", "Result", "Success")
	}

	//step12 delete DB project_template
	project, notExist := c.repository.Project().FindByNsNameExist(ns, name)
	if notExist == false {
		if err = c.repository.ProjectTemplate().DeleteByProjectId(project.ID); err != nil {
			_ = level.Error(c.logger).Log("DeleteStep12", "Delete DB ProjectTemplates", "err", err.Error())
		} else {
			_ = level.Info(c.logger).Log("DeleteStep12", "Delete DB ProjectTemplates", "Result", "Success")
		}
	}

	//step13 delete DB project && groups
	if err = c.repository.Project().Delete(ns, name); err != nil {
		_ = level.Error(c.logger).Log("DeleteStep13", "Delete DB Project", "err", err.Error())
	} else {
		_ = level.Info(c.logger).Log("DeleteStep13", "Delete DB Project", "Result", "Success")
	}

	//step14 delete jenkins job
	if job, err := c.jenkins.GetJob(name + "." + ns); err == nil {
		if err = c.jenkins.DeleteJob(job); err != nil {
			_ = level.Error(c.logger).Log("DeleteStep14", "Delete Jenkins Job", "err", err.Error())
		} else {
			_ = level.Info(c.logger).Log("DeleteStep14", "Delete Jenkins Job", "Result", "Success")
		}
	} else {
		_ = level.Error(c.logger).Log("DeleteSetep14", "Delete Jenkins Job", "err", err.Error())
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.DeleteEvent,
			name, ns,
			fmt.Sprintf("项目删除: %v.%v", name, ns)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	//发送邮件通知
	return nil
}

func (c *service) parseTemplate(tpl *types.ProjectTemplate, wg *sync.WaitGroup) interface{} {
	defer wg.Done()

	switch repository.Kind(tpl.Kind) {
	case repository.Deployment:
		var err error
		var deployment *v1.Deployment
		_ = yaml.Unmarshal([]byte(tpl.FinalTemplate), &deployment)
		if deployment, err = c.k8sClient.Do().AppsV1().Deployments(deployment.Namespace).Get(deployment.Name, metav1.GetOptions{}); err != nil {
			_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		}
		return deployment
	case repository.Service:
		var svc *corev1.Service
		_ = yaml.Unmarshal([]byte(tpl.FinalTemplate), &svc)
		return svc
	case repository.Ingress:
		var ing *v1beta1.Ingress
		_ = yaml.Unmarshal([]byte(tpl.FinalTemplate), &ing)
		return ing
	case repository.ConfigMap:
		var cm corev1.ConfigMap
		_ = yaml.Unmarshal([]byte(tpl.FinalTemplate), &cm)
		return cm
	case repository.VirtualService:

	}
	return nil
}

func (c *service) parsePods(pod *corev1.Pod) map[string]interface{} {
	var isReady = true
	var message, lastMessage string
	var restartCount int32

	metrics := make(chan map[string]interface{})
	go pods.GetPodsMetrics(pod.Namespace, pod.Name, c.config.GetString("server", "heapster_url"), metrics)

	for _, val := range pod.Status.ContainerStatuses {
		if val.Name == pod.Labels["app"] {
			restartCount = val.RestartCount
			break
		}
	}

	for _, val := range pod.Status.Conditions {
		if val.Type != corev1.PodReady && val.Status == corev1.ConditionFalse {
			isReady = false
			message = val.Message
			break
		}
	}

	if !isReady {
		for _, val := range pod.Status.ContainerStatuses {
			if val.Name == pod.Labels["app"] {
				if val.State.Waiting != nil {
					message = val.State.Waiting.Message
				}
				if val.LastTerminationState.Terminated != nil {
					lastMessage = val.LastTerminationState.Terminated.Message
				}
				break
			}
		}
	}

	var memory, currMemory, currCpu, cpu interface{}
	for {
		data, ok := <-metrics
		if !ok {
			break
		}
		if m, ok := data["memory"]; ok {
			memory = m
		}
		if u, ok := data["cpu"]; ok {
			cpu = u
		}
		if m, ok := data["curr_memory"]; ok {
			currMemory = m
		}
		if u, ok := data["curr_cpu"]; ok {
			currCpu = u
		}
	}

	return map[string]interface{}{
		"name":          pod.Name,
		"node_name":     pod.Spec.NodeName,
		"status":        pod.Status.Phase,
		"restart_count": restartCount,
		"create_at":     pod.CreationTimestamp,
		"message":       message,
		"last_message":  lastMessage,
		"memory":        memory,
		"cpu":           cpu,
		"curr_cpu":      currCpu,
		"curr_memory":   currMemory,
	}
}

func (c *service) rewriteTemplate(projectId int64, kind repository.TplKind, data map[string]interface{}) error {
	template, err := c.repository.Template().FindByKindType(kind)
	if err != nil {
		_ = level.Error(c.logger).Log("RewriteTemplate", "FindByKindType", "Kind", kind, "err", err.Error())
		return ErrTemplateGet
	}
	finalTemplate, err := encode.EncodeTemplate(kind.ToString(), template.Detail, data)
	if err != nil {
		_ = level.Error(c.logger).Log("RewriteTemplate", "EncodeTemplate", "err", err.Error())
		return ErrTemplateEncodeGet
	}

	field, _ := json.Marshal(data)
	if err = c.repository.ProjectTemplate().CreateOrUpdate(&types.ProjectTemplate{
		Kind:          kind.ToString(),
		ProjectID:     projectId,
		FinalTemplate: finalTemplate,
		Fields:        string(field),
	}); err != nil {
		_ = level.Error(c.logger).Log("RewriteTemplate", "CreateOrUpdate", "err", err.Error())
		return ErrProjectTemplateUpdate
	}
	return nil
}

func (c *service) saveJenkins(project *types.Project, req basicRequest) error {
	var kind repository.TplKind
	name := strings.Split(helper.GitName(req.GitAddr), "/")
	dat := map[string]string{
		"app_name":  project.Name,
		"git_name":  name[1],
		"git_path":  helper.GitUrl(req.GitAddr),
		"namespace": project.Namespace,
	}

	switch project.Language {
	case "Java":
		dat["build_path"] = helper.FormatBuildPath(req.BuildPath)
		kind = repository.JenkinsJavaCommand
	case "NodeJs":
		kind = repository.JenkinsNodeCommandKind
	case "Python":
		dat["build_path"] = req.BuildPath
		kind = repository.PythonKind
	case "Golang":
		dat["build_path"] = req.BuildPath
		kind = repository.JenkinsCommand
	default:
		kind = repository.JenkinsCommand
	}

	template, err := c.repository.Template().FindByKindType(kind)
	if err != nil {
		_ = level.Error(c.logger).Log("SaveJenkins", "FindTemplate", "err", err.Error())
		return ErrTemplateGet
	}

	// add projectJenkins DB
	finalTemplate, err := encode.EncodeTemplate(kind.ToString(), template.Detail, dat)
	if err != nil {
		_ = level.Error(c.logger).Log("SaveJenkins", "EncodeTemplate", "err", err.Error())
		return ErrTemplateEncodeGet
	}
	if err = c.repository.ProjectJenkins().CreateOrUpdate(&types.ProjectJenkins{
		Name:       project.Name,
		Namespace:  project.Namespace,
		GitAddr:    req.GitAddr,
		GitType:    req.GitType,
		GitVersion: req.GitVersion,
		Command:    finalTemplate,
	}); err != nil {
		_ = level.Error(c.logger).Log("SaveJenkins", "Create ProjectJenkins", "err", err.Error())
		return ErrProjectTemplateUpdate
	}

	return nil
}

func (c *service) Workspace(ctx context.Context) (res []map[string]interface{}, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	isAdmin := ctx.Value(middleware.IsAdmin).(bool)
	memberId := ctx.Value(middleware.UserIdContext).(int64)

	var projectList []*types.Project
	if isAdmin == true {
		projectList, err = c.repository.Project().GetProjectByNsLimit(ns)
	} else {
		projectList, err = c.repository.Groups().GetIndexProjectByMemberIdAndNs(memberId, ns)
	}

	if err != nil {
		_ = level.Error(c.logger).Log("Workspace", "get project lists", "err", err.Error())
	}

	for _, v := range projectList {

		res = append(res, map[string]interface{}{
			"id":          v.ID,
			"title":       v.Name,
			"logo":        "https://niu.yirendai.com/kpl-logo-blue.png",
			"description": v.Desc,
			"member":      v.Member.Username,
			"updatedAt":   v.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			"href":        "/project/detail/" + v.Namespace + "/" + v.Name,
			"memberLink":  "/project/detail/" + v.Namespace + "/" + v.Name,
		})
	}

	return
}

func (c *service) Config(ctx context.Context) (res map[string]interface{}, err error) {
	gitAddr := c.config.GetString("git", "git_addr")
	res = map[string]interface{}{
		"git_type": c.config.GetString("git", "git_type"),
		"git_addr": helper.GitUrl(gitAddr),
		"domain":   c.config.GetString("server", "domain_suffix"),
	}
	return
}
