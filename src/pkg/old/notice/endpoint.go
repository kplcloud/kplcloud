/**
 * @Time : 2019-07-02 10:41
 * @Author : soupzhb@gmail.com
 * @File : endpoint.go
 * @Software: GoLand
 */

package notice

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"gopkg.in/guregu/null.v3"
)

type noticeListRequest struct {
	Title       string
	Type        string
	IsRead      string
	Page, Limit int
}

type detailRequest struct {
	NoticeId int
}

type noticeInfo struct {
	Id       int64     `json:"id"`
	Avatar   string    `json:"avatar"`
	Title    string    `json:"title"`
	Datetime null.Time `json:"datetime"`
	Type     string    `json:"type"`
}

type clearRequest struct {
	ClearType int
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(noticeListRequest)
		param := map[string]string{}
		param["title"] = req.Title
		param["type"] = req.Type
		res, err := s.List(ctx, param, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeTipsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Tips(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeCountReadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(noticeListRequest)
		param := map[string]string{}
		param["is_read"] = req.IsRead
		param["type"] = req.Type
		res, err := s.CountRead(ctx, param)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeClearAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(clearRequest)
		err := s.ClearAll(ctx, req.ClearType)
		return encode.Response{Err: err, Data: nil}, err
	}
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(detailRequest)
		res, err := s.Detail(ctx, req.NoticeId)
		return encode.Response{Err: err, Data: res}, err
	}
}
