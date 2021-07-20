/**
 * @Time: 2020/12/27 12:19
 * @Author: solacowa@gmail.com
 * @File: service_tracer
 * @Software: GoLand
 */

package server

import (
	"fmt"
	"io"

	"github.com/icowan/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

// 使用jaeger
func newJaegerTracer(config *config.Config) (tracer opentracing.Tracer, closer io.Closer, err error) {
	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  config.GetString("tracer", "jaeger.type"),        //固定采样
			Param: float64(config.GetInt("tracer", "jaeger.param")), //1=全采样、0=不采样
		},
		Reporter: &jaegerConfig.ReporterConfig{
			//QueueSize:          200, // 缓冲区越大内存消耗越大,默认100
			LogSpans:           config.GetBool("tracer", "jaeger.logspans"),
			LocalAgentHostPort: config.GetString("tracer", "jaeger.host"),
		},
		ServiceName: fmt.Sprintf("%s.%s", appName, namespace),
	}
	metricsFactory := prometheus.New()
	tracer, closer, err = cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger), jaegerConfig.Metrics(metricsFactory))
	if err != nil {
		return
	}
	opentracing.SetGlobalTracer(tracer)
	return
}
