FROM golang:1.25 AS builder
RUN set -eux && \
    apt-get update && \
    apt-get install -y git

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,https://mirrors.aliyun.com/goproxy/,direct
ENV GOSUMDB=off
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN rm -f go.mod && \
    go mod  init k8s-platform-go &&\
    go mod  tidy && \
    go mod download && \
    go build -ldflags="-w -s" -o k8s-dev ./cmd/api/

# 启动阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata && \
    mkdir /app

ENV GIN_MODE=release
WORKDIR /app
COPY --from=builder /app/k8s-dev .
COPY --from=builder /app/configs .
RUN chmod 755  k8s-dev && \
    ls -l /app/

EXPOSE 8081
# 启动应用
CMD ["/app/k8s-dev"]