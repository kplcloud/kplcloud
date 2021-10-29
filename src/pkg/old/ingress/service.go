package ingress

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/hooks"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var (
	ErrIngressGet                   = errors.New("ingress信息获取失败")
	ErrIngressGetProject            = errors.New("project信息获取失败")
	ErrIngressGetProjectTemplate    = errors.New("projectTemplate信息获取失败")
	ErrIngressUpdateProjectTemplate = errors.New("projectTemplate信息更新失败")
	ErrIngressGetTemplate           = errors.New("template信息获取失败")
	ErrIngressCreateTemplate        = errors.New("模板匹配失败")
	ErrIngressK8sCreate             = errors.New("Ingress生成错误")
	ErrIngressExists                = errors.New("Ingress已经存在无法再次生成")
)

type Service interface {
	// 获取Ingress详情
	Get(ctx context.Context, ns string, name string) (res map[string]interface{}, err error)

	// Ingress列表
	List(ctx context.Context, ns string, page, limit int) (res map[string]interface{}, err error)

	// 创建Ingress
	Post(ctx context.Context, req postRequest) error

	// 获取没有Ingress 的Project列表
	GetNoIngressProject(ctx context.Context, ns string) (res []map[string]interface{}, err error)

	// 同步Ingress
	Sync(ctx context.Context, ns string) error

	// 初始化生成Ingress
	Generate(ctx context.Context) error
}

type service struct {
	logger       log.Logger
	config       *config.Config
	k8sClient    kubernetes.K8sClient
	repository   repository.Repository
	hookQueueSvc hooks.ServiceHookQueue
}

func NewService(logger log.Logger, config *config.Config,
	client kubernetes.K8sClient,
	repository repository.Repository,
	hookQueueSvc hooks.ServiceHookQueue) Service {
	return &service{
		logger:       logger,
		config:       config,
		k8sClient:    client,
		repository:   repository,
		hookQueueSvc: hookQueueSvc,
	}
}

/**
 * @Title 生成Ingress
 */
func (c *service) Generate(ctx context.Context) (err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	if tpl, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Ingress); err == nil && tpl != nil {
		return ErrIngressExists
	}

	ingTpl, err := c.repository.Template().FindByKindType(repository.IngressKind)
	if err != nil {
		_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
		return ErrIngressGetTemplate
	}

	projectTpl, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindByProjectId", "err", err.Error())
		return ErrIngressGetProjectTemplate
	}

	tplData := map[string]interface{}{
		"name":      project.Name,
		"namespace": project.Namespace,
		"rules": []map[string]interface{}{
			map[string]interface{}{
				"domain": c.domain(project.Namespace, project.Name),
				"paths": []map[string]interface{}{
					map[string]interface{}{
						"serviceName": project.Name,
						"port":        projectTpl.FieldStruct.Ports[0].Port,
					},
				},
			},
		},
	}
	tpl, err := encode.EncodeTemplate(repository.IngressKind.ToString(), ingTpl.Detail, tplData)
	if err != nil {
		_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
		return ErrIngressCreateTemplate
	}

	var ingress *v1beta1.Ingress

	if err = yaml.Unmarshal([]byte(tpl), &ingress); err != nil {
		_ = level.Error(c.logger).Log("yaml", "Unmarshal", "err", err.Error())
		return ErrIngressCreateTemplate
	}

	if ingress, err = c.k8sClient.Do().ExtensionsV1beta1().Ingresses(project.Namespace).Create(ingress); err != nil {
		_ = level.Error(c.logger).Log("Ingresses", "Create", "err", err.Error())
		return ErrIngressK8sCreate
	}

	go func() {
		// projecttemplate 生成
		b, _ := yaml.Marshal(ingress)
		bb, _ := json.Marshal(tplData)
		if err = c.repository.ProjectTemplate().Create(&types.ProjectTemplate{
			ProjectID:     project.ID,
			FinalTemplate: string(b),
			Fields:        string(bb),
			Kind:          repository.Ingress.String(),
			State:         1,
		}); err != nil {
			_ = level.Warn(c.logger).Log("projectTemplateRepository", "Create", "err", err.Error())
		}

		if vsTpl, err := c.repository.Template().FindByKindType(repository.VirtualServiceKind); err == nil {
			vsData := map[string]interface{}{
				"name":      project.Name,
				"namespace": project.Namespace,
				"gateways":  []string{project.Namespace},
				"hosts":     []string{c.domain(project.Namespace, project.Name)},
				"http": []map[string]interface{}{
					map[string]interface{}{
						"route": []map[string]interface{}{
							map[string]interface{}{
								"host":   project.Name,
								"port":   "8080",
								"weight": 100,
							},
						},
					},
				},
			}
			if _, err := encode.EncodeTemplate(repository.VirtualServiceKind.ToString(), vsTpl.Detail, vsData); err != nil {
				_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
			}
		}

		// todo 创建virtualservice
		// c.repository.Template().FindByKindType(repository.VirtualServiceKind)
	}()

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.IngressGateway,
			project.Name, project.Namespace,
			fmt.Sprintf("初始化 Ingress \n 应用: %v.%v", project.Name, project.Namespace)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}

