/**
 * @Time : 2019-07-23 18:34
 * @Author : solacowa@gmail.com
 * @File : Service
 * @Software: GoLand
 */

package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	crd "github.com/kplcloud/kplcloud/src/istio/types/v1beta1"
	"github.com/kplcloud/kplcloud/src/jenkins"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"gopkg.in/guregu/null.v3"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgtypes "k8s.io/apimachinery/pkg/types"
	"time"
)

var (
	ErrToolDeploymentK8sGet    = errors.New("项目获取错误,可能不存在")
	ErrToolDeploymentK8sUpdate = errors.New("修改容器时间错误，请重试")
	ErrToolFakeTimeErr         = errors.New("时间格式错误,请重新提交")
	ErrToolProjectGet          = errors.New("源项目获取错误,可能不存在")
	ErrToolProjectExists       = errors.New("目标空间已经存在该应用无法克隆")
	ErrToolProjectTemplateGet  = errors.New("应用模版获取错误")
	ErrToolProjectCreate       = errors.New("应用创建错误")
	ErrToolJenkinsGet          = errors.New("项目已克隆完成,但获取源项目jenkins错误,请联系管理员")
)

type Service interface {
	// clone a project
	// 克隆一个服务
	Duplication(ctx context.Context, sourceNamespace, sourceAppName, destinationNamespace string) (err error)

	// 调整容器时间
	FakeTime(ctx context.Context, fakeTime time.Time, method FakeTimeMethod) (err error)
}

type service struct {
	logger     log.Logger
	config     *config.Config
	jenkins    jenkins.Jenkins
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func NewService(logger log.Logger, config *config.Config,
	jenkins jenkins.Jenkins,
	k8sClient kubernetes.K8sClient,
	repository repository.Repository) Service {
	return &service{level.Error(logger), config,
		jenkins,
		k8sClient,
		repository}
}

func (c *service) Duplication(ctx context.Context, sourceNamespace, sourceAppName, destinationNamespace string) (err error) {
	// 源是否存在
	// 查询目标ns 是否已经有该应用了
	// 读取源ns下的应用所有数据,包括所有kind
	// 替换模版所有namespace
	// 生成新的模版
	// 调用k8s生成
	// 调用jenkins生成job

	userId := ctx.Value(middleware.UserIdContext).(int64)

	project, err := c.repository.Project().FindByNsName(sourceNamespace, sourceAppName)
	if err != nil {
		_ = level.Error(c.logger).Log("projectRepository", "FindByNsName", "err", err.Error())
		return ErrToolProjectGet
	}

	if destination, err := c.repository.Project().FindByNsName(destinationNamespace, sourceAppName); err == nil && destination != nil {
		return ErrToolProjectExists
	}

	tpls, err := c.repository.ProjectTemplate().FindProjectTemplateByProjectId(project.ID)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindProjectTemplateByProjectId", "err", err.Error())
		return ErrToolProjectTemplateGet
	}

	m, _ := c.repository.Member().FindById(userId)

	project.Namespace = destinationNamespace
	project.MemberID = userId
	project.ID = 0
	project.Member = *m
	project.CreatedAt = null.NewTime(time.Now(), true)
	project.UpdatedAt = null.NewTime(time.Now(), true)

	if err = c.repository.Project().Create(project); err != nil {
		_ = level.Error(c.logger).Log("projectRepository", "Create", "err", err.Error())
		return ErrToolProjectCreate
	}

	var deploymentFields types.TemplateField

	for _, tpl := range tpls {
		newTpl := &types.ProjectTemplate{
			ProjectID: project.ID,
			Project:   *project,
			Kind:      tpl.Kind,
			State:     1,
		}
		var fields map[string]interface{}
		if err = json.Unmarshal([]byte(tpl.Fields), &fields); err == nil {
			fields["namespace"] = destinationNamespace
			b, _ := json.Marshal(fields)
			newTpl.Fields = string(b)

			if tpl.Kind == repository.Deployment.String() {
				_ = json.Unmarshal([]byte(newTpl.Fields), &deploymentFields)
			}
		}
		newTpl.FinalTemplate = c.copyOnWrite(tpl, sourceAppName, sourceNamespace, destinationNamespace)
		if err = c.repository.ProjectTemplate().Create(newTpl); err != nil {
			_ = level.Error(c.logger).Log("projectTemplateRepository", "Create", "err", err.Error())
			continue
		}
	}

	// Jenkins{Language}Command
	commandKey := "JenkinsCommand"
	tpl, err := c.repository.Template().FindByKindType(repository.TplKind(commandKey))
	if err != nil {
		_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
		return ErrToolProjectTemplateGet
	}

	jobItem, err := c.jenkins.GetJobConfig(sourceAppName + "." + sourceNamespace)
	if err != nil {
		_ = level.Error(c.logger).Log("jenkins", "GetJobConfig", "err", err.Error())
		return ErrToolJenkinsGet
	}

	tplStr, err := encode.EncodeTemplate(commandKey, tpl.Detail, map[string]interface{}{
		"build_path": deploymentFields.BuildPath,
		"app_name":   project.Name,
		"namespace":  project.Namespace,
	})
	if err != nil {
		return err
	}

	jobItem.Builders.HudsonTasksShell.Command = tplStr
	jobItem.Description = jobItem.Description + fmt.Sprintf(" ps: 从 %s 克隆而来", sourceNamespace)

	if err = c.jenkins.CreateJob(jobItem, sourceAppName+"."+destinationNamespace); err != nil {
		_ = level.Error(c.logger).Log("jenkins", "CreateJob", "err", err.Error())
		return err
	}

	/*
		if _, err = c.jenkins.GetView("local"); err != nil {
			if e := c.jenkins.CreateView(jenkins.ListView{Name: "local", Properties: jenkins.ViewProperties{
				Class: "hudson.model.View$PropertyList",
			}}); e != nil {
				_ = level.Error(c.logger).Log("jenkins", "CreateView", "err", e.Error())
			}
		}

		job, err := c.jenkins.GetJob("hello.local")
		if err != nil {
			return err
		}

		if err = c.jenkins.AddJobToView("local", job); err != nil {
			_ = level.Error(c.logger).Log("jenkins", "AddJobToView", "err", err.Error())
			return err
		}*/

	return
}

