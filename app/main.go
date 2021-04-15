package main

import (
	"io"
	"net/http"

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
	e.GET("/external_request/:delay", ExternalRequest)
}

// controllers

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func ExternalRequest(c *gin.Context) {
	delay := c.Params.ByName("delay")
	if delay == "" {
		delay = "5"
	}
	resp, _ := http.Get("http://httpbin/delay/" + delay)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.JSON(resp.StatusCode, gin.H{"body": body})
}