func (c *service) domain(ns, name string) string {
	return fmt.Sprintf(c.config.GetString("server", "domain_suffix"), name, ns)
}

/**
 * @Title 获取单个Ingress信息
 */
func (c *service) Get(ctx context.Context, ns string, name string) (res map[string]interface{}, err error) {
	ing, err := c.k8sClient.Do().ExtensionsV1beta1().Ingresses(ns).Get(name, v1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Ingress", "Get", "err", err.Error())
		return nil, ErrIngressGet
	}
	ing.Kind = "Ingress"
	ing.APIVersion = "extensions/v1beta1"
	res = map[string]interface{}{
		"spec":               ing.Spec,
		"name":               name,
		"namespace":          ns,
		"createionTimestamp": ing.CreationTimestamp,
		"yaml":               "",
	}
	if ingressYaml, err := yaml.Marshal(ing); err == nil {
		res["yaml"] = string(ingressYaml)
	}
	return
}

/**
 * @Title 获取Ingress列表
 */
func (c *service) List(ctx context.Context, ns string, page, limit int) (res map[string]interface{}, err error) {
	count, err := c.repository.ProjectTemplate().Count(ns, repository.Ingress)
	if err != nil {
		_ = level.Error(c.logger).Log("Ingress", "List Count", "err", err.Error())
		return nil, ErrIngressGetProjectTemplate
	}
	p := paginator.NewPaginator(page, limit, count)
	list, err := c.repository.ProjectTemplate().FindOffsetLimit(ns, repository.Ingress, p.Offset(), limit)
	if err != nil {
		_ = level.Error(c.logger).Log("Ingress", "List ProjectTemplate", "err", err.Error())
		return nil, ErrIngressGetProjectTemplate
	}
	var listData []map[string]interface{}

	for _, v := range list {
		var result = map[string]interface{}{
			"id":          v.ID,
			"projectId":   v.Project.ID,
			"projectName": v.Project.Name,
			"createdAt":   v.CreatedAt,
			"namespace":   ns,
			"spec":        v.IngressFieldStruct.Rules,
		}
		listData = append(listData, result)
	}
	res = map[string]interface{}{
		"list": listData,
		"page": p.Result(),
	}
	return
}

/**
 * @Title 创建更新Ingress
 */
