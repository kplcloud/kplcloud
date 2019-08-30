/**
 * @Time : 2019-07-09 16:04
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package build

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type buildRequest struct {
	Version             string `json:"version"`
	BuildEnv            string `json:"build_env"`
	BuildEnvDescription string `json:"build_env_description"`
	BuildTime           string `json:"build_time"`
	GitType             string `json:"git_type"`
}

type abortRequest struct {
	Number int
}

type buildConsoleRequest struct {
	Number, Start int
}

type listRequest struct {
	Page  int
	Limit int
	Types string
}
type confRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func makeBuildEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(buildRequest)
		err := s.Build(ctx, req.GitType, req.Version, req.BuildEnv, req.BuildEnvDescription, req.BuildTime)
		return encode.Response{Err: err}, err
	}
}

func makeBuildConsoleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(buildConsoleRequest)
		res, end, err := s.BuildConsole(ctx, req.Number, req.Start)
		return encode.Response{Err: err, Data: map[string]interface{}{
			"output": res,
			"start":  req.Start,
			"end":    end,
		}}, err
	}
}

func makeCronBuildConsoleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(buildConsoleRequest)
		res, end, err := s.CronBuildConsole(ctx, req.Number, req.Start)
		return encode.Response{Err: err, Data: map[string]interface{}{
			"output": res,
			"start":  req.Start,
			"end":    end,
		}}, err
	}
}

func makeAbortBuildEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(abortRequest)
		err := s.AbortBuild(ctx, req.Number)
		return encode.Response{Err: err}, err
	}
}

func makeCronHistoryEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, err := s.CronHistory(ctx, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeHistoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		res, err := s.History(ctx, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeRollBackEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(abortRequest)
		err := s.Rollback(ctx, int64(req.Number))
		return encode.Response{Err: err}, err
	}
}

func makeBuildConfEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(confRequest)
		data, err := s.BuildConf(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err, Data: data}, err
	}
}
