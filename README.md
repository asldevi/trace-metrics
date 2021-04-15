
### About

This repo aims to check discrepancies in request latencies as measured 
by tracing and metrics instrumentations.


### Usage

- Run application, Prometheus and Jaeger with the following.
    $ docker-compose up --build -d

- Check that the app is responding to requests
    - hit few requests with different delays
    $ curl http://0.0.0.0:8080/external_request/1         #   send request with delay of 1 sec

- Check that the app is traced by pointing your browser to http://0.0.0.0:16686/

- Check the metrics on grafana - http://0.0.0.0:3001
     - set a new datasource as http://prometheus:9090/
     - create a new panel with metrics "sum(rate(request_duration_ms_sum[1m])) / sum(rate(request_duration_ms_count[1m]))" for measuring the request latency

