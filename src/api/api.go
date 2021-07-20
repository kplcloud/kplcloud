/**
 * @Time: 2020/12/27 22:06
 * @Author: solacowa@gmail.com
 * @File: api
 * @Software: GoLand
 */

package api

import (
	"github.com/go-kit/kit/log"
	"github.com/icowan/config"
	kitcache "github.com/icowan/kit-cache"
	"github.com/opentracing/opentracing-go"
)

type Service interface {
}

type api struct {
	logger log.Logger
}

// 中间件有顺序,在后面的会最先执行
func NewApi(logger log.Logger, traceId string, tracer opentracing.Tracer, cfg *config.Config, cache kitcache.Service) Service {
	logger = log.With(logger, "api", "Api")

	// 如果tracer有的话
	if tracer != nil {

	}

	// 如果有cache的话
	if cache != nil {
	}

	return &api{}
}
