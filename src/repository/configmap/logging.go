/**
 * @Time : 8/19/21 1:29 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package configmap

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) List(ctx context.Context, clusterId int64, ns, name string, page, pageSize int) (res []types.ConfigMap, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "clusterId",clusterId,"ns",ns,"name",name,"page",page,"pageSize",pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx,clusterId,ns,name,page,pageSize)
}

func (s *logging) FindBy(ctx context.Context, clusterId int64, ns, name string) (res types.ConfigMap, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindBy",
			"clusterId", clusterId,
			"ns", ns,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindBy(ctx, clusterId, ns, name)
}

func (s *logging) Save(ctx context.Context, configMap *types.ConfigMap, data []types.Data) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, configMap, data)
}

func (s *logging) SaveData(ctx context.Context, configMapId int64, key, value string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SaveData",
			"configMapId", configMapId,
			"key", key,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.SaveData(ctx, configMapId, key, value)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "configmap", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
