SHELL = /bin/bash

#SCRIPT_DIR         = $(shell pwd)/etc/script
# Select the Go version
BUILD_IMAGE_SERVER  = golang:1.22
# Select the Node.js version
BUILD_IMAGE_WEB     = node:20
# Project name
PROJECT_NAME        = github.com/huuhoaitvn/gin-vue-admin/server
# Config file path
CONFIG_FILE         = config.yaml
# Image registry namespace
IMAGE_NAME          = gva
# Image registry address
REPOSITORY          = registry.cn-hangzhou.aliyuncs.com/${IMAGE_NAME}
# Image tag/version
TAGS_OPT           ?= latest
PLUGIN             ?= email

# Build both web and server inside containers
build: build-web build-server
	docker run --name build-local --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE_SERVER} make build-local

# Build the web frontend inside a container
build-web:
	docker run --name build-web-local --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE_WEB} make build-web-local

# Build the server backend inside a container
build-server:
	docker run --name build-server-local --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE_SERVER} make build-server-local

# Build the web Docker image
build-image-web:
	@cd web/ && docker build -t ${REPOSITORY}/web:${TAGS_OPT} .

# Build the server Docker image
build-image-server:
	@cd server/ && docker build -t ${REPOSITORY}/server:${TAGS_OPT} .

# Build both web and server on the local host
build-local:
	if [ -d "build" ];then rm -rf build; else echo "OK!"; fi \
	&& if [ -f "/.dockerenv" ];then echo "OK!"; else  make build-web-local && make build-server-local; fi \
	&& mkdir build && cp -r web/dist build/ && cp server/server build/ && cp -r server/resource build/resource

# Build the web frontend on the local host
build-web-local:
	@cd web/ && if [ -d "dist" ];then rm -rf dist; else echo "OK!"; fi \
	&& yarn config set registry http://mirrors.cloud.tencent.com/npm/ && yarn install && yarn build

# Build the server backend on the local host
build-server-local:
	@cd server/ && if [ -f "server" ];then rm -rf server; else echo "OK!"; fi \
	&& go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.io,direct \
	&& go env -w CGO_ENABLED=0 && go env  && go mod tidy \
	&& go build -ldflags "-B 0x$(shell head -c8 /dev/urandom|od -An -tx1|tr -d ' \n') -X main.Version=${TAGS_OPT}" -v

# Build a combined web + server Docker image
image: build
	docker build -t ${REPOSITORY}/gin-vue-admin:${TAGS_OPT} -f deploy/docker/Dockerfile .

# Preview build: produce combined plus split web/server images
images: build build-image-web build-image-server
	docker build -t ${REPOSITORY}/all:${TAGS_OPT} -f deploy/docker/Dockerfile .

# Generate the Swagger API documentation
doc:
	@cd server && swag init

# Quick-package a plugin: make plugin PLUGIN="<plugin folder name, defaults to email>"
plugin:
	if [ -d ".plugin" ];then rm -rf .plugin ; else echo "OK!"; fi && mkdir -p .plugin/${PLUGIN}/{server/plugin,web/plugin} \
	&& if [ -d "server/plugin/${PLUGIN}" ];then cp -r server/plugin/${PLUGIN} .plugin/${PLUGIN}/server/plugin/ ; else echo "OK!"; fi \
	&& if [ -d "web/src/plugin/${PLUGIN}" ];then cp -r web/src/plugin/${PLUGIN} .plugin/${PLUGIN}/web/plugin/ ; else echo "OK!"; fi \
	&& cd .plugin && zip -r ${PLUGIN}.zip ${PLUGIN} && mv ${PLUGIN}.zip ../ && cd ..
