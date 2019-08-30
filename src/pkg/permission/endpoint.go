/**
 * @Time : 2019-07-11 17:32
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package permission

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type permissionRequest struct {
	Icon     string `json:"icon"`
	Id       int64  `json:"id"`
	KeyType  string `json:"keyType"`
	Menu     bool   `json:"menu"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Method   string `json:"method"`
	ParentId int64  `json:"parent_id"`
}

type dragPermissionRequest struct {
	DragKey int64 `json:"drag_key"`
	DropKey int64 `json:"drop_key"`
}

type deletePermissionRequest struct {
	Id int64 `json:"id"`
}

func makeMenuEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Menu(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deletePermissionRequest)
		res, err := s.Delete(ctx, req.Id)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(permissionRequest)
		res, err := s.Update(ctx, req.Id, req.Icon, req.KeyType, req.Menu, req.Name, req.Path, req.Method)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(permissionRequest)
		err := s.Post(ctx, req.Name, req.Path, req.Method, req.Icon, req.Menu, req.ParentId)
		return encode.Response{Err: err}, err
	}
}

func makeDragEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dragPermissionRequest)
		res, err := s.Drag(ctx, int64(req.DragKey), int64(req.DropKey))
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.List(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
