/**
 * @Time : 2019-07-29 11:21
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package market

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	"time"
)

type listRequest struct {
	Page     int
	Limit    int
	language []string
	name     string
	status   int
}

type postRequest struct {
	Id         int64               `json:"id"`
	Name       string              `json:"name"`
	Language   repository.Language `json:"language"`   // 语言
	Version    string              `json:"version"`    // 版本
	Detail     string              `json:"detail"`     // 使用的例子
	Desc       string              `json:"desc"`       // 描述
	Dockerfile string              `json:"dockerfile"` // 做镜像的Dockerfile 点详情的时候显示 必填
	Status     int64               `json:"status"`     // 状态是否可用 或是否审核通过
	Score      int64               `json:"score"`      // 评分
	FullPath   string              `json:"full_path"`  // 全路径 hub.kpaas.nsini.com/golang/goalng-1.11.2:v1.0
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		err := s.Post(ctx, req.Name, req.Language.String(), req.Version, req.Detail, req.Desc, req.Dockerfile, req.FullPath, req.Status)
		return encode.Response{Err: err}, err
	}
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		res, err := s.Detail(ctx, req.Id)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		res, count, err := s.List(ctx, req.Page, req.Limit, req.language, req.status, req.name)
		return encode.Response{Err: err, Data: map[string]interface{}{
			"page":  paginator.NewPaginator(req.Page, req.Limit, int(count)).Result(),
			"items": res,
		}}, err
	}
}

func makePutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		err := s.Put(ctx, req.Id, req.Name, req.Language.String(), req.Version, req.Detail, req.Desc, req.Dockerfile, req.FullPath, req.Status)
		return encode.Response{Err: err}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		err := s.Delete(ctx, req.Id)
		return encode.Response{Err: err}, err
	}
}
