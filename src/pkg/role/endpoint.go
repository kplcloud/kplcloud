/**
 * @Time : 2019-07-16 18:00
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package role

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type getRoleRequest struct {
	Id int64
}

type permissionIds struct {
	PermissionIds []string `json:"permission_ids"`
}

type rolePermissionsRequest struct {
	RoleId        int64   `json:"role_id"`
	PermissionIds []int64 `json:"permission_ids"`
}

type roleRequest struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Desc  string `json:"desc"`
}

func makePermissionSelectedEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRoleRequest)
		res, err := s.PermissionSelected(ctx, req.Id)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(roleRequest)
		err := s.Update(ctx, req.Id, req.Name, req.Desc, req.Level)
		return encode.Response{Err: err}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(roleRequest)
		err := s.Post(ctx, req.Name, req.Desc, req.Level)
		return encode.Response{Err: err}, err
	}
}

func makeAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.All(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRoleRequest)
		res, err := s.Detail(ctx, req.Id)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRoleRequest)
		err := s.Delete(ctx, req.Id)
		return encode.Response{Err: err}, err
	}
}

func makeRolePermissionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(rolePermissionsRequest)
		err := s.RolePermission(ctx, req.RoleId, req.PermissionIds)
		return encode.Response{Err: err}, err
	}
}
