package discovery

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
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/convert"
	"github.com/kplcloud/kplcloud/src/util/encode"
	utilpods "github.com/kplcloud/kplcloud/src/util/pods"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
)

var (
	ErrServiceNamePattern = errors.New("名称格式不符合要求,请重新填写")
	ErrServiceK8sDelete   = errors.New("删除Service资源错误,请查询是否存在")
	ErrServiceK8sList     = errors.New("获取Service资源列表错误,请查询是否存在")
	ErrServiceK8sGet      = errors.New("获取Service资源错误,请查询是否存在")
	ErrServiceTplGet      = errors.New("获取Service模版错误,请联系管理员")
	ErrServiceTplParse    = errors.New("解析Service模版错误,请联系管理员")
	ErrServiceK8sCreate   = errors.New("创建Service错误,可能已经存在")
	ErrServiceK8sUpdate   = errors.New("更新Service错误,可能已经存在")
	ErrEndpointsTplGet    = errors.New("获取Endpoints模版错误,请联系管理员")
	ErrEndpointsParse     = errors.New("创建Endpoints错误,请联系管理员")
	ErrEndpointsK8sGet    = errors.New("获取Endpoints错误,可能不存在")
	ErrEndpointsK8sCreate = errors.New("创建Endpoints错误,请联系管理员")
	ErrProjectGet         = errors.New("项目可能不存在,无法与之关联")
)

type ResourceType string

const (
	ServiceResourceType  ResourceType = "service"
	EndpointResourceType ResourceType = "endpoint"
)

type serviceList struct {
	Name             string                 `json:"name"`
	Labels           map[string]string      `json:"labels"`
	ClusterIP        string                 `json:"cluster_ip"`
	InsideEndpoint   []v1.ServicePort       `json:"inside_endpoint"`
	ExternalEndpoint map[string]interface{} `json:"external_endpoint"`
	CreatedAt        metav1.Time            `json:"created_at"`
	Namespace        string                 `json:"namespace"`
}

type Service interface {
	// 删除service
	Delete(ctx context.Context, svcName string) error

	// 查看service详情
	Detail(ctx context.Context, svcName string) (map[string]interface{}, error)

	// Service 列表页
	List(ctx context.Context, page, limit int) (res []*serviceList, err error)

	// 创建Service
	Create(ctx context.Context, req createRequest) (err error)

	// 更新service
	Update(ctx context.Context, req createRequest) (err error)

	// 创建service
	// Deprecated: 这个功能可以去掉
	PostYaml(ctx context.Context, body []byte) error
}

type service struct {
	logger     log.Logger
	k8sClient  kubernetes.K8sClient
	config     *config.Config
	repository repository.Repository
}

func NewService(logger log.Logger, k8sClient kubernetes.K8sClient,
	config *config.Config, store repository.Repository) Service {
	return &service{logger,
		k8sClient,
		config,
		store,
	}
}

