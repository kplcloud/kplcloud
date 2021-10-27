package audit

import (
	"context"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/jenkins"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/build"
	"github.com/kplcloud/kplcloud/src/pkg/hooks"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	ErrUpdateProject         = errors.New("project update error")
	ErrUpdateProjectTemplate = errors.New("project template update error")
	ErrAduitStepRefused      = errors.New("当前模块不支持重新发布")
)

type Service interface {
	AccessAudit(ctx context.Context, ns, name string) error
	AuditStep(ctx context.Context, ns, name, kind string) error
	Refused(ctx context.Context, ns, name string) error
}

type service struct {
	logger       log.Logger
	config       *config.Config
	jenkins      jenkins.Jenkins
	k8sClient    kubernetes.K8sClient
	amqpClient   amqp.AmqpClient
	repository   repository.Repository
	hookQueueSvc hooks.ServiceHookQueue
	buildSvc     build.Service
}

func NewService(logger log.Logger,
	config *config.Config,
	jenkins jenkins.Jenkins,
	k8sClient kubernetes.K8sClient,
	amqpClient amqp.AmqpClient,
	store repository.Repository,
	hookQueueSvc hooks.ServiceHookQueue,
	buildSvc build.Service) Service {
	return &service{
		logger,
		config,
		jenkins,
		k8sClient,
		amqpClient,
		store,
		hookQueueSvc,
		buildSvc,
	}
}

func (c *service) AccessAudit(ctx context.Context, ns, name string) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)
	projectTemplates, err := c.repository.ProjectTemplate().FindProjectTemplateByProjectId(project.ID)
	if err != nil {
		_ = level.Error(c.logger).Log("AccessAudit", "FindProjectTemplateByProjectId", "err", err.Error())
		return err
	}

	projectJenkins, err := c.repository.ProjectJenkins().Find(project.Namespace, project.Name)
	if err != nil {
		_ = level.Error(c.logger).Log("Audit", "createJob", "projectJenkins", "Find", "err", err.Error())
		return err
	}

	//创建Jenkins
	jenkinsParams := jenkins.Params{
		Namespace:  ns,
		Name:       name,
		GitAddr:    projectJenkins.GitAddr,
		GitType:    projectJenkins.GitType,
		GitVersion: projectJenkins.GitVersion,
		Command:    projectJenkins.Command,
	}
	if project.Language != repository.Java.String() {
		if err := c.jenkins.CreateJobParams(jenkinsParams); err != nil {
			_ = level.Error(c.logger).Log("Aduit", "CreateJob", "err", err.Error())
			return err
		}
	} else {
		if err := c.jenkins.CreateJavaJobParams(jenkinsParams); err != nil {
			_ = level.Error(c.logger).Log("Audit", "CreateJob", "err", err.Error())
			return err
		}
	}

	// jenkins build
	if err := c.buildSvc.Build(ctx, projectJenkins.GitType, projectJenkins.GitVersion, "", "", ""); err != nil {
		_ = level.Error(c.logger).Log("Audit", "Build", "err", err.Error())
	}

	for _, v := range projectTemplates {
		// 创建configMap
		if v.Kind == repository.ConfigMap.String() {
			var configMap *corev1.ConfigMap
			_ = yaml.Unmarshal([]byte(v.FinalTemplate), &configMap)
			if _, err := c.k8sClient.Do().CoreV1().ConfigMaps(project.Namespace).Create(configMap); err != nil {
				_ = level.Error(c.logger).Log("AccessAudit", "ConfigMap Create", "err", err.Error())
				v.State = 2
			} else {
				v.State = 1
			}
			if err := c.repository.ProjectTemplate().UpdateState(v); err != nil {
				_ = level.Error(c.logger).Log("AccessAudit", "ConfigMap UpdateProjectTemplate", "err", err.Error())
			}
			continue
		}

		// 创建Deployment
		if v.Kind == repository.Deployment.String() {
			var deployment *v1.Deployment
			_ = yaml.Unmarshal([]byte(v.FinalTemplate), &deployment)
			if _, err := c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Create(deployment); err != nil {
				_ = level.Error(c.logger).Log("AccessAudit", "Deployment Create", "err", err.Error())
				v.State = 2
			} else {
				v.State = 1
			}
			if err := c.repository.ProjectTemplate().UpdateState(v); err != nil {
				_ = level.Error(c.logger).Log("AccessAudit", "Deployment UpdateProjectTemplate", "err", err.Error())
			}
			continue
		}

		// 创建Service
		if v.Kind == repository.Service.String() {
			var service *corev1.Service
			_ = yaml.Unmarshal([]byte(v.FinalTemplate), &service)
			if _, err := c.k8sClient.Do().CoreV1().Services(project.Namespace).Create(service); err != nil {
				_ = level.Error(c.logger).Log("AccessAudit", "Services Create", "err", err.Error())
				v.State = 2
			} else {
				v.State = 1
			}
			if err := c.repository.ProjectTemplate().UpdateState(v); err != nil {
				_ = level.Error(c.logger).Log("AccessAudit", "Services UpdateProjectTemplate", "err", err.Error())
			}
			continue
		}
	}

	project.PublishState = int64(repository.PublishPass)
	project.AuditState = int64(repository.AuditPass)
	if err := c.repository.Project().UpdateProjectById(project); err != nil {
		_ = level.Error(c.logger).Log("AccessAudit", "Update Project State", "err", err.Error())
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.AuditEvent,
			name, ns,
			fmt.Sprintf("项目审核通过: %v.%v", name, ns)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}

