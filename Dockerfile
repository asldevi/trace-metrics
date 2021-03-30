FROM golang:alpine

# ENV GO111MODULE=on

RUN set -ex

RUN mkdir -p /app

WORKDIR /app

ADD ./go.mod /app
ADD ./go.sum /app

RUN go mod download

ADD . /app

RUN go build -o app/main app/*.go

CMD [ "./app/main" ]

