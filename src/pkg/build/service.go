package build

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	amqpClient "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/jenkins"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/hooks"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	"gopkg.in/guregu/null.v3"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/url"
	"strings"
	"time"
)

var (
	ErrBuildGet                 = errors.New("获取Build记录错误,请查询是否存在")
	ErrBuilding                 = errors.New("需要等上一个构建完成才能进行新的构建,你可以选择暂停")
	ErrBuildJenkinsJob          = errors.New("Jenkins Build Job错误:")
	ErrBuildProjectGet          = errors.New("项目获取失败，可能不存在或未审核")
	ErrBuildJenkinsJobGet       = errors.New("Jenkins Job 获取错误: ")
	ErrBuildCreate              = errors.New("创建Build记录错误,请联系管理员")
	ErrBuildQueuePublish        = errors.New("构建入列出错了,请联系管理员")
	ErrBuildAbort               = errors.New("终止错误,可能这个job状态已经是终止态")
	ErrBuildCount               = errors.New("构建记录统计出错")
	ErrBuildList                = errors.New("构建列表获取出错")
	ErrBuildDeploymentK8sGet    = errors.New("基础信息获取错误")
	ErrBuildDeploymentK8sUpdate = errors.New("更新基础镜像出错")
)

type Service interface {
	// build 应用
	Build(ctx context.Context, gitType, version, buildEnv, buildEnvDesc, buildTime string) error

	// 获取build输出的信息
	BuildConsole(ctx context.Context, number, start int) (string, int, error)

	// 处理消费出来的 build 的数据
	// desc: 如果还在build 或出有一点啥差错 重新放回队列
	ReceiverBuild(ctx context.Context, data string) error

	// 终止构建
	AbortBuild(ctx context.Context, jenkinsBuildId int) error

	// 获取build 记录
	History(ctx context.Context, page, limit int) (map[string]interface{}, error)

	// 回滚版本
	Rollback(ctx context.Context, buildId int64) error

	// 获取Jenkins build配置信息
	BuildConf(ctx context.Context, ns, name string) (res interface{}, err error)

	// 获取cronjob build 记录
	CronHistory(ctx context.Context, page, limit int) (map[string]interface{}, error)

	// 获取cronjob build输出的信息
	CronBuildConsole(ctx context.Context, number, start int) (string, int, error)
}

type service struct {
	logger       log.Logger
	jenkins      jenkins.Jenkins
	amqpClient   amqpClient.AmqpClient
	k8sClient    kubernetes.K8sClient
	config       *config.Config
	repository   repository.Repository
	hookQueueSvc hooks.ServiceHookQueue
}

func NewService(logger log.Logger,
	jenkins jenkins.Jenkins,
	amqpClient amqpClient.AmqpClient,
	k8sClient kubernetes.K8sClient,
	config *config.Config,
	store repository.Repository,
	hookQueueSvc hooks.ServiceHookQueue) Service {
	return &service{
		logger,
		jenkins,
		amqpClient,
		k8sClient,
		config,
		store,
		hookQueueSvc,
	}
}

func (c *service) Rollback(ctx context.Context, buildId int64) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)
	build, err := c.repository.Build().FindById(project.Namespace, project.Name, buildId)
	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "FindById", "err", err.Error())
		return ErrBuildGet
	}
	var deployment *v1.Deployment

	defer func() {
		if err == nil {
			if tpl, e := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment); e == nil {
				b, _ := yaml.Marshal(deployment)
				tpl.FinalTemplate = string(b)
				tpl.FieldStruct.Image = build.Version
				if ee := c.repository.ProjectTemplate().UpdateTemplate(tpl); ee != nil {
					_ = level.Error(c.logger).Log("projectTemplateRepository", "UpdateTemplate", "err", ee.Error())
				}
			}
		}
	}()

	deployment, err = c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Get(project.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrBuildDeploymentK8sGet
	}

	for k, v := range deployment.Spec.Template.Spec.Containers {
		if v.Name != project.Name {
			continue
		}
		image := strings.Split(v.Image, ":")
		deployment.Spec.Template.Spec.Containers[k].Image = image[0] + ":" + build.Version
		break
	}

	deployment, err = c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Update(deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrBuildDeploymentK8sUpdate
	}

	build.Status = null.StringFrom("ROLLBACK")
	if _, e := c.repository.Build().CreateBuild(&build); e != nil {
		_ = level.Error(c.logger).Log("buildRepository", "CreateBuild", "err", err.Error())
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.RollbackEvent,
			project.Name, project.Namespace,
			fmt.Sprintf("项目回滚: %v.%v, 版本：%v", project.Name, project.Namespace, build.Version)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}

