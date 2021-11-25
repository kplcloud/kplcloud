/**
 * @Time : 5/7/21 9:45 AM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package syspermission

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"time"
)

type (
	result struct {
		Name      string    `json:"name"`
		Path      string    `json:"path"`
		Method    string    `json:"method"`
		Alias     string    `json:"alias"`
		Remark    string    `json:"remark"`
		ParentId  int64     `json:"parentId"`
		Id        int64     `json:"id"`
		Menu      bool      `json:"menu"`
		Sort      int       `json:"sort"`
		CreatedAt time.Time `json:"createdAt"`
		Children  []result  `json:"children"`
	}

	addRequest struct {
		Name     string `json:"name" valid:"required"`
		Alias    string `json:"alias" valid:"required"`
		Icon     string `json:"icon"`
		Path     string `json:"path" valid:"required"`
		Method   string `json:"method" valid:"required"`
		Desc     string `json:"desc"`
		ParentId int64  `json:"parentId"`
		Menu     bool   `json:"menu"`
	}
	updateRequest struct {
		Id int64 `json:"id"`
		addRequest
	}

	deleteRequest struct {
		Id int64 `json:"id"`
	}

	dragRequest struct {
		DragId int64 `json:"dragId" valid:"required"`
		DropId int64 `json:"dropId" valid:"required"`
	}
)

type Endpoints struct {
	AddEndpoint    endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	AllEndpoint    endpoint.Endpoint
	DragEndpoint   endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		AddEndpoint:    makeAddEndpoint(s),
		UpdateEndpoint: makeUpdateEndpoint(s),
		DeleteEndpoint: makeDeleteEndpoint(s),
		AllEndpoint:    makeAllEndpoint(s),
		DragEndpoint:   makeDragEndpoint(s),
	}

	for _, m := range dmw["Add"] {
		eps.AddEndpoint = m(eps.AddEndpoint)
	}
	for _, m := range dmw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}
	for _, m := range dmw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range dmw["All"] {
		eps.AllEndpoint = m(eps.AllEndpoint)
	}
	for _, m := range dmw["Drag"] {
		eps.DragEndpoint = m(eps.DragEndpoint)
	}

	return eps
}

func makeDragEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(dragRequest)
		res, err := s.Drag(ctx, req.DragId, req.DropId)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.All(ctx)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
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
		req := request.(updateRequest)
		err = s.Update(ctx, req.Id, req.Name, req.Alias, req.Icon, req.Path, req.Method, req.Desc, req.ParentId, req.Menu)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addRequest)
		err = s.Add(ctx, req.Name, req.Alias, req.Icon, req.Path, req.Method, req.Desc, req.ParentId, req.Menu)
		return encode.Response{
			Error: err,
		}, err
	}
}
