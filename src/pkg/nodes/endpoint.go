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
		Name   string `json:"name"`
		Memory int64  `json:"memory"`
		Cpu    int64  `json:"cpu"`
	}

	listRequest struct {
		page, pageSize int
	}

	infoRequest struct {
		Name string
	}
	infoResult struct {
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
		PodNum           int64             `json:"podNum"`
	}
)

type Endpoints struct {
	SyncEndpoint endpoint.Endpoint
	ListEndpoint endpoint.Endpoint
	InfoEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SyncEndpoint: makeSyncEndpoint(s),
		InfoEndpoint: makeInfoEndpoint(s),
	}

	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range dmw["Info"] {
		eps.InfoEndpoint = m(eps.InfoEndpoint)
	}
	return eps
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