func (c *service) Post(ctx context.Context, req postRequest) error {

	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	//get template and encode it
	serviceTemp, err := c.repository.Template().FindByKindType(repository.IngressKind)
	if err != nil {
		_ = level.Error(c.logger).Log("Ingress", "Post", "err", err.Error())
		return ErrIngressGetTemplate
	}
	var a map[string]interface{}
	b, _ := json.Marshal(req)
	_ = json.Unmarshal(b, &a)
	finalTemplate, err := encode.EncodeTemplate(repository.IngressKind.ToString(), serviceTemp.Detail, a)
	if err != nil {
		_ = level.Error(c.logger).Log("Ingress", "Post EncodeTemplate", "err", err.Error())
		return ErrIngressCreateTemplate
	}

	var ingress *v1beta1.Ingress
	if err = yaml.Unmarshal([]byte(finalTemplate), &ingress); err != nil {
		_ = level.Error(c.logger).Log("Ingress", "Post", "Yaml", "Unmarshal", "err", err.Error())
		return ErrIngressCreateTemplate
	}
	if _, err = c.k8sClient.Do().ExtensionsV1beta1().Ingresses(req.Namespace).Get(req.Name, v1.GetOptions{}); err == nil {
		ingress, err = c.k8sClient.Do().ExtensionsV1beta1().Ingresses(req.Namespace).Update(ingress)
		_ = level.Error(c.logger).Log("Ingress", "Post", "Ingress", "Update")
	} else {
		ingress, err = c.k8sClient.Do().ExtensionsV1beta1().Ingresses(req.Namespace).Create(ingress)
		_ = level.Error(c.logger).Log("Ingress", "Post", "Ingress", "Create")
	}
	if err != nil {
		_ = level.Error(c.logger).Log("Ingress", "Post", "err", err.Error())
		return err
	}

	// update projectTemplate database
	projectTemplate, err := c.repository.ProjectTemplate().FirstOrCreate(project.ID, repository.Ingress, string(b), finalTemplate, 1)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FirstOrCreate", "err", err.Error())
		return ErrIngressGetProjectTemplate
	}

	projectTemplate.FinalTemplate = finalTemplate
	projectTemplate.Fields = string(b)
	if err = c.repository.ProjectTemplate().UpdateTemplate(projectTemplate); err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "UpdateTemplate", "err", err.Error())
		return ErrIngressUpdateProjectTemplate
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.IngressGateway,
			project.Name, project.Namespace,
			fmt.Sprintf("创建或更新 Ingress \n 应用: %v.%v", project.Name, project.Namespace)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}

/**
 * @Title 获取不含Ingress的项目列表
 */
func (c *service) GetNoIngressProject(ctx context.Context, ns string) (res []map[string]interface{}, err error) {
	projects, err := c.repository.Project().GetProjectByNs(ns)
	if err != nil {
		_ = level.Error(c.logger).Log("Ingress", "GetNoIngressProject", "err", err.Error())
		return nil, ErrIngressGetProject
	}

	for _, v := range projects {
		if _, err = c.repository.ProjectTemplate().FindByProjectId(v.ID, repository.Ingress); err == nil {
			continue
		}
		res = append(res, map[string]interface{}{
			"id":        v.ID,
			"name":      v.Name,
			"namespace": v.Namespace,
		})
	}
	return res, nil

}

func (c *service) Sync(ctx context.Context, ns string) error {

	list, err := c.k8sClient.Do().ExtensionsV1beta1().Ingresses(ns).List(v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, v := range list.Items {
		name := strings.Replace(v.Name, "-ingress", "", -1)
		project, err := c.repository.Project().FindByNsName(v.Namespace, name)
		if err != nil {
			_ = level.Warn(c.logger).Log("Ingress", "Sync", "Name", name, "err", err.Error())
			continue
		}
		var rule []types.RuleStruct

		for _, val := range v.Spec.Rules {
			var path []*types.Paths
			for _, vv := range val.HTTP.Paths {
				path = append(path, &types.Paths{
					Path:        vv.Path,
					PortName:    vv.Backend.ServicePort.IntValue(),
					ServiceName: vv.Backend.ServiceName,
				})
			}
			rule = append(rule, types.RuleStruct{
				Domain: val.Host,
				Paths:  path,
			})
		}

		fields, _ := json.Marshal(rule)
		ingressYml, _ := yaml.Marshal(v)

		_, err = c.repository.ProjectTemplate().FirstOrCreate(project.ID, repository.Ingress, string(fields), string(ingressYml), 1)
		if err != nil {
			_ = level.Warn(c.logger).Log("Ingress", "Sync", "ProjectTemplate", "FirstOrCreate", "err", err.Error())
		}

		_ = level.Info(c.logger).Log("Ingress", "Sync", "Name", name)
	}
	return nil
}
