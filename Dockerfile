# FROM hub.hitry.io/base/golang:1.21-alpine3.17 AS golang-builder
# LABEL maintainer="aichy@sf.com"

# docker pull hub.hitry.io/base/golang@sha256:17638a3d703bf2143b7fea818ac1db841ab2b18ee60aa96b1305001ba053d8bd

# ARG GOPROXY
# ENV GOPROXY ${GOPROXY:-direct}
# ENV GOPROXY=https://proxy.golang.com.cn,direct

# ENV GOPATH /go
# ENV GOROOT /usr/local/go
# ENV GOPROXY=https://goproxy.cn,direct
# ENV PACKAGE wstester
# ENV BUILD_DIR ${GOPATH}/src/${PACKAGE}

# COPY . ${BUILD_DIR}
# WORKDIR ${BUILD_DIR}
# RUN apk --no-cache add build-base && make clean build

# RUN CGO_ENABLE=0 GO111MODULE=on go build -o wstester .

# RUN chmod 755 wstester
# RUN cp wstester /usr/bin/wstester
# RUN cp entrypoint.sh /entrypoint.sh
# RUN cp -r config /config

FROM hub.hitry.io/base/alpine:latest

#docker pull hub.hitry.io/base/alpine@sha256:d7342993700f8cd7aba8496c2d0e57be0666e80b4c441925fc6f9361fa81d10e
LABEL maintainer="lifue"

ENV TZ "Asia/Shanghai"
RUN echo "Asia/Shanghai" > /etc/timezone

COPY ./wstester /usr/bin/wstester
COPY ./entrypoint.sh /entrypoint.sh
COPY ./config/config_docker.yaml /config/config.yaml
RUN chmod 755 /entrypoint.sh

WORKDIR /

EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]