func (c *service) History(ctx context.Context, page, limit int) (map[string]interface{}, error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	total, err := c.repository.Build().Count(project.Namespace, project.Name)
	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "Count", "err", err.Error())
		return nil, ErrBuildCount
	}

	p := paginator.NewPaginator(page, limit, int(total))
	builds, err := c.repository.Build().FindOffsetLimit(project.Namespace, project.Name, p.Offset(), limit)
	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "FindOffsetLimit", "err", err.Error())
		return nil, ErrBuildList
	}

	return map[string]interface{}{
		"list": builds,
		"page": p.Result(),
	}, nil
}

func (c *service) AbortBuild(ctx context.Context, jenkinsBuildId int) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	build, err := c.repository.Build().FindById(project.Namespace, project.Name, int64(jenkinsBuildId))
	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "FindById", "err", err.Error())
		return ErrBuildAbort
	}

	defer func() {
		// 更新数据库
		build.Status = null.StringFrom("ABORTED")
		if err = c.repository.Build().Update(&build); err != nil {
			_ = level.Error(c.logger).Log("buildRepository", "Update", "err", err.Error())
		}
	}()
	jobName := project.Name + "." + project.Namespace
	if err = c.jenkins.AbortJob(jobName, int(build.BuilderID)); err != nil {
		_ = level.Error(c.logger).Log("jenkins", "AbortJob", "err", err.Error())
		return ErrBuildAbort
	}

	// TODO event 谁谁谁终止了构建

	return nil
}

func (c *service) Build(ctx context.Context, gitType, version, buildEnv, buildEnvDesc, buildTime string) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	var err error
	var buildTimer time.Time
	if buildTime != "" {
		buildTimer, err = time.ParseInLocation("2006-01-02 15:04:05", buildTime, time.Local)
		if err != nil {
			_ = level.Error(c.logger).Log("time", "ParseInLocation", "err", err.Error())
		}
	}

	if buildTimer != (time.Time{}) && buildTimer.Unix() > time.Now().Unix() {
		// todo 可以入列 暂不实现
	} else {
		buildTimer = time.Now()
	}

	if build, err := c.repository.Build().FirstByTag(project.Namespace, project.Name, version); err == nil && build != nil {
		return ErrBuilding
	}

	projectTpl, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindByProjectId", "err", err.Error())
		return ErrBuildProjectGet
	}

	jenkinsJobName := project.Name + "." + project.Namespace
	job, err := c.jenkins.GetJob(jenkinsJobName)
	if err != nil {
		_ = level.Error(c.logger).Log("jenkins", "GetJob", "err", err.Error())
		return ErrBuildJenkinsJobGet
	}

	var branch = "tags/" + version
	var tagName = version
	if gitType == "branch" {
		branch = "*/" + version
		tagName = version + "-" + time.Now().Format("2006012150405")
	}

	params := url.Values{
		"TAGNAME": []string{tagName},
		"BRANCH":  []string{branch},
	}

	if repository.Language(project.Language) == repository.Java {
		if buildEnv != "" {
			buildEnv = "-P" + buildEnv
		}
		params.Set("POMFILE", projectTpl.FieldStruct.PomFile)
		params.Set("BUILD_ENV", buildEnv)
	}

	if err := c.jenkins.Build(jenkinsJobName, params); err != nil {
		return errors.New(ErrBuildJenkinsJob.Error() + err.Error())
	}

	lastBuild, err := c.jenkins.GetLastBuild(job)
	if err != nil {
		_ = level.Warn(c.logger).Log("jenkins", "GetLastBuild", "err", err.Error())
		lastBuild = jenkins.Build{}
	}

	// jenkins 会静莫几秒
	lastBuildId := int64(lastBuild.Number + 1)

	userId := ctx.Value(middleware.UserIdContext).(int64)
	build, err := c.repository.Build().CreateBuild(&types.Build{
		Address:   null.StringFrom(projectTpl.FieldStruct.GitAddr),
		BuildEnv:  null.StringFrom(buildEnv),
		BuildID:   null.IntFrom(lastBuildId),
		BuildTime: null.TimeFrom(buildTimer),
		BuilderID: userId,
		Name:      project.Name,
		Namespace: project.Namespace,
		GitType:   null.StringFrom(gitType),
		Status:    null.StringFrom("BUILDING"),
		Version:   version,
	})

	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "CreateBuild", "err", err.Error())
		return ErrBuildCreate
	}

	if err := c.amqpClient.PublishOnQueue(amqpClient.BuildTopic, func() []byte {
		b, _ := json.Marshal(&amqpClient.BuildPublishData{
			Name:           project.Name,
			Namespace:      project.Namespace,
			JenkinJobName:  jenkinsJobName,
			JenkinsBuildId: lastBuildId,
			Builder:        userId,
			BuildId:        build.ID,
			Version:        version,
		})
		return b
	}); err != nil {
		_ = level.Error(c.logger).Log("amqpClient", "PublishOnQueue", "err", err.Error())
		return ErrBuildQueuePublish
	}

	//event history
	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.BuildEvent,
			project.Name, project.Namespace,
			fmt.Sprintf("项目 %v.%v, Build版本：%v", project.Name, project.Namespace, version)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}