func (c *service) Update(ctx context.Context, req createRequest) (err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)

	svc, err := c.k8sClient.Do().CoreV1().Services(ns).Get(req.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Services", "Get", "err", err.Error())
		return ErrServiceK8sGet
	}

	var ports []map[string]interface{}
	for _, v := range req.Routes {
		if v.Name != "" && !convert.IsEnNameString(v.Name) {
			return ErrServiceNamePattern
		}
		ports = append(ports, map[string]interface{}{
			"port":       v.Port,
			"protocol":   v.Protocol,
			"targetport": v.TargetPort,
			"name":       v.Name,
		})
	}

	serviceTpl, err := c.repository.Template().FindByKindType(repository.ServiceKind)
	if err != nil {
		_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
		return ErrServiceTplGet
	}
	svcYaml, err := encode.EncodeTemplate(repository.ServiceKind.ToString(), serviceTpl.Detail, map[string]interface{}{
		"name":         req.Name,
		"namespace":    ns,
		"ports":        ports,
		"resourceType": req.ResourceType,
	})

	_ = yaml.Unmarshal([]byte(svcYaml), &svc)

	if err != nil {
		_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
		return ErrServiceTplParse
	}

	if svc, err = c.k8sClient.Do().CoreV1().Services(ns).Update(svc); err != nil {
		_ = level.Error(c.logger).Log("Services", "Update", "err", err.Error())
		return ErrServiceK8sUpdate
	}

	if ResourceType(req.ResourceType) == ServiceResourceType {
		project, err := c.repository.Project().FindByNsName(ns, req.ServiceProject)
		if err != nil {
			_ = level.Error(c.logger).Log("projectRepository", "FindByNsName", "err", err.Error())
			return ErrProjectGet
		}

		if project != nil {
			// 创建projectTemplate
			fields, _ := json.Marshal(req.Routes)
			finalTpl, _ := yaml.Marshal(svc)
			go func() {
				if svcTpl, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Service); err == nil {
					svcTpl.FinalTemplate = string(finalTpl)
					svcTpl.Fields = string(fields)
					if err = c.repository.ProjectTemplate().UpdateTemplate(svcTpl); err != nil {
						_ = level.Error(c.logger).Log("projectTemplate", "FirstOrCreate", "err", err.Error())
					}
				}
			}()
		}
	}

	if ResourceType(req.ResourceType) == EndpointResourceType {
		// 如果选择的是 端点的话
		endpointTpl, err := c.repository.Template().FindByKindType(repository.EndpointsKind)
		if err != nil {
			_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
			return ErrEndpointsTplGet
		}
		var addressPorts []map[string]interface{}
		for _, v := range req.Address.Ports {
			addressPorts = append(addressPorts, map[string]interface{}{
				"port": v.Port,
				"name": v.Name,
			})
		}
		var address []map[string]interface{}
		address = append(address, map[string]interface{}{
			"ips":   req.Address.Ips,
			"ports": addressPorts,
		})
		epYaml, err := encode.EncodeTemplate(repository.EndpointsKind.ToString(), endpointTpl.Detail, map[string]interface{}{
			"name":      req.Name,
			"namespace": ns,
			"ports":     addressPorts,
			"addresses": address,
		})
		if err != nil {
			_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
			return ErrEndpointsParse
		}
		ep, err := c.k8sClient.Do().CoreV1().Endpoints(ns).Get(req.Name, metav1.GetOptions{})
		if err != nil {
			_ = level.Error(c.logger).Log("endpoints", "Get", "err", err.Error())
			return ErrEndpointsK8sGet
		}
		_ = yaml.Unmarshal([]byte(epYaml), &ep)
		if ep, err = c.k8sClient.Do().CoreV1().Endpoints(ns).Update(ep); err != nil {
			_ = level.Error(c.logger).Log("Endpoints", "Create", "err", err.Error())
			return ErrEndpointsK8sCreate
		}
	}

	return nil
}

