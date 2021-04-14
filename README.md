
### About

This repo aims to check discrepancies in request latencies as measured 
by tracing and metrics instrumentations.


### Usage

$ docker-compose up --build -d
$ curl http://0.0.0.0:8080/ping         #   check app
$ curl http://0.0.0.0:9090/metrics      #   check prom metrics
$ curl http://0.0.0.0:16686/search      #   check jaeger traces


