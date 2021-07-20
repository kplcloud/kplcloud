/**
 * @Time : 3/10/21 3:30 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package sysrole

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/kplcloud/kplcloud/src/encode"
)

type (
	listResult struct {
		Id          int64     `json:"id"`
		Alias       string    `json:"alias"`
		Name        string    `json:"name"`
		Enabled     bool      `json:"enabled"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	listRequest struct {
		page, pageSize int
	}

	addRequest struct {
		Id          int64  `json:"id"`
		Alias       string `json:"alias"`
		Name        string `json:"name"`
		Enabled     bool   `json:"enabled"`
		Description string `json:"description"`
	}
	permissionRequest struct {
		Id      int64   `json:"id"`
		PermIds []int64 `json:"permIds"`
	}
	userRequest struct {
		Id      int64   `json:"id"`
		UserIds []int64 `json:"userIds"`
	}
	deleteRequest struct {
		Id int64 `json:"id"`
	}
)

type Endpoints struct {
	ListEndpoint        endpoint.Endpoint
	AddEndpoint         endpoint.Endpoint
	PermissionsEndpoint endpoint.Endpoint
	PermissionEndpoint  endpoint.Endpoint
	//UsersEndpoint       endpoint.Endpoint
	UserEndpoint   endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint:       makeListEndpoint(s),
		AddEndpoint:        makeAddEndpoint(s),
		PermissionEndpoint: makePermissionEndpoint(s),
		UserEndpoint:       makeUserEndpoint(s),
		UpdateEndpoint:     makeUpdateEndpoint(s),
		DeleteEndpoint:     makeDeleteEndpoint(s),
	}

	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range dmw["Add"] {
		eps.AddEndpoint = m(eps.AddEndpoint)
	}
	for _, m := range dmw["Permission"] {
		eps.PermissionEndpoint = m(eps.PermissionEndpoint)
	}
	for _, m := range dmw["User"] {
		eps.UserEndpoint = m(eps.UserEndpoint)
	}
	for _, m := range dmw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}
	for _, m := range dmw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}

	return eps
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteRequest)
		err = s.Delete(ctx, req.Id)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addRequest)
		err = s.Update(ctx, req.Id, req.Alias, req.Name, req.Description, req.Enabled)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(userRequest)
		err = s.User(ctx, req.Id, req.UserIds)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makePermissionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(permissionRequest)
		err = s.Permission(ctx, req.Id, req.PermIds)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, total, err := s.List(ctx, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addRequest)
		err = s.Add(ctx, req.Alias, req.Name, req.Description, req.Enabled)
		return encode.Response{
			Error: err,
		}, err
	}
}
