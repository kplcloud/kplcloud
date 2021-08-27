FROM golang:1.14.15-alpine3.11 as build-env

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache \
		ca-certificates \
		tzdata \
		git \
		openssh \
		vim \
		make \
		mercurial \
		subversion \
		bzr \
		fossil \
		&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
		&& echo "Asia/Shanghai" > /etc/timezone \
		&& apk del tzdata \
		&& rm -rf /var/cache/apk/*

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
ENV BUILDPATH=github.com/icowan/kplcloud
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
WORKDIR /go/src/${BUILDPATH}/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v

# 前端打包基础镜像
#FROM node:12.13.0-alpine AS node-build-env
#
#RUN mkdir /opt/build
#COPY ./web/admin/ /opt/build
#WORKDIR /opt/build
#
#RUN yarn config set registry https://registry.npm.taobao.org
#RUN yarn build --registry https://registry.npm.taobao.org

# 运行镜像
FROM alpine:3.11
RUN apk add --no-cache \
		ca-certificates \
		curl \
		tzdata \
		&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
		&& echo "Asia/Shanghai" > /etc/timezone \
		&& apk del tzdata \
		&& rm -rf /var/cache/apk/*

COPY --from=build-env /go/bin/cmd /usr/local/kplcloud/bin/kplcloud
#COPY --from=node-build-env /opt/build/dist/ /usr/local/kplcloud/web/admin

WORKDIR /usr/local/kplcloud/
ENV PATH=$PATH:/usr/local/kplcloud/bin/

#COPY ./app.prod.cfg /usr/local/kplcloud/etc/app.cfg

CMD ["kplcloud", "start", "-c", "/usr/local/kplcloud/etc/app.cfg"]