/**
 * @Time : 2019-07-09 11:32
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package git

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

func (s *loggingService) Tags(ctx context.Context) (res []string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Tags",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"err", err,
		)
	}(time.Now())
	return s.Service.Tags(ctx)
}

func (s *loggingService) Branches(ctx context.Context) (res []string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Branches",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"err", err,
		)
	}(time.Now())
	return s.Service.Branches(ctx)
}

func (s *loggingService) TagsByGitPath(ctx context.Context, gitPath string) (res []string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "TagsByGitPath",
			"gitPath", gitPath,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.TagsByGitPath(ctx, gitPath)
}

func (s *loggingService) BranchesByGitPath(ctx context.Context, gitPath string) (res []string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "BranchesByGitPath",
			"gitPath", gitPath,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.BranchesByGitPath(ctx, gitPath)
}

func (s *loggingService) GetDockerfile(ctx context.Context, fileName string) (res string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "GetDockerfile",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"fileName", fileName,
			"err", err,
		)
	}(time.Now())
	return s.Service.GetDockerfile(ctx, fileName)
}