func (c *service) AuditStep(ctx context.Context, ns, name, kind string) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	switch kind {
	case repository.ConfigMap.String():
		projectTpl, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.ConfigMap)
		if err != nil {
			return err
		}
		var configMap *corev1.ConfigMap
		_ = yaml.Unmarshal([]byte(projectTpl.FinalTemplate), &configMap)
		if _, err := c.k8sClient.Do().CoreV1().ConfigMaps(project.Namespace).Create(configMap); err != nil {
			_ = level.Error(c.logger).Log("AuditStep", "ConfigMap Create", "err", err.Error())
			projectTpl.State = 2
		} else {
			projectTpl.State = 1
		}
		if err := c.repository.ProjectTemplate().UpdateState(projectTpl); err != nil {
			_ = level.Error(c.logger).Log("AuditStep", "ConfigMap UpdateProjectTemplate", "err", err.Error())
			return err
		}
	case repository.Deployment.String():
		projectTpl, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
		if err != nil {
			return err
		}
		var deployment *v1.Deployment
		_ = yaml.Unmarshal([]byte(projectTpl.FinalTemplate), &deployment)
		if _, err := c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Create(deployment); err != nil {
			_ = level.Error(c.logger).Log("AccessAudit", "Deployment Create", "err", err.Error())
			projectTpl.State = 2
		} else {
			projectTpl.State = 1
		}
		if err := c.repository.ProjectTemplate().UpdateState(projectTpl); err != nil {
			_ = level.Error(c.logger).Log("AccessAudit", "Deployment UpdateProjectTemplate", "err", err.Error())
			return err
		}
	case repository.Service.String():
		projectTpl, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Service)
		if err != nil {
			return err
		}
		var service *corev1.Service
		_ = yaml.Unmarshal([]byte(projectTpl.FinalTemplate), &service)
		if _, err := c.k8sClient.Do().CoreV1().Services(project.Namespace).Create(service); err != nil {
			_ = level.Error(c.logger).Log("AccessAudit", "Services Create", "err", err.Error())
			projectTpl.State = 2
		} else {
			projectTpl.State = 1
		}
		if err := c.repository.ProjectTemplate().UpdateState(projectTpl); err != nil {
			_ = level.Error(c.logger).Log("AccessAudit", "Services UpdateProjectTemplate", "err", err.Error())
			return err
		}
	default:
		_ = level.Error(c.logger).Log("AuditStep", "Kind Refused")
		return ErrAduitStepRefused
	}
	return nil
}

func (c *service) Refused(ctx context.Context, ns, name string) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	// 1. update project DB
	project.AuditState = int64(repository.AuditFail)
	project.Step = 1
	if err := c.repository.Project().UpdateProjectById(project); err != nil {
		_ = level.Error(c.logger).Log("Refused", "UpdateProject", "err", err.Error())
		return ErrUpdateProject
	}

	// 2. update projectTemplate DB
	if err := c.repository.ProjectTemplate().DeleteByProjectId(project.ID); err != nil {
		_ = level.Error(c.logger).Log("Refused", "Delete ProjectTemplate", "err", err.Error())
		return ErrUpdateProjectTemplate
	}

	// 3. delete projectJenkins
	if err := c.repository.ProjectJenkins().Delete(project.Namespace, project.Name); err != nil {
		_ = level.Error(c.logger).Log("Refused", "Delete ProjectJenkins", "err", err.Error())
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.AuditEvent,
			name, ns,
			fmt.Sprintf("项目审核驳回: %v.%v", name, ns)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()
	return nil
}
