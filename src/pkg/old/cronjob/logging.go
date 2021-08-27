/**
 * Created by GoLand.
 * Email: xzghua@gmail.com
 * Date: 2019-07-09
 * Time: 15:01
 */
package cronjob

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport/http"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) AddCronJob(ctx context.Context, acj addCronJob) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AddCronJob",
			"name", acj.Name,
			"args", acj.Args,
			"gitType", acj.GitType,
			"gitPath", acj.GitPath,
			"image", acj.Image,
			"namespace", acj.Namespace,
			"schedule", acj.Schedule,
			"confMap", acj.ConfMap,
			"logPath", acj.LogPath,
			"addType", acj.AddType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AddCronJob(ctx, acj)
}

func (s *loggingService) CronJobList(ctx context.Context, cjl cronJobList) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "CronJobList",
			"name", cjl.Name,
			"group", cjl.Group,
			"page", cjl.Page,
			"limit", cjl.Limit,
			"namespace", cjl.Namespace,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, cjl.Name, cjl.Namespace, cjl.Group, cjl.Page, cjl.Limit)
}

func (s *loggingService) CronJobDel(ctx context.Context, cjd cronJobDel) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Delete",
			"name", cjd.Name,
			"namespace", cjd.Namespace,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, cjd.Name, cjd.Namespace)
}

func (s *loggingService) CronJobAllDel(ctx context.Context, cjd cronJobAllDel) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "DeleteJobAll",
			"namespace", cjd.Namespace,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteJobAll(ctx, cjd.Namespace)
}

func (s *loggingService) CronJobUpdate(ctx context.Context, acj addCronJob) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Put",
			"name", acj.Name,
			"args", acj.Args,
			"gitType", acj.GitType,
			"gitPath", acj.GitPath,
			"image", acj.Image,
			"namespace", acj.Namespace,
			"schedule", acj.Schedule,
			"confMap", acj.ConfMap,
			"logPath", acj.LogPath,
			"addType", acj.AddType,
			"paramName", acj.ParamName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Put(ctx, acj.ParamName, acj)
}

func (s *loggingService) CronJobDetail(ctx context.Context, cjl cronJobDetail) (res *DetailReturnData, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Detail",
			"name", cjl.Name,
			"namespace", cjl.Namespace,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, cjl.Name, cjl.Namespace)
}

func (s *loggingService) CronJobUpdateLog(ctx context.Context, cjl cronJobLogUpdate) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "UpdateLog",
			"name", cjl.Name,
			"namespace", cjl.Namespace,
			"logPath", cjl.LogPath,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UpdateLog(ctx, cjl)
}
