/**
 * @Time : 3/9/21 5:58 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package template

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) Add(ctx context.Context, kind, alias, rules, content string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Add",
			"kind", kind,
			"alias", alias,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Add(ctx, kind, alias, rules, content)
}

func (s *logging) Delete(ctx context.Context, kind string) (err error) {
	panic("implement me")
}

func (s *logging) Update(ctx context.Context, kind, alias, rules, content string) (err error) {
	panic("implement me")
}

func (s *logging) List(ctx context.Context, searchValue string, page, pageSize int64) (res []interface{}, total int, err error) {
	panic("implement me")
}

func (s *logging) Info(ctx context.Context, kind string) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Info",
			"kind", kind,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Info(ctx, kind)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "template", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
