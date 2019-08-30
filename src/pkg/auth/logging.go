package auth

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) Login(ctx context.Context, name, password string) (rs string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "login",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Login(ctx, name, password)
}
