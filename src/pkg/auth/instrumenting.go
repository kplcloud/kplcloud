package auth

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) Login(ctx context.Context, name, password string) (rs string, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "detail").Add(1)
		s.requestLatency.With("method", "detail").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Login(ctx, name, password)
}
