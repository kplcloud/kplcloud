package auth

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func (s *instrumentingService) Register(ctx context.Context, username, email, password, mobile, remark string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Register").Add(1)
		s.requestLatency.With("method", "Register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.Register(ctx, username, email, password, mobile, remark)
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram) Middleware {
	return func(s Service) Service {
		return &instrumentingService{
			requestCount:   counter,
			requestLatency: latency,
			next:           s,
		}
	}
}

func (s *instrumentingService) Login(ctx context.Context, username, password string) (rs string, sessionTimeout int64, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Login").Add(1)
		s.requestLatency.With("method", "Login").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.Login(ctx, username, password)
}