func (c *service) copyOnWrite(tpl *types.ProjectTemplate, appName, sourceNamespace, destinationNamespace string) string {

	var err error
	switch repository.Kind(tpl.Kind) {
	case repository.Deployment:
		var deployment *appsv1.Deployment
		_ = yaml.Unmarshal([]byte(tpl.FinalTemplate), &deployment)
		deployment.ObjectMeta.Namespace = destinationNamespace
		deployment.ObjectMeta.Annotations = nil
		deployment.ObjectMeta.CreationTimestamp = metav1.Time{}
		deployment.ObjectMeta.Generation = 0
		deployment.ObjectMeta.ResourceVersion = ""
		deployment.ObjectMeta.UID = pkgtypes.UID("")
		deployment.SelfLink = ""
		deployment.Status = appsv1.DeploymentStatus{}

		if deployment, err = c.k8sClient.Do().AppsV1().Deployments(destinationNamespace).Create(deployment); err != nil {
			_ = level.Error(c.logger).Log("Deployments", "Create", "err", err.Error())
		}

		b, _ := yaml.Marshal(deployment)
		return string(b)
	case repository.Service:
		var svc *v1.Service
		_ = yaml.Unmarshal([]byte(tpl.FinalTemplate), &svc)
		svc.ObjectMeta.Namespace = destinationNamespace
		svc.ObjectMeta.Annotations = nil
		svc.ObjectMeta.CreationTimestamp = metav1.Time{}
		svc.ObjectMeta.Generation = 0
		svc.ObjectMeta.ResourceVersion = ""
		svc.SelfLink = ""
		svc.Status = v1.ServiceStatus{}
		svc.Spec.ClusterIP = ""

		if svc, err = c.k8sClient.Do().CoreV1().Services(destinationNamespace).Create(svc); err != nil {
			_ = level.Error(c.logger).Log("Services", "Create", "err", err.Error())
		}
		b, _ := yaml.Marshal(svc)
		return string(b)
	case repository.Ingress:
		var ing *v1beta1.Ingress
		_ = yaml.Unmarshal([]byte(tpl.FinalTemplate), &ing)
		ing.ObjectMeta.SelfLink = ""
		ing.ObjectMeta.ResourceVersion = ""
		ing.ObjectMeta.Generation = 0
		ing.Namespace = destinationNamespace
		for k, _ := range ing.Spec.Rules {
			ing.Spec.Rules[k].Host = fmt.Sprintf(c.config.GetString("server", "domain_suffix"), appName, destinationNamespace)
		}
		if ing, err = c.k8sClient.Do().ExtensionsV1beta1().Ingresses(destinationNamespace).Create(ing); err != nil {
			_ = level.Error(c.logger).Log("Ingresses", "Create", "err", err.Error())
		}
		b, _ := yaml.Marshal(ing)
		return string(b)
	case repository.ConfigMap:
		var configmap *v1.ConfigMap
		_ = yaml.Unmarshal([]byte(tpl.FinalTemplate), &configmap)
		configmap.ObjectMeta.Namespace = destinationNamespace
		configmap.ObjectMeta.Annotations = nil
		configmap.ObjectMeta.CreationTimestamp = metav1.Time{}
		configmap.ObjectMeta.Generation = 0
		configmap.ObjectMeta.ResourceVersion = ""
		configmap.SelfLink = ""
		if configmap, err = c.k8sClient.Do().CoreV1().ConfigMaps(destinationNamespace).Create(configmap); err != nil {
			_ = level.Error(c.logger).Log("ConfigMaps", "Create", "err", err.Error())
		}
		b, _ := yaml.Marshal(configmap)
		return string(b)
	case repository.VirtualService:
		var out interface{}
		var obj crd.IstioObject
		if err = c.k8sClient.Do().RESTClient().Post().Namespace(destinationNamespace).
			Resource(crd.VirtualServiceProtoSchema.String()).Body(out).Do().Into(obj); err != nil {
			_ = level.Error(c.logger).Log("ConfigMaps", "Create", "err", err.Error())
		}
	}

	return ""
}