func (c *service) BuildConsole(ctx context.Context, number, start int) (string, int, error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	buildInfo, err := c.repository.Build().FindById(project.Namespace, project.Name, int64(number))
	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "FindBuildByBuildId", "err", err.Error())
		return "", start, ErrBuildGet
	}

	output := buildInfo.Output.String
	end := len(output)
	if start != 0 {
		output = output[start:end]
	}
	if len(output) == 0 {
		output = "."
	}

	return output, end, nil
}

func (c *service) ReceiverBuild(ctx context.Context, data string) (err error) {
	// 获取jenkins build 状态 然后入库
	// build成功之后 调用deployment更新版本
	if len(data) <= 0 {
		return nil
	}

	defer func() {
		if err != nil {
			if e := c.amqpClient.PublishOnQueue(amqpClient.BuildTopic, func() []byte {
				return []byte(data)
			}); e != nil {
				_ = level.Error(c.logger).Log("amqpClient", "PublishOnQueue", "err", err.Error())
			}
		}
		time.Sleep(time.Second * 2)
	}()

	var receiverData amqpClient.BuildPublishData
	_ = json.Unmarshal([]byte(data), &receiverData)

	job, err := c.jenkins.GetJob(receiverData.JenkinJobName)
	if err != nil {
		_ = level.Error(c.logger).Log("jenkins", "GetJob", "err", err.Error())
		return err
	}

	build, err := c.jenkins.GetBuild(job, int(receiverData.JenkinsBuildId))
	if err != nil {
		_ = level.Error(c.logger).Log("jenkins", "GetBuild", "err", err.Error())
		return err
	}

	buildRes, err := c.repository.Build().FindBuildByBuildId(receiverData.Namespace, receiverData.Name, int(receiverData.JenkinsBuildId))
	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "FindBuildByBuildId", "err", err.Error())
		return nil
	}

	resBody, err := c.jenkins.GetBuildConsoleOutput(build)
	if err != nil {
		_ = level.Error(c.logger).Log("jenkins", "GetBuildConsoleOutput", "err", err.Error())
		return err
	}

	if build.Building {
		buildRes.Status = null.StringFrom("BUILDING")
	} else {
		buildRes.Status = null.StringFrom(strings.ToUpper(build.Result))
	}

	buildRes.Output = null.StringFrom(string(resBody))
	go func() {
		if e := c.repository.Build().Update(&buildRes); e != nil {
			_ = level.Error(c.logger).Log("buildRepository", "Update", "err", e.Error())
		}
	}()

	buildResult := strings.ToUpper(build.Result)

	defer func() {
		if "SUCCESS" == buildResult || "FAILURE" == buildResult {
			go func() {
				m, _ := c.repository.Member().FindById(receiverData.Builder)
				p, _ := c.repository.Project().FindByNsNameOnly(receiverData.Namespace, receiverData.Name)
				ctx = context.WithValue(ctx, middleware.ProjectContext, p)
				ctx = context.WithValue(ctx, middleware.UserIdContext, receiverData.Builder)
				msg := fmt.Sprintf(`
版本: %s
Build时间: %s
构建状态： %s
操作人: %s
详情: %s
`, receiverData.Version, time.Now().Format("2006/01/02 15:04:05"), buildResult, m.Username, build.URL)

				if err := c.hookQueueSvc.SendHookQueue(ctx,
					repository.BuildEvent,
					receiverData.Name, receiverData.Namespace,
					msg); err != nil {
					_ = level.Error(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
				}
			}()
			return
		}
	}()
	if "SUCCESS" == buildResult {
		// 调用deployment api 更新image
		deployment, e := c.k8sClient.Do().AppsV1().Deployments(receiverData.Namespace).Get(receiverData.Name, metav1.GetOptions{})
		if e != nil {
			_ = level.Error(c.logger).Log("Deployments", "Get", "err", e.Error())
			return nil
		}
		dockerHub := c.config.GetString("server", "docker_repo")
		for key, container := range deployment.Spec.Template.Spec.Containers {
			if container.Name == receiverData.Name {
				deployment.Spec.Template.Spec.Containers[key].Image = dockerHub + "/" + receiverData.Namespace + "/" + receiverData.Name + ":" + receiverData.Version
				break
			}
		}
		if deployment, e = c.k8sClient.Do().AppsV1().Deployments(receiverData.Namespace).Update(deployment); e == nil {
			if project, eer := c.repository.Project().FindByNsName(receiverData.Namespace, receiverData.Name); eer == nil {
				if projectTpl, eeer := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment); eeer == nil {
					projectTpl.FieldStruct.Image = dockerHub + "/" + receiverData.Namespace + "/" + receiverData.Name + ":" + receiverData.Version
					b, _ := json.Marshal(deployment)
					projectTpl.FinalTemplate = string(b)
					if eeeer := c.repository.ProjectTemplate().UpdateTemplate(projectTpl); eeeer != nil {
						_ = level.Error(c.logger).Log("build", "status", "projectTemplateRepository", "UpdateTemplate", "eeeer", eeeer.Error())
					}
				} else {
					_ = level.Error(c.logger).Log("projectTemplateRepository", "FindByProjectId", "err", eeer.Error())
				}
			} else {
				_ = level.Error(c.logger).Log("projectRepository", "FindByNsName", "err", eer.Error())
			}
		} else {
			// todo 给管理员推 更新镜像失败
			_ = level.Error(c.logger).Log("Deployments", "Update", "err", e.Error())
		}
		// todo 监听Pods的起动情况

		return nil
	}
	if buildResult == "" {
		buildResult = "BUILDING"
	}

	if "FAILURE" == buildResult {
		return nil
	}

	return errors.New("build status:" + buildResult)
}

