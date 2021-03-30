package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var defaultMetricPath = "/metrics"
var reqCnt *prometheus.CounterVec
var reqDur *prometheus.HistogramVec

func init() {
	reqCnt = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Count of HTTP requests processed",
		},
		[]string{"code", "method", "url"},
	)

	reqDur = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_ms",
			Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192},
			Help:    "HTTP request latencies in ms",
		},
		[]string{"code", "method", "url"},
	)
}

// UsePrometheus adds the middleware to a gin engine.
func UsePrometheus(e *gin.Engine) {
	e.Use(prometheusMiddleware())
	e.GET(defaultMetricPath, prometheusHandler())
}

// prometheusMiddleware defines handler function for middleware
func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == defaultMetricPath {
			c.Next()
			return
		}

		start := time.Now()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := time.Since(start) / time.Millisecond

		reqDur.WithLabelValues(
			status,
			c.Request.Method,
			c.Request.URL.Path).Observe(float64(elapsed))
		reqCnt.WithLabelValues(status,
			c.Request.Method,
			c.Request.URL.Path).Inc()
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
