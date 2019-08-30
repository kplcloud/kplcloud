/**
 * @Time : 2019-07-09 11:31
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package git

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type getDockerfile struct {
	FileName string `json:"file_name"`
}

type gitPath struct {
	Git string `json:"git"`
}

func makeTagsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Tags(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeBranchesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Branches(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeTagsByGitPathEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(gitPath)
		res, err := s.TagsByGitPath(ctx, req.Git)
		return encode.Response{Data: res, Err: err}, err
	}
}

func makeBranchesByGitPathEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(gitPath)
		res, err := s.BranchesByGitPath(ctx, req.Git)
		return encode.Response{Data: res, Err: err}, err
	}
}

func makeGetDockerfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getDockerfile)
		res, err := s.GetDockerfile(ctx, req.FileName)
		return encode.Response{Err: err, Data: res}, err
	}
}
