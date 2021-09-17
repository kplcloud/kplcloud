/**
 * @Time : 3/5/21 5:33 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package sysuser

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/kplcloud/kplcloud/src/encode"
)

type (
	listRequest struct {
		email          string
		page, pageSize int
	}
	listResult struct {
		Username     string     `json:"username"`
		Email        string     `json:"email"`
		Locked       bool       `json:"locked"`
		WechatOpenId string     `json:"wechatOpenId"`
		LastLogin    *time.Time `json:"lastLogin"`
		CreatedAt    time.Time  `json:"createdAt"`
		UpdatedAt    time.Time  `json:"updatedAt"`
		Roles        []string   `json:"roles"`
		Namespaces   []string   `json:"namespaces"`
	}

	addRequest struct {
		Username      string  `json:"username"`
		Email         string  `json:"email"`
		Locked        bool    `json:"locked"`
		NamespacesIds []int64 `json:"namespaces"`
		RoleIds       []int64 `json:"roles"`
		ClusterIds    []int64 `json:"clusterIds" valid:"required"`
	}
	updateRequest struct {
		UserId     int64   `json:"userId" valid:"required"`
		Username   string  `json:"username" valid:"required"`
		Email      string  `json:"email" valid:"required"`
		Locked     bool    `json:"locked"`
		ClusterIds []int64 `json:"clusters" valid:"required"`
		RoleIds    []int64 `json:"roles" valid:"required"`
	}
)

type Endpoints struct {
	ListEndpoint   endpoint.Endpoint
	AddEndpoint    endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
	LockedEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint:   makeListEndpoint(s),
		AddEndpoint:    makeAddEndpoint(s),
		DeleteEndpoint: nil,
		UpdateEndpoint: makeUpdateEndpoint(s),
		LockedEndpoint: nil,
	}

	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range dmw["Add"] {
		eps.AddEndpoint = m(eps.AddEndpoint)
	}
	for _, m := range dmw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range dmw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}
	for _, m := range dmw["Locked"] {
		eps.LockedEndpoint = m(eps.LockedEndpoint)
	}

	return eps
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateRequest)
		err = s.Update(ctx, req.UserId, req.Username, req.Email, req.Locked, req.ClusterIds, req.RoleIds)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addRequest)
		err = s.Add(ctx, req.Username, req.Email, req.Locked, req.ClusterIds, req.RoleIds)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, total, err := s.List(ctx, req.email, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}
