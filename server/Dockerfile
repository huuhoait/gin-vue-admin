FROM golang:alpine as builder

WORKDIR /go/src/github.com/huuhoait/gin-vue-admin/server
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest

LABEL MAINTAINER="SliverHorn@sliver_horn@qq.com"
# Set timezone
ENV TZ=Asia/Shanghai
RUN apk update && apk add --no-cache tzdata openntpd \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /go/src/github.com/huuhoait/gin-vue-admin/server

COPY --from=0 /go/src/github.com/huuhoait/gin-vue-admin/server/server ./
COPY --from=0 /go/src/github.com/huuhoait/gin-vue-admin/server/resource ./resource/
COPY --from=0 /go/src/github.com/huuhoait/gin-vue-admin/server/config.docker.yaml ./

# Mount volume: if using sqlite, example command: docker run -d -v /host-path/gva.db:/go/src/github.com/huuhoait/gin-vue-admin/server/gva.db -p 8888:8888 --name gva-server-v1 gva-server:1.0
# VOLUME ["/go/src/github.com/huuhoait/gin-vue-admin/server"]

EXPOSE 8888
ENTRYPOINT ./server -c config.docker.yaml
