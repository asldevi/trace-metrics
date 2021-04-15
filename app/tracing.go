package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	// jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func InitTracing(serviceName string) opentracing.Tracer {
	// take JAEGER_AGENT_HOST / JAEGER_AGENT_PORT from env
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, _ := config.FromEnv()
	fmt.Printf("%+v\n", cfg)

	cfg.ServiceName = serviceName
	cfg.Sampler.Type = jaeger.SamplerTypeConst
	cfg.Sampler.Param = 1

	jaegerLogger := jaegerlog.StdLogger
	metricsFactory := metrics.NullFactory

	tracer, _, _ := cfg.NewTracer(
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	return tracer
}

func TracingMiddleware() gin.HandlerFunc {
	var serverSpan opentracing.Span
	return func(c *gin.Context) {
		spanName := c.FullPath()

		if spanName == "/metrics" {
			// execute other middlewares and return without creating span
			c.Next()
			return
		}

		spanCtx, _ := Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		serverSpan = Tracer.StartSpan(spanName, ext.RPCServerOption(spanCtx))
		serverSpan.SetTag("request.host", c.Request.Host)
		serverSpan.SetTag("http.method", c.Request.Method)

		defer serverSpan.Finish()
		c.Next()
		statusCode := c.Writer.Status()
		serverSpan.SetTag("http.status_code", strconv.Itoa(statusCode))
		if len(c.Errors) > 0 {
			serverSpan.SetTag("error", c.Errors.String())
		}
	}
}

func UseTracing(e *gin.Engine) {
	e.Use(TracingMiddleware())
}