func (c *service) Create(ctx context.Context, req createRequest) (err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)

	var ports []map[string]interface{}
	for _, v := range req.Routes {
		if v.Name != "" && !convert.IsEnNameString(v.Name) {
			return ErrServiceNamePattern
		}
		ports = append(ports, map[string]interface{}{
			"port":       v.Port,
			"protocol":   v.Protocol,
			"targetport": v.TargetPort,
			"name":       v.Name,
		})
	}

	serviceTpl, err := c.repository.Template().FindByKindType(repository.ServiceKind)
	if err != nil {
		_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
		return ErrServiceTplGet
	}
	svcYaml, err := encode.EncodeTemplate(repository.ServiceKind.ToString(), serviceTpl.Detail, map[string]interface{}{
		"name":         req.Name,
		"namespace":    ns,
		"ports":        ports,
		"resourceType": req.ResourceType,
	})
	if err != nil {
		_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
		return ErrServiceTplParse
	}

	var svc *v1.Service
	_ = yaml.Unmarshal([]byte(svcYaml), &svc)

	svc, err = c.k8sClient.Do().CoreV1().Services(ns).Create(svc)
	if err != nil {
		_ = level.Error(c.logger).Log("Services", "Create", "err", err.Error())
		return ErrServiceK8sCreate
	}

	if ResourceType(req.ResourceType) == ServiceResourceType {
		project, err := c.repository.Project().FindByNsName(ns, req.ServiceProject)
		if err != nil {
			_ = level.Error(c.logger).Log("projectRepository", "FindByNsName", "err", err.Error())
			return ErrProjectGet
		}

		if project != nil {
			// 创建projectTemplate
			fields, _ := json.Marshal(req.Routes)
			finalTpl, _ := yaml.Marshal(svc)
			go func() {
				_, err = c.repository.ProjectTemplate().FirstOrCreate(project.ID, repository.Service, string(fields), string(finalTpl), 1)
				if err != nil {
					_ = level.Error(c.logger).Log("projectTemplate", "FirstOrCreate", "err", err.Error())
				}
			}()
		}
	}

	if ResourceType(req.ResourceType) == EndpointResourceType {
		// 如果选择的是 端点的话
		endpointTpl, err := c.repository.Template().FindByKindType(repository.EndpointsKind)
		if err != nil {
			_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
			return ErrEndpointsTplGet
		}
		var addressPorts []map[string]interface{}
		for _, v := range req.Address.Ports {
			addressPorts = append(addressPorts, map[string]interface{}{
				"port": v.Port,
				"name": v.Name,
			})
		}
		var address []map[string]interface{}
		address = append(address, map[string]interface{}{
			"ips":   req.Address.Ips,
			"ports": addressPorts,
		})
		epYaml, err := encode.EncodeTemplate(repository.EndpointsKind.ToString(), endpointTpl.Detail, map[string]interface{}{
			"name":      req.Name,
			"namespace": ns,
			"ports":     addressPorts,
			"addresses": address,
		})
		if err != nil {
			_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
			return ErrEndpointsParse
		}
		var ep *v1.Endpoints
		_ = yaml.Unmarshal([]byte(epYaml), &ep)
		if ep, err = c.k8sClient.Do().CoreV1().Endpoints(ns).Create(ep); err != nil {
			_ = level.Error(c.logger).Log("Endpoints", "Create", "err", err.Error())
			return ErrEndpointsK8sCreate
		}
	}

	return nil
}

func (c *service) PostYaml(ctx context.Context, body []byte) error {
	ns := ctx.Value(middleware.NamespaceContext).(string)

	var err error
	var svc *v1.Service
	if err := yaml.Unmarshal(body, &svc); err != nil {
		return err
	}

	if svc, err = c.k8sClient.Do().CoreV1().Services(ns).Create(svc); err != nil {
		return err
	}

	if project, err := c.repository.Project().FindByNsName(ns, svc.Name); err == nil && project.ID != 0 {
		// projectTemplate增加模版
		fields, _ := json.Marshal(svc.Spec.Ports)
		final, _ := yaml.Marshal(svc)

		// todo templateId 没什么用，可以去掉
		if _, err = c.repository.ProjectTemplate().FirstOrCreate(project.ID, repository.Service, string(fields), string(final), 1); err != nil {
			_ = level.Error(c.logger).Log("projectTemplate", "FirstOrCreate", "err", err.Error())
		}
	}

	return nil
}

