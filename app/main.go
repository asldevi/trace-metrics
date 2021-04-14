package main

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

var (
	Tracer opentracing.Tracer
)

func main() {
	Tracer = InitTracing("sample-app")

	router := gin.Default()

	UseTracing(router)
	UsePrometheus(router)
	registerRoutes(router)
	router.Run()
}

func registerRoutes(e *gin.Engine) {
	e.GET("/ping", Ping)
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
