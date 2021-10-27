/**
 * @Time: 2019-06-29 09:34
 * @Author: solacowa@gmail.com
 * @File: instrumenting
 * @Software: GoLand
 */

package public

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

func (s *instrumentingService) GitPost(ctx context.Context, namespace, name, token, keyWord, branch string, req gitlabHook) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GitPost").Add(2)
		s.requestLatency.With("method", "GitPost").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GitPost(ctx, namespace, name, token, keyWord, branch, req)
}

func (s *instrumentingService) PrometheusAlert(ctx context.Context, req *prometheusAlerts) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "PrometheusAlert").Add(2)
		s.requestLatency.With("method", "PrometheusAlert").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PrometheusAlert(ctx, req)
}
