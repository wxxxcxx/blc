FROM golang:latest as builder

ENV GO111MODULE=on \
    GOBIN=/app \
    GIN_MODE=release

WORKDIR /app

COPY . .

RUN go install github.com/iawia002/lux@latest

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM jrottenberg/ffmpeg:4.1-centos7

ENV ROOT=/data \
    COOKIE=/data/cookie \
    INTERVAL=3600

VOLUME [ "/data" ]

WORKDIR /app

COPY --from=builder ["/app/lux", "/app/blc", "/app/start.sh", "./"]


ENTRYPOINT [ "/app/start.sh" ]