func (c *service) List(ctx context.Context, page, limit int) (res []*serviceList, err error) {
	// todo 好像还有一个搜索功能没做
	ns := ctx.Value(middleware.NamespaceContext).(string)
	page = (page - 1) * limit

	svcList, err := c.k8sClient.Do().CoreV1().Services(ns).List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Services", "List", "err", err.Error())
		return nil, ErrServiceK8sList
	}
	length := len(svcList.Items)
	if length == 0 || length < page {
		return nil, nil
	}
	var lenPageSize = page + limit
	if lenPageSize > length {
		lenPageSize = length
	}
	var sourceList svcSort
	sourceList = svcList.Items
	sort.Sort(sourceList)

	//list := sourceList[page:lenPageSize]
	list := sourceList
	for _, svc := range list {
		res = append(res, &serviceList{
			Name:           svc.Name,
			Labels:         svc.Labels,
			ClusterIP:      svc.Spec.ClusterIP,
			InsideEndpoint: svc.Spec.Ports,
			CreatedAt:      svc.CreationTimestamp,
			ExternalEndpoint: map[string]interface{}{
				"ExternalIPs":           svc.Spec.ExternalIPs,
				"ExternalName":          svc.Spec.ExternalName,
				"ExternalTrafficPolicy": svc.Spec.ExternalTrafficPolicy,
			},
			Namespace: ns,
		})
	}

	return res, nil
}

func (c *service) Detail(ctx context.Context, svcName string) (map[string]interface{}, error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	svc, err := c.k8sClient.Do().CoreV1().Services(ns).Get(svcName, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Services", "Get", "err", err.Error())
		return nil, ErrServiceK8sGet
	}
	var pods []map[string]interface{}

	var endpointList []v1.EndpointSubset
	{
		if endpoints, err := c.k8sClient.Do().CoreV1().Endpoints(ns).Get(svcName, metav1.GetOptions{}); err == nil {
			for _, subset := range endpoints.Subsets {
				endpointList = append(endpointList, subset)
			}
		} else {
			_ = level.Error(c.logger).Log("Endpoints", "Get", "err", err.Error())
		}
	}
	var selectorKey, selectorVal string
	for key, val := range svc.Spec.Selector {
		selectorKey = key
		selectorVal = val
	}

	if podList, err := c.k8sClient.Do().CoreV1().Pods(ns).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", selectorKey, selectorVal),
	}); err == nil {
		for _, pod := range podList.Items {
			var (
				isReady              = true
				message, lastMessage string
				restartCount         int32
			)
			for _, v := range pod.Status.ContainerStatuses {
				if v.Name == svcName {
					restartCount = v.RestartCount
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
			metrics := make(chan map[string]interface{})
			go utilpods.GetPodsMetrics(pod.Namespace, pod.Name, c.config.GetString("server", "heapster_url"), metrics)
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

			pods = append(pods, map[string]interface{}{
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
			})
		}
	} else {
		_ = level.Error(c.logger).Log("Pods", "List", "err", err.Error())
	}

	var eps v1.EndpointSubset
	if len(endpointList) > 0 {
		eps = endpointList[0]
	}

	return map[string]interface{}{
		"service":   svc,
		"endpoints": eps,
		"pods":      pods,
	}, nil
}

func (c *service) Delete(ctx context.Context, svcName string) error {
	ns := ctx.Value(middleware.NamespaceContext).(string)

	var err error

	defer func() {
		if err == nil {
			if project, err := c.repository.Project().FindByNsName(ns, svcName); err == nil && project.ID != 0 {
				if err = c.repository.ProjectTemplate().Delete(project.ID, repository.Service); err != nil {
					_ = level.Error(c.logger).Log("projectTemplate", "delete", "err", err)
				}
			}
		}
	}()

	if err = c.k8sClient.Do().CoreV1().Services(ns).Delete(svcName, &metav1.DeleteOptions{}); err != nil {
		_ = level.Error(c.logger).Log("Services", "Delete", "err", err.Error())
		return ErrServiceK8sDelete
	}

	go func() {
		if err := c.k8sClient.Do().CoreV1().Endpoints(ns).Delete(svcName, &metav1.DeleteOptions{}); err != nil {
			_ = level.Error(c.logger).Log("Endpoints", "Delete", "err", err.Error())
		}
	}()

	// todo 操作记录 event history

	return nil
}

type svcSort []v1.Service

func (m svcSort) Len() int {
	return len(m)
}
func (m svcSort) Less(i, j int) bool {
	return m[i].CreationTimestamp.Unix() > m[j].CreationTimestamp.Unix()
}
func (m svcSort) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
