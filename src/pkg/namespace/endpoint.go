/**
 * @Time : 8/13/21 2:55 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package namespace

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/middleware"
	"time"
)

type (
	createRequest struct {
		Name    string   `json:"name" valid:"required"`
		Alias   string   `json:"alias" valid:"required"`
		Remark  string   `json:"remark"`
		Secrets []string `json:"secrets"`
	}

	result struct {
		Name         string    `json:"name"`
		Alias        string    `json:"alias"`
		Remark       string    `json:"remark"`
		Status       string    `json:"status"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
		ImageSecrets []string  `json:"imageSecrets"` // 空间获取镜像的证书
		RegSecrets   []string  `json:"regSecrets"`   // 所有镜像仓库
		ClusterAlias string    `json:"clusterAlias,omitempty"`
		ClusterName  string    `json:"clusterName,omitempty"`
	}

	listRequest struct {
		query          string
		page, pageSize int
	}

	deleteRequest struct {
		force bool
	}

	issueSecretRequest struct {
		Registry string `json:"registry" valid:"required"`
	}
)

type Endpoints struct {
	SyncEndpoint        endpoint.Endpoint
	CreateEndpoint      endpoint.Endpoint
	ListEndpoint        endpoint.Endpoint
	UpdateEndpoint      endpoint.Endpoint
	DeleteEndpoint      endpoint.Endpoint
	InfoEndpoint        endpoint.Endpoint
	IssueSecretEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SyncEndpoint:        makeSyncEndpoint(s),
		CreateEndpoint:      makeCreateEndpoint(s),
		ListEndpoint:        makeListEndpoint(s),
		UpdateEndpoint:      makeUpdateEndpoint(s),
		DeleteEndpoint:      makeDeleteEndpoint(s),
		InfoEndpoint:        makeInfoEndpoint(s),
		IssueSecretEndpoint: makeIssueSecretEndpoint(s),
	}

	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range dmw["Create"] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
	}
	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
		eps.InfoEndpoint = m(eps.InfoEndpoint)
	}
	for _, m := range dmw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
		eps.IssueSecretEndpoint = m(eps.IssueSecretEndpoint)
	}
	return eps
}

func makeIssueSecretEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, ok := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		req := request.(issueSecretRequest)
		err = s.IssueSecret(ctx, clusterId, ns, req.Registry)
		return encode.Response{
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
		ns, ok := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		res, err := s.Info(ctx, clusterId, ns)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, ok := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		req := request.(deleteRequest)
		err = s.Delete(ctx, clusterId, ns, req.force)
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
		namespaces, ok := ctx.Value(middleware.ContextKeyNamespaceList).([]string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}

		req := request.(listRequest)
		res, total, err := s.List(ctx, clusterId, namespaces, req.query, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, ok := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		req := request.(createRequest)
		err = s.Update(ctx, clusterId, ns, req.Alias, req.Remark, "", req.Secrets)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(createRequest)
		err = s.Create(ctx, clusterId, req.Name, req.Alias, req.Remark, req.Secrets)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		err = s.Sync(ctx, clusterId)
		return encode.Response{
			Error: err,
		}, err
	}
}
