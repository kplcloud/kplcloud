/**
 * @Time : 8/11/21 4:24 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package nodes

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/middleware"
)

type (
	syncRequest struct {
		ClusterName string `json:"clusterName"  valid:"required"`
	}

	nodeResult struct {
		Name             string `json:"name"`
		Memory           string `json:"memory"`
		Cpu              int64  `json:"cpu"`
		EphemeralStorage string `json:"ephemeralStorage"` // 临时存储单位字节
		InternalIp       string `json:"internalIp"`       // 节点内部IP
		ExternalIp       string `json:"externalIp"`       // 节点外部IP
		KubeletVersion   string `json:"kubeletVersion"`   // kubelet版本
		KubeProxyVersion string `json:"kubeProxyVersion"` // kubeproxy版本
		ContainerVersion string `json:"containerVersion"` // 容器版本
		OsImage          string `json:"osImage"`          // 系统镜像
		Status           string `json:"status"`           // 状态
		Scheduled        bool   `json:"scheduled"`        // 是否调度
		Remark           string `json:"remark"`           // 备注
	}

	listRequest struct {
		page, pageSize int
		query          string
	}

	infoRequest struct {
		Name string
	}
	drainRequest struct {
		Name  string
		Force bool
	}
	infoResult struct {
		ClusterName      string            `json:"clusterName"`
		ClusterAlias     string            `json:"clusterAlias"`
		ClusterStatus    int               `json:"clusterStatus"`
		Name             string            `json:"name"`
		Status           string            `json:"status"`
		CPU              string            `json:"cpu"`
		Memory           string            `json:"memory"`
		UsedCPU          string            `json:"usedCpu"`
		UsedMemory       string            `json:"usedMemory"`
		Bandwidth        string            `json:"bandwidth"`  // 带宽
		SystemDisk       string            `json:"systemDisk"` // 系统盘大小
		InternalIp       string            `json:"internalIp"`
		ExternalIp       string            `json:"externalIp"`
		KubeletVersion   string            `json:"kubeletVersion"`
		KubeProxyVersion string            `json:"kubeProxyVersion"`
		ContainerVersion string            `json:"containerVersion"`
		OsImage          string            `json:"osImage"`
		Scheduled        bool              `json:"scheduled"`
		Remark           string            `json:"remark"`
		Labels           map[string]string `json:"labels"`
		Annotations      map[string]string `json:"annotations"`
		PodNum           int64             `json:"podNum"`
	}
)

type Endpoints struct {
	SyncEndpoint   endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	InfoEndpoint   endpoint.Endpoint
	CordonEndpoint endpoint.Endpoint
	DrainEndpoint  endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SyncEndpoint:   makeSyncEndpoint(s),
		InfoEndpoint:   makeInfoEndpoint(s),
		ListEndpoint:   makeListEndpoint(s),
		CordonEndpoint: makeCordonEndpoint(s),
		DrainEndpoint:  makeDrainEndpoint(s),
		DeleteEndpoint: makeDeleteEndpoint(s),
	}

	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range dmw["Info"] {
		eps.InfoEndpoint = m(eps.InfoEndpoint)
	}
	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range dmw["Cordon"] {
		eps.CordonEndpoint = m(eps.CordonEndpoint)
		eps.DrainEndpoint = m(eps.DrainEndpoint)
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	return eps
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(infoRequest)
		err = s.Delete(ctx, clusterId, req.Name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDrainEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(drainRequest)
		err = s.Drain(ctx, clusterId, req.Name, req.Force)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeCordonEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(infoRequest)
		err = s.Cordon(ctx, clusterId, req.Name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(listRequest)
		res, total, err := s.List(ctx, clusterId, req.query, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(infoRequest)
		res, err := s.Info(ctx, clusterId, req.Name)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterName, ok := ctx.Value(middleware.ContextKeyClusterName).(string)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		err = s.Sync(ctx, clusterName)
		return encode.Response{
			Error: err,
		}, err
	}
}
