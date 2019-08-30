/**
 * @Time : 2019/7/5 11:03 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package configmap

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
	return &loggingService{logger: level.Info(logger), Service: s}
}

func (s *loggingService) GetOne(ctx context.Context, ns, name string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Get",
			"name", name,
			"namespace", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetOne(ctx, ns, name)
}

func (s *loggingService) GetOnePull(ctx context.Context, ns, name string) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Get",
			"name", name,
			"namespace", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetOnePull(ctx, ns, name)
}

func (s *loggingService) List(ctx context.Context, req listRequest) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Get",
			"name", req.Name,
			"namespace", req.Namespace,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, req)
}

func (s *loggingService) Post(ctx context.Context, req postRequest) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Post",
			"name", req.Name,
			"namespace", req.Namespace,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Post(ctx, req)
}

func (s *loggingService) Update(ctx context.Context, req postRequest) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Update",
			"name", req.Name,
			"namespace", req.Namespace,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Update(ctx, req)
}

func (s *loggingService) Delete(ctx context.Context, ns, name string) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Delete",
			"name", name,
			"namespace", ns,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Delete(ctx, ns, name)
}

func (s *loggingService) Sync(ctx context.Context, ns string) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Sync",
			"namespace", ns,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Sync(ctx, ns)
}

func (s *loggingService) CreateConfigMap(ctx context.Context, req createConfigMapRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "CreateConfigMap",
			"namespace", req.Namespace,
			"name", req.Name,
			"desc", req.Desc,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.CreateConfigMap(ctx, req)
}

func (s *loggingService) GetConfigMap(ctx context.Context, ns, name string) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "GetConfigMap",
			"namespace", ns,
			"name", name,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetConfigMap(ctx, ns, name)
}

func (s *loggingService) GetConfigMapData(ctx context.Context, ns, name string, page, limit int) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "GetConfigMapData",
			"namespace", ns,
			"name", name,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetConfigMapData(ctx, ns, name, page, limit)
}

func (s *loggingService) CreateConfigMapData(ctx context.Context, req createConfigMapDataRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "CreateConfigMapData",
			"configMapId", req.ConfigMapId,
			"key", req.Key,
			"value", req.Value,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.CreateConfigMapData(ctx, req)
}

func (s *loggingService) UpdateConfigMapData(ctx context.Context, req configMapDataRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "UpdateConfigMapData",
			"configMapId", req.ConfigMapId,
			"configMapDataId", req.ConfigMapDataId,
			"key", req.Key,
			"value", req.Value,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.UpdateConfigMapData(ctx, req)
}

func (s *loggingService) DeleteConfigMapData(ctx context.Context, req configMapDataRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "DeleteConfigMapData",
			"configMapId", req.ConfigMapId,
			"configMapDataId", req.ConfigMapDataId,
			"key", req.Key,
			"value", req.Value,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.DeleteConfigMapData(ctx, req)
}

func (s *loggingService) ConfEnvList(ctx context.Context, req listRequest) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "GetConfigEnv",
			"name", req.Name,
			"namespace", req.Namespace,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetConfigEnv(ctx, req.Namespace, req.Name, req.Page, req.Limit)
}

func (s *loggingService) CreateConfEnv(ctx context.Context, req configEnvRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "CreateConfigEnv",
			"name", req.Name,
			"namespace", req.Namespace,
			"envDesc", req.EnvDesc,
			"envVar", req.EnvVar,
			"envKey", req.EnvKey,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CreateConfigEnv(ctx, req)
}

func (s *loggingService) UpdateConfEnv(ctx context.Context, req configEnvRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "ConfigEnvUpdate",
			"name", req.Name,
			"namespace", req.Namespace,
			"envDesc", req.EnvDesc,
			"envVar", req.EnvVar,
			"envKey", req.EnvKey,
			"id", req.Id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CreateConfigEnv(ctx, req)
}

func (s *loggingService) DelConfEnv(ctx context.Context, req configEnvRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "ConfigEnvUpdate",
			"name", req.Name,
			"namespace", req.Namespace,
			"id", req.Id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CreateConfigEnv(ctx, req)
}
