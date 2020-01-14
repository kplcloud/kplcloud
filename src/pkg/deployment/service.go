/**
 * @Time : 2019-06-27 10:10
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package deployment

import (
	"context"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/hooks"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/configmapyaml"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"strings"
)

var (
	ErrDeploymentK8sGet                      = errors.New("Kubernetes获取错误,请查询是否存在")
	ErrDeploymentK8sUpdate                   = errors.New("Kubernetes更新错误,请联系管理员")
	ErrDeploymentK8sScale                    = errors.New("Kubernetes伸缩错误,请联系管理员")
	ErrDeploymentK8sUnBindPvc                = errors.New("取消挂载持久化存储失败，请联系管理员")
	ErrDeploymentPvcBindNotVary              = errors.New("挂载持久化存储失败，没有修改")
	ErrDeploymentPvcGet                      = errors.New("存储卷声明获取错,可能没有创建")
	ErrDeploymentPortNum                     = errors.New("至少需要留一个端口")
	ErrDeploymentConfigMapCreate             = errors.New("配置字典创建错误")
	ErrDeploymentConfigMapUpdate             = errors.New("配置字典更新错误")
	ErrDeploymentFilebeatTplGet              = errors.New("Filebeat 模版获取错误,可能不存在")
	ErrDeploymentFilebeatTplEncode           = errors.New("Filebeat 模版解析错误,请联系管理员")
	ErrDeploymentProjectTemplateGet          = errors.New("项目Deplooyment获取错误")
	ErrDeploymentInitContainerTemplateGet    = errors.New("初始化容器模版获取错误")
	ErrDeploymentInitContainerTemplateParse  = errors.New("初始化容器模版解析错误")
	ErrDeploymentProxyContainerTemplateGet   = errors.New("代理容器模版获取错误")
	ErrDeploymentProxyContainerTemplateParse = errors.New("代理容器模版解析错误")
	ErrDeploymentServiceMesh                 = errors.New("您没有启用服务网格,请配置ServiceMesh启用状态")
	//ErrDeploymentConfigMapGet      = errors.New("配置字典获取错误")
)

type Service interface {
	// 获取deployment的yaml信息
	GetYaml(ctx context.Context) (res interface{}, err error)

	// 调整command及args
	CommandArgs(ctx context.Context, commands []string, args []string) (err error)

	// 扩容服务，CPU和内存
	Expansion(ctx context.Context, requestCpu, limitCpu, requestMemory, limitMemory string) (err error)

	// 容器伸缩
	Stretch(ctx context.Context, num int) (err error)

	// 获取项目的存储卷
	GetPvc(ctx context.Context, ns, name string) (map[string]interface{}, error)

	// 持久化存储卷挂载到 deployment
	BindPvc(ctx context.Context, ns, name, path, claimName string) (err error)

	// 取消挂载持久化存储卷
	UnBindPvc(ctx context.Context, ns, name, claimName string) (err error)

	// 增加端口
	AddPort(ctx context.Context, ns, name string, req portRequest) (err error)

	// 删除端口
	DelPort(ctx context.Context, ns, name string, portName string, port int32) (err error)

	// 修改日志规则,调整日志采集
	Logging(ctx context.Context, ns, name, pattern, suffix string, paths []string) (err error)

	// 添加探针或调整探针
	Probe(ctx context.Context, ns, name string, req probeRequest) (err error)

	// 服务网络切换
	Mesh(ctx context.Context, ns, name, model string) (err error)

	// 挂载Hosts
	Hosts(ctx context.Context, hosts []string) error

	// 挂载配置文件
	VolumeConfig(ctx context.Context, mountPath, subPath string) error
}

type service struct {
	logger       log.Logger
	k8sClient    kubernetes.K8sClient
	repository   repository.Repository
	hookQueueSvc hooks.ServiceHookQueue
}

func NewService(logger log.Logger,
	k8sClient kubernetes.K8sClient,
	store repository.Repository,
	hookQueueSvc hooks.ServiceHookQueue) Service {
	return &service{
		logger,
		k8sClient,
		store,
		hookQueueSvc,
	}
}

func (c *service) Hosts(ctx context.Context, hosts []string) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	var host struct {
		Host []v1.HostAlias `json:"host"`
	}

	for _, v := range hosts {
		h := strings.Split(v, " ")
		if len(h) < 2 {
			continue
		}
		var hostnames []string
		for k, val := range h {
			if k == 0 {
				continue
			}
			if len(val) < 1 {
				continue
			}
			hostnames = append(hostnames, val)
		}

		host.Host = append(host.Host, v1.HostAlias{
			IP:        h[0],
			Hostnames: hostnames,
		})
	}

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Get(project.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	defer func() {
		if err == nil {
			if tpl, e := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment); e == nil {
				b, _ := yaml.Marshal(deployment)
				tpl.FinalTemplate = string(b)
				if ee := c.repository.ProjectTemplate().UpdateTemplate(tpl); ee != nil {
					_ = level.Error(c.logger).Log("projectTemplateRepository", "UpdateTemplate", "err", err.Error())
				}
			}

		}
	}()

	deployment.Spec.Template.Spec.HostAliases = host.Host
	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.VolumeHosts,
			project.Name, project.Namespace,
			fmt.Sprintf("服务挂载hosts\n 应用: %s.%s", project.Name, project.Name)); err != nil {
			_ = level.Error(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}

func (c *service) GetYaml(ctx context.Context) (res interface{}, err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Get(project.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return nil, ErrDeploymentK8sGet
	}

	return deployment, nil
}

func (c *service) CommandArgs(ctx context.Context, commands []string, args []string) (err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Get(project.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	for key, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name != project.Name {
			continue
		}
		if len(commands) > 0 {
			deployment.Spec.Template.Spec.Containers[key].Command = commands
		} else {
			deployment.Spec.Template.Spec.Containers[key].Command = nil
		}
		if len(args) > 0 {
			deployment.Spec.Template.Spec.Containers[key].Args = args
		} else {
			deployment.Spec.Template.Spec.Containers[key].Args = nil
		}
	}

	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.CommandEvent,
			project.Name, project.Namespace,
			fmt.Sprintf("更新启动参数\n 应用: %s.%s", project.Name, project.Namespace)); err != nil {
			_ = level.Error(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return
}

func (c *service) Expansion(ctx context.Context, requestCpu, limitCpu, requestMemory, limitMemory string) (err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Get(project.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	maxCpu := resource.MustParse(limitCpu)
	reqCpu := resource.MustParse(requestCpu)
	maxMemory := resource.MustParse(limitMemory)
	reqMemory := resource.MustParse(requestMemory)
	resources := v1.ResourceList{}
	limits := v1.ResourceList{}

	if maxMemory.Value() > 0 {
		limits[v1.ResourceMemory] = *resource.NewQuantity(maxMemory.Value(), resource.BinarySI)
	}

	if maxCpu.MilliValue() > 100 && maxCpu.MilliValue() < 1000 {
		limits[v1.ResourceCPU] = *resource.NewMilliQuantity(maxCpu.MilliValue(), resource.BinarySI)
	} else if maxCpu.Value() > 0 {
		limits[v1.ResourceCPU] = *resource.NewQuantity(maxCpu.Value(), resource.BinarySI)
	}

	if reqCpu.MilliValue() > 100 && reqCpu.MilliValue() < 1000 {
		resources[v1.ResourceCPU] = *resource.NewMilliQuantity(reqCpu.MilliValue(), resource.BinarySI)
	} else if reqCpu.Value() > 0 {
		resources[v1.ResourceCPU] = *resource.NewQuantity(reqCpu.Value(), resource.BinarySI)
	}

	resources[v1.ResourceMemory] = *resource.NewQuantity(reqMemory.Value(), resource.BinarySI)

	for key, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name != project.Name {
			continue
		}
		container.Resources.Requests = resources
		container.Resources.Limits = limits
		deployment.Spec.Template.Spec.Containers[key] = container
		break
	}

	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.Expansion,
			project.Name, project.Namespace,
			fmt.Sprintf("服务扩容\n 应用: %s.%s --> 内存最小值: %v  内存最大值: %v CPU最小值: %v CPU最大值: %v", project.Name, project.Name, requestMemory, limitMemory, requestCpu, limitCpu)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}

func (c *service) Stretch(ctx context.Context, num int) (err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Get(project.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	replace := int32(num)
	deployment.Spec.Replicas = &replace
	b, _ := json.Marshal(deployment)

	deployment, err = c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Patch(project.Name, apitypes.MergePatchType, b, "scale")
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Patch", "err", err.Error())
		return ErrDeploymentK8sScale
	}
	//@todo update ProjectTemplate
	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.Extend,
			project.Name, project.Namespace,
			fmt.Sprintf("服务伸缩\n 应用: %v.%v --> 副本数: %d", project.Name, project.Namespace, num)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}

func (c *service) GetPvc(ctx context.Context, ns, name string) (res map[string]interface{}, err error) {
	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return nil, ErrDeploymentK8sGet
	}

	// 暂时不考虑多个挂载的情况
	var volumeName, pvcName, mountPath string
	for _, volume := range deployment.Spec.Template.Spec.Volumes {
		if volume.PersistentVolumeClaim == nil {
			continue
		}
		volumeName = volume.Name
		pvcName = volume.PersistentVolumeClaim.ClaimName
		break
	}

	type pvcChannel struct {
		Pvc   chan *v1.PersistentVolumeClaim
		Error chan error
	}

	type pvChannel struct {
		Pv    chan *v1.PersistentVolume
		Error chan error
	}

	channel := pvcChannel{
		Pvc:   make(chan *v1.PersistentVolumeClaim, 1),
		Error: make(chan error, 1),
	}

	pvCh := pvChannel{
		Pv:    make(chan *v1.PersistentVolume, 1),
		Error: make(chan error, 1),
	}

	if pvcName != "" {
		go func() {
			pvc, err := c.k8sClient.Do().CoreV1().PersistentVolumeClaims(ns).Get(pvcName, metav1.GetOptions{})
			channel.Pvc <- pvc
			channel.Error <- err
		}()
	} else {
		channel.Pvc <- nil
		channel.Error <- nil
	}

	for _, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name != name {
			continue
		}
		for _, volume := range container.VolumeMounts {
			if volume.Name != name+"-pvc" {
				continue
			}
			mountPath = volume.MountPath
			break
		}
		break
	}

	if err = <-channel.Error; err != nil {
		_ = level.Warn(c.logger).Log("channel.Error", "<-", "err", err.Error())
	}

	pvc := <-channel.Pvc

	if pvc != nil {
		go func() {
			volume, err := c.k8sClient.Do().CoreV1().PersistentVolumes().Get(pvc.Spec.VolumeName, metav1.GetOptions{})
			pvCh.Pv <- volume
			pvCh.Error <- err
		}()
	} else {
		pvCh.Pv <- nil
		pvCh.Error <- nil
	}

	return map[string]interface{}{
		"volumeName": volumeName,
		"pvcName":    pvcName,
		"volumePath": mountPath,
		"pvc":        pvc,
		"pv":         <-pvCh.Pv,
	}, nil
}

func (c *service) BindPvc(ctx context.Context, ns, name, path, claimName string) (err error) {
	_, err = c.repository.Pvc().Find(ns, claimName)
	if err != nil {
		_ = level.Error(c.logger).Log("pvcRepository", "Find", "err", err.Error())
		return ErrDeploymentPvcGet
	}

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	var pvc *v1.Volume
	index := -1
	for k, volume := range deployment.Spec.Template.Spec.Volumes {
		if volume.PersistentVolumeClaim == nil {
			continue
		}
		pvc = &volume
		index = k
		break
	}

	containerIndex := -1
	volumeIndex := -1
	for key, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name != name {
			continue
		}
		containerIndex = key
		for k, v := range container.VolumeMounts {
			if v.Name != name+"-pvc" {
				continue
			}
			if v.MountPath == path {
				_ = level.Error(c.logger).Log("v.MountPath", v.MountPath, "path", path)
				return ErrDeploymentPvcBindNotVary
			}
			volumeIndex = k
			break
		}
		break
	}

	if pvc != nil {
		pvc.PersistentVolumeClaim.ClaimName = claimName
		pvc.Name = name + "-pvc"
	} else {
		pvc = &v1.Volume{
			Name: name + "-pvc",
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: claimName,
				},
			},
		}
	}
	if index != -1 {
		deployment.Spec.Template.Spec.Volumes[index] = *pvc
	} else {
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, *pvc)
	}

	if containerIndex != -1 {
		volumeMount := v1.VolumeMount{
			Name:      pvc.Name,
			MountPath: path,
		}
		if volumeIndex != -1 {
			deployment.Spec.Template.Spec.Containers[containerIndex].VolumeMounts[volumeIndex] = volumeMount
		} else {
			deployment.Spec.Template.Spec.Containers[containerIndex].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[containerIndex].VolumeMounts, volumeMount)
		}
	}

	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(ns).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sScale
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.StorageEvent,
			name, ns,
			fmt.Sprintf("绑定持久化存储\n 应用: %s.%s", name, ns)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return
}

func (c *service) UnBindPvc(ctx context.Context, ns, name, claimName string) (err error) {
	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	var pvcName string
	{
		var volumes []v1.Volume
		for _, volume := range deployment.Spec.Template.Spec.Volumes {
			if volume.PersistentVolumeClaim == nil {
				volumes = append(volumes, volume)
				continue
			}
			if volume.PersistentVolumeClaim.ClaimName != claimName {
				volumes = append(volumes, volume)
				continue
			}
			pvcName = volume.Name
			break
		}
		deployment.Spec.Template.Spec.Volumes = volumes
	}

	{
		for key, container := range deployment.Spec.Template.Spec.Containers {
			if container.Name != name {
				continue
			}
			var volumes []v1.VolumeMount
			for _, v := range container.VolumeMounts {
				if v.Name != pvcName {
					volumes = append(volumes, v)
					continue
				}
			}
			deployment.Spec.Template.Spec.Containers[key].VolumeMounts = volumes
			break
		}
	}

	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(ns).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUnBindPvc
	}
	//@todo update ProjectTemplate

	return nil
}

func (c *service) AddPort(ctx context.Context, ns, name string, req portRequest) (err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	svc, err := c.k8sClient.Do().CoreV1().Services(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Services", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	var ports []v1.ServicePort
	var containerPorts []v1.ContainerPort
	for _, v := range req.Ports {
		var protocol v1.Protocol
		v.Protocol = strings.ToUpper(v.Protocol)
		if v.Protocol == "UDP" {
			protocol = v1.ProtocolUDP
		}
		ports = append(ports, v1.ServicePort{Name: v.Name, Protocol: protocol, Port: v.Port})
		containerPorts = append(containerPorts, v1.ContainerPort{Name: v.Name, ContainerPort: v.Port, Protocol: protocol})
	}
	svc.Spec.Ports = ports

	defer func() {
		if err == nil {
			var e error
			if svc, e = c.k8sClient.Do().CoreV1().Services(ns).Update(svc); e != nil {
				_ = level.Error(c.logger).Log("Services", "Update", "err", e.Error())
			}
			fields, _ := json.Marshal(req)
			svcTpl, _ := yaml.Marshal(svc)
			if e = c.repository.ProjectTemplate().UpdateFieldsByNsProjectId(project.ID, repository.Service, string(fields), string(svcTpl)); e != nil {
				_ = level.Error(c.logger).Log("projectTemplateRepository", "UpdateFieldsByNsProjectId", "err", e.Error())
			}

		}
	}()

	for k, v := range deployment.Spec.Template.Spec.Containers {
		if v.Name != name {
			continue
		}
		deployment.Spec.Template.Spec.Containers[k].Ports = containerPorts
		break
	}

	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(ns).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	return
}

func (c *service) DelPort(ctx context.Context, ns, name string, portName string, port int32) (err error) {
	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	svc, err := c.k8sClient.Do().CoreV1().Services(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Services", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	var ports []v1.ServicePort
	for _, v := range svc.Spec.Ports {
		if v.Port == port || v.Name == name {
			continue
		}
		ports = append(ports, v)
	}

	if len(ports) < 1 {
		// 必须留一个端口？
		return ErrDeploymentPortNum
	}

	svc.Spec.Ports = ports
	if svc, err = c.k8sClient.Do().CoreV1().Services(ns).Update(svc); err != nil {
		_ = level.Error(c.logger).Log("Services", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	for k, v := range deployment.Spec.Template.Spec.Containers {
		if v.Name != name {
			continue
		}
		var containerPorts []v1.ContainerPort
		for _, val := range v.Ports {
			if val.ContainerPort == port || val.Name == name {
				continue
			}
			containerPorts = append(containerPorts, val)
		}
		deployment.Spec.Template.Spec.Containers[k].Ports = containerPorts
	}

	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(ns).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	//@todo update ProjectTemplate
	return
}

func (c *service) Logging(ctx context.Context, ns, name, pattern, suffix string, paths []string) (err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	var filePaths []string

	for k, v := range paths {
		paths[k] = strings.TrimRight(v, "/")
		paths[k] += "/"
		filePaths = append(filePaths, paths[k]+suffix)
	}

	tpl, err := c.repository.Template().FindByKindType(repository.FilebeatConfigMap)
	if err != nil {
		_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
		return ErrDeploymentFilebeatTplGet
	}

	filebeat, err := encode.EncodeTemplate(repository.FilebeatConfigMap.ToString(), tpl.Detail, map[string]interface{}{
		"name":      name,
		"namespace": ns,
		"pattern":   pattern,
		"paths":     filePaths,
	})
	if err != nil {
		_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
		return ErrDeploymentFilebeatTplEncode
	}

	configMap, err := c.k8sClient.Do().CoreV1().ConfigMaps(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		configMap = new(v1.ConfigMap)
		configMap.APIVersion = "v1"
		configMap.Kind = "ConfigMap"
		configMap.Name = name
		configMap.Namespace = ns
		configMap.Data = map[string]string{}
	}

	configMap.Data["filebeat.yml"] = filebeat
	if configMap.ResourceVersion == "" {
		if configMap, err = c.k8sClient.Do().CoreV1().ConfigMaps(ns).Create(configMap); err != nil {
			_ = level.Error(c.logger).Log("ConfigMaps", "Create", "err", err.Error())
			return ErrDeploymentConfigMapCreate
		}
	} else {
		if configMap, err = c.k8sClient.Do().CoreV1().ConfigMaps(ns).Update(configMap); err != nil {
			_ = level.Error(c.logger).Log("ConfigMaps", "Update", "err", err.Error())
			return ErrDeploymentConfigMapUpdate
		}
	}

	//同步远程configMapYaml 数据
	go configmapyaml.SyncConfigMapYaml(ns, name, c.logger, c.k8sClient, c.repository)

	cmp, _ := yaml.Marshal(configMap)
	if e := c.repository.ProjectTemplate().CreateOrUpdate(&types.ProjectTemplate{
		FinalTemplate: string(cmp),
		Kind:          repository.ConfigMap.String(),
		ProjectID:     project.ID,
	}); e != nil {
		_ = level.Warn(c.logger).Log("projectTemplateRepository", "CreateOrUpdate", "err", e.Error())
	}

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	var filebeatExists bool
	for _, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == "filebeat" {
			filebeatExists = true
			break
		}
	}

	if !filebeatExists {
		tpl, err = c.repository.Template().FindByKindType(repository.FilebeatContainer)
		if err != nil {
			_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
			return ErrDeploymentFilebeatTplGet
		}

		if filebeat, err = encode.EncodeTemplate(repository.FilebeatContainer.ToString(), tpl.Detail, map[string]interface{}{
			"name": name,
		}); err != nil {
			_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
			return ErrDeploymentFilebeatTplEncode
		}

		var filebeatContainer v1.Container
		if err = yaml.Unmarshal([]byte(filebeat), &filebeatContainer); err != nil {
			_ = level.Error(c.logger).Log("yaml", "Unmarshal", "err", err.Error())
			return ErrDeploymentFilebeatTplEncode
		}
		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, filebeatContainer)
	}

	// 暂时只考虑一个路径
	var volumes []v1.Volume
	for _, volume := range deployment.Spec.Template.Spec.Volumes {
		if volume.Name == "app-logs" || (volume.Name == name && volume.ConfigMap != nil) {
			continue
		}
		volumes = append(volumes, volume)
	}
	volumes = append(volumes, v1.Volume{
		Name: "app-logs",
	}, v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name,
				},
			},
		},
	})
	deployment.Spec.Template.Spec.Volumes = volumes

	for key, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name != name {
			continue
		}
		var volumeMounts []v1.VolumeMount
		for _, volumeMount := range container.VolumeMounts {
			if volumeMount.Name == "app-logs" {
				continue
			}
			volumeMounts = append(volumeMounts, volumeMount)
		}
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      "app-logs",
			MountPath: paths[0],
		})
		deployment.Spec.Template.Spec.Containers[key].VolumeMounts = volumeMounts
	}

	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(ns).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.LoggingEvent,
			name, ns,
			fmt.Sprintf("调整日志采集器: %s.%s", name, ns)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return
}

func (c *service) Probe(ctx context.Context, ns, name string, req probeRequest) (err error) {

	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	for key, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name != name {
			continue
		}
		var readiness, liveness *v1.Probe
		for _, val := range req.Probe {
			if val == ProbeReadiness.String() {
				readiness = &v1.Probe{
					Handler: v1.Handler{
						TCPSocket: &v1.TCPSocketAction{
							Port: intstr.IntOrString{IntVal: req.Port},
						},
					},
					InitialDelaySeconds: req.InitialDelaySeconds,
					PeriodSeconds:       req.PeriodSeconds,
					TimeoutSeconds:      req.TimeoutSeconds,
					SuccessThreshold:    req.SuccessThreshold,
					FailureThreshold:    req.FailureThreshold,
				}
			} else if val == ProbeLiveness.String() {
				liveness = &v1.Probe{
					Handler: v1.Handler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.IntOrString{IntVal: req.Port},
							Path: req.Path,
						},
					},
					InitialDelaySeconds: req.InitialDelaySeconds,
					PeriodSeconds:       req.PeriodSeconds,
					TimeoutSeconds:      req.TimeoutSeconds,
					SuccessThreshold:    req.SuccessThreshold,
					FailureThreshold:    req.FailureThreshold,
				}
			}
		}

		deployment.Spec.Template.Spec.Containers[key].LivenessProbe = liveness
		deployment.Spec.Template.Spec.Containers[key].ReadinessProbe = readiness
	}

	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(ns).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	if tpl, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment); err == nil && tpl != nil {
		b, _ := yaml.Marshal(deployment)
		tpl.FinalTemplate = string(b)
		if e := c.repository.ProjectTemplate().UpdateTemplate(tpl); e != nil {
			_ = level.Error(c.logger).Log("ProjectTemplate", "UpdateTemplate", "err", e.Error())
		}
	}

	//@todo update ProjectTemplate
	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.ReadinessProbe,
			name, ns,
			fmt.Sprintf("修改探针\n 应用: %s.%s --> 端口: %d", name, ns, req.Port)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()
	return
}

func (c *service) Mesh(ctx context.Context, ns, name, model string) (err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	projectTemplate, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindByProjectId", "err", err.Error())
		return ErrDeploymentProjectTemplateGet
	}

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	if model == repository.FieldMesh {
		// todo 如果没有启用istio的话无法使用servicesmesh模式

		// 初始化相关容器
		initContainerTpl, err := c.repository.Template().FindByKindType(repository.InitContainersKind)
		if err != nil {
			_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
			return ErrDeploymentInitContainerTemplateGet
		}
		proxyContainerTpl, err := c.repository.Template().FindByKindType(repository.IstioProxyKind)
		if err != nil {
			_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
			return ErrDeploymentProxyContainerTemplateGet
		}
		if deployment.Spec.Template.Annotations == nil {
			annotations := map[string]string{
				"sidecar.istio.io/status": `{"version":"` + encode.HashString([]byte(ns+name)) + `","initContainers":["istio-init","enable-core-dump"],"containers":["istio-proxy"],"volumes":["istio-envoy","istio-certs"]}`,
			}
			deployment.Spec.Template.Annotations = annotations
		} else {
			deployment.Spec.Template.Annotations["sidecar.istio.io/status"] = `{"version":"` + encode.HashString([]byte(ns+name)) + `","initContainers":["istio-init","enable-core-dump"],"containers":["istio-proxy"],"volumes":["istio-envoy","istio-certs"]}`
		}
		fields := map[string]interface{}{
			"name":     name,
			"language": project.Language,
		}
		var ports []map[string]interface{}
		var initPorts []string
		for _, v := range deployment.Spec.Template.Spec.Containers {
			for _, containerPort := range v.Ports {
				ports = append(ports, map[string]interface{}{
					"port": containerPort.ContainerPort,
				})
				initPorts = append(initPorts, encode.String(containerPort.ContainerPort))
			}
		}
		fields["ports"] = ports
		fields["initPorts"] = strings.Join(initPorts, ",")

		initContainerYml, err := encode.EncodeTemplate(repository.InitContainersKind.ToString(), initContainerTpl.Detail, fields)
		if err != nil {
			_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
			return ErrDeploymentInitContainerTemplateParse
		}
		proxyContainerYml, err := encode.EncodeTemplate(repository.IstioProxyKind.ToString(), proxyContainerTpl.Detail, fields)
		if err != nil {
			_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
			return ErrDeploymentProxyContainerTemplateParse
		}
		var container v1.Container
		var initContainer v1.PodSpec
		if err = yaml.Unmarshal([]byte(initContainerYml), &initContainer); err != nil {
			_ = level.Error(c.logger).Log("yaml", "Unmarshal", "err", err.Error())
			return ErrDeploymentInitContainerTemplateParse
		}
		if err = yaml.Unmarshal([]byte(proxyContainerYml), &container); err != nil {
			_ = level.Error(c.logger).Log("yaml", "Unmarshal", "err", err.Error())
			return ErrDeploymentProxyContainerTemplateParse
		}
		deployment.Spec.Template.Spec.InitContainers = append(deployment.Spec.Template.Spec.InitContainers, initContainer.InitContainers[0])
		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, container)

		// 相关挂载
		optional := true
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, v1.Volume{
			Name: "istio-envoy",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{
					Medium: "Memory",
				},
			},
		}, v1.Volume{
			Name: "istio-certs",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: "istio.default",
					Optional:   &optional,
				},
			},
		})
		projectTemplate.FieldStruct.Mesh = repository.FieldMesh
	} else {
		delete(deployment.Spec.Template.Annotations, "sidecar.istio.io/status")
		// 去掉 istio-proxy
		var containers []v1.Container
		for _, v := range deployment.Spec.Template.Spec.Containers {
			if v.Name == "istio-proxy" {
				continue
			}
			containers = append(containers, v)
		}
		deployment.Spec.Template.Spec.Containers = containers

		// 去掉初始容器
		var initContainers []v1.Container
		for _, v := range deployment.Spec.Template.Spec.InitContainers {
			if v.Name == "istio-init" || v.Name == "enable-core-dump" {
				continue
			}
			initContainers = append(initContainers, v)
		}
		deployment.Spec.Template.Spec.InitContainers = initContainers

		// 去除 volumes
		var volumes []v1.Volume
		for _, v := range deployment.Spec.Template.Spec.Volumes {
			if v.Name == "istio-envoy" || v.Name == "istio-certs" {
				continue
			}
			volumes = append(volumes, v)
		}
		deployment.Spec.Template.Spec.Volumes = volumes
		projectTemplate.FieldStruct.Mesh = repository.FieldNormal
	}

	if deployment, err = c.k8sClient.Do().AppsV1().Deployments(ns).Update(deployment); err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	if err = c.repository.ProjectTemplate().UpdateTemplate(projectTemplate); err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "UpdateTemplate", "err", err.Error())
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.SwitchModel,
			name, ns,
			fmt.Sprintf("修改服务模式\n 应用: %s.%s --> 模式类型: %d", name, ns, model)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}

func (c *service) VolumeConfig(ctx context.Context, mountPath, subPath string) error {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	name := ctx.Value(middleware.NameContext).(string)

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrDeploymentK8sGet
	}

	for key, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name != name {
			continue
		}

		volumeMounts := deployment.Spec.Template.Spec.Containers[key].VolumeMounts
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      name + "-config",
			ReadOnly:  true,
			MountPath: mountPath,
			SubPath:   subPath,
		})
		deployment.Spec.Template.Spec.Containers[key].VolumeMounts = volumeMounts
	}

	deployment, err = c.k8sClient.Do().AppsV1().Deployments(ns).Update(deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
		return ErrDeploymentK8sUpdate
	}

	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.VolumeConfig,
			name, ns,
			fmt.Sprintf("增加挂载配置文件\n 应用: %s.%s --> 文件名: %s \n 路径: %s\n", name, ns, subPath, mountPath)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return nil
}