func (c *service) FakeTime(ctx context.Context, fakeTime time.Time, method FakeTimeMethod) (err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	var filePath = "/usr/local/lib/libfaketime.so.1"
	var mountPath = "/usr/local/lib/"

	diffTime := fakeTime.Sub(time.Now()).Seconds()
	diffTimeStr := fmt.Sprintf("%fs", diffTime)
	if diffTime > 0 {
		diffTimeStr = "+" + diffTimeStr
	}

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Get(project.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrToolDeploymentK8sGet
	}

	var volumes []v1.Volume
	for _, v := range deployment.Spec.Template.Spec.Volumes {
		if v.Name == "faketime" {
			continue
		}
		volumes = append(volumes, v)
	}

	if method == FakeTimeAdd {
		volumes = append(volumes, v1.Volume{
			Name: "faketime",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/usr/local/lib/faketime",
				},
			},
		})
	}

	deployment.Spec.Template.Spec.Volumes = volumes

	for k, v := range deployment.Spec.Template.Spec.Containers {
		if v.Name != project.Name {
			continue
		}
		var envs []v1.EnvVar
		for _, val := range v.Env {
			if val.Name == "FAKETIME" || val.Name == "LD_PRELOAD" {
				continue
			}
			envs = append(envs, val)
		}

		var volumeMounts []v1.VolumeMount
		for _, val := range v.VolumeMounts {
			if val.Name == "faketime" {
				continue
			}
			volumeMounts = append(volumeMounts, val)
		}
		if method == FakeTimeAdd {
			envs = append(envs, v1.EnvVar{
				Name:  "FAKETIME",
				Value: diffTimeStr,
			}, v1.EnvVar{
				Name:  "LD_PRELOAD",
				Value: filePath,
			})
			volumeMounts = append(volumeMounts, v1.VolumeMount{
				Name:      "faketime",
				MountPath: mountPath,
			})
		}

		deployment.Spec.Template.Spec.Containers[k].Env = envs
		deployment.Spec.Template.Spec.Containers[k].VolumeMounts = volumeMounts
	}

	deployment, err = c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Update(deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return errors.New(ErrToolDeploymentK8sUpdate.Error() + err.Error())
	}

	go func() {
		b, _ := yaml.Marshal(deployment)
		if e := c.repository.ProjectTemplate().UpdateProjectTemplate(&types.ProjectTemplate{
			ProjectID:     project.ID,
			Kind:          repository.Deployment.String(),
			FinalTemplate: string(b),
		}); e != nil {
			_ = level.Error(c.logger).Log("projectTemplateRepository", "UpdateProjectTemplate", "err", err.Error())
		}
	}()

	return
}