func (c *service) BuildConf(ctx context.Context, ns, name string) (res interface{}, err error) {
	projectJenkins, err := c.repository.ProjectJenkins().Find(ns, name)
	return projectJenkins, err
}

func (c *service) CronHistory(ctx context.Context, page, limit int) (map[string]interface{}, error) {
	cronjob := ctx.Value(middleware.CronJobContext).(*types.Cronjob)

	total, err := c.repository.Build().Count(cronjob.Namespace, cronjob.Name+"-cronjob")
	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "Count", "err", err.Error())
		return nil, ErrBuildCount
	}

	p := paginator.NewPaginator(page, limit, int(total))
	builds, err := c.repository.Build().FindOffsetLimit(cronjob.Namespace, cronjob.Name+"-cronjob", p.Offset(), limit)
	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "FindOffsetLimit", "err", err.Error())
		return nil, ErrBuildList
	}

	return map[string]interface{}{
		"list": builds,
		"page": p.Result(),
	}, nil
}

func (c *service) CronBuildConsole(ctx context.Context, number, start int) (string, int, error) {
	cronjob := ctx.Value(middleware.CronJobContext).(*types.Cronjob)

	buildInfo, err := c.repository.Build().FindById(cronjob.Namespace, cronjob.Name+"-cronjob", int64(number))
	if err != nil {
		_ = level.Error(c.logger).Log("buildRepository", "FindBuildByBuildId", "err", err.Error())
		return "", start, ErrBuildGet
	}

	output := buildInfo.Output.String
	end := len(output)
	if start != 0 {
		output = output[start:end]
	}
	if len(output) == 0 {
		output = "."
	}

	return output, end, nil
}
