/**
 * @Time : 2019-06-28 10:33
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package deployment

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/middleware"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) GetYaml(ctx context.Context) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "GetYaml",
			"took", time.Since(begin),
			"name", ctx.Value(middleware.NameContext),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetYaml(ctx)
}

func (s *loggingService) CommandArgs(ctx context.Context, commands []string, args []string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "CommandArgs",
			"took", time.Since(begin),
			"name", ctx.Value(middleware.NameContext),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"err", err,
		)
	}(time.Now())
	return s.Service.CommandArgs(ctx, commands, args)
}

func (s *loggingService) Expansion(ctx context.Context, requestCpu, limitCpu, requestMemory, limitMemory string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Expansion",
			"took", time.Since(begin),
			"name", ctx.Value(middleware.NameContext),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"requestCpu", requestCpu,
			"limitCpu", limitCpu,
			"requestMemory", requestMemory,
			"limitLimit", limitMemory,
			"err", err,
		)
	}(time.Now())
	return s.Service.Expansion(ctx, requestCpu, limitCpu, requestMemory, limitMemory)
}

func (s *loggingService) Stretch(ctx context.Context, num int) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Stretch",
			"took", time.Since(begin),
			"name", ctx.Value(middleware.NameContext),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"num", num,
			"err", err,
		)
	}(time.Now())
	return s.Service.Stretch(ctx, num)
}

func (s *loggingService) GetPvc(ctx context.Context, ns, name string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "GetPvc",
			"took", time.Since(begin),
			"name", name,
			"namespace", ns,
			"err", err,
		)
	}(time.Now())
	return s.Service.GetPvc(ctx, ns, name)
}

func (s *loggingService) BindPvc(ctx context.Context, ns, name, path, claimName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "BindPvc",
			"took", time.Since(begin),
			"name", name,
			"namespace", ns,
			"claimName", claimName,
			"path", path,
			"err", err,
		)
	}(time.Now())
	return s.Service.BindPvc(ctx, ns, name, path, claimName)
}

func (s *loggingService) UnBindPvc(ctx context.Context, ns, name, claimName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "UnBindPvc",
			"took", time.Since(begin),
			"name", name,
			"namespace", ns,
			"claimName", claimName,
			"err", err,
		)
	}(time.Now())
	return s.Service.UnBindPvc(ctx, ns, name, claimName)
}

func (s *loggingService) AddPort(ctx context.Context, ns, name string, req portRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "AddPort",
			"took", time.Since(begin),
			"name", name,
			"namespace", ns,
			"err", err,
		)
	}(time.Now())
	return s.Service.AddPort(ctx, ns, name, req)
}

func (s *loggingService) DelPort(ctx context.Context, ns, name string, portName string, port int32) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "DelPort",
			"took", time.Since(begin),
			"name", name,
			"namespace", ns,
			"portNaeme", portName,
			"port", port,
			"err", err,
		)
	}(time.Now())
	return s.Service.DelPort(ctx, ns, name, portName, port)
}

func (s *loggingService) Logging(ctx context.Context, ns, name, pattern, suffix string, paths []string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Logging",
			"took", time.Since(begin),
			"name", name,
			"namespace", ns,
			"pattern", pattern,
			"suffix", suffix,
			"err", err,
		)
	}(time.Now())
	return s.Service.Logging(ctx, ns, name, pattern, suffix, paths)
}

func (s *loggingService) Probe(ctx context.Context, ns, name string, req probeRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Probe",
			"took", time.Since(begin),
			"name", name,
			"namespace", ns,
			"port", req.Port,
			"probe", req.Probe,
			"path", req.Path,
			"err", err,
		)
	}(time.Now())
	return s.Service.Probe(ctx, ns, name, req)
}

func (s *loggingService) Mesh(ctx context.Context, ns, name, model string) (err error) {
	//ctx = context.WithValue(ctx,"OperationMethod", repository.SwitchModel)

	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Mesh",
			"took", time.Since(begin),
			"name", name,
			"namespace", ns,
			"model", model,
			"err", err,
		)
	}(time.Now())
	return s.Service.Mesh(ctx, ns, name, model)
}

func (s *loggingService) Hosts(ctx context.Context, hosts []string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Hosts",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Hosts(ctx, hosts)
}

func (s *loggingService) VolumeConfig(ctx context.Context, mountPath, subPath string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "VolumeConfig",
			"mountPath", mountPath,
			"subPath", subPath,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.VolumeConfig(ctx, mountPath, subPath)
}
