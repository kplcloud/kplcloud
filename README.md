# kplcloud

开普勒平台开源版

kplcloud是一个基于了kubernetes的应用管理系统，通过可视化的页面对应用进行管理，降低容器化成本，同时也降低了Docker及Kubernetes的学习门槛。

kplcloud已服务于宜人财富部分业务，稳定管理了上百个应用，近千个容器。

![](http://source.qiniu.cnd.nsini.com/images/2019/08/70/29/eb/20190813-e6d29094aab8be96ecb77ad029a70896.jpeg?imageView2/2/w/1280/interlace/0/q/100)

## 架构设计

该平台提供了一整套解决方案。

## 平台演示

演示地址: [https://kplcloud.nsini.com/about.html](https://kplcloud.nsini.com/about.html)

- 感谢 [@icowan](https://github.com/icowan) 赞助三台服务器
- 感谢 [@xzghua](https://github.com/xzghua) 赞助一台服务器

所用到的相关服务，组件分别部署在阿里云，腾讯云服务器上。资源非常有限，仅供大家体验，希望不用过度使用。

## 安装说明

平台后端基于[go-kit](https://github.com/go-kit/kit)、前端基于[ant-design](https://github.com/ant-design/ant-design)(版本略老)框架进行开发。

后端所使用到的依赖全部都在[go.mod](go.mod)里，前端的依赖在`package.json`，详情的请看`yarn.lock`，感谢开源社区的贡献。

后端代码: [https://github.com/kplcloud/kplcloud](https://github.com/kplcloud/kplcloud)

前端代码: [https://github.com/kplcloud/kpaas-frontend](https://github.com/kplcloud/kpaas-frontend)

### 安装教程

[安装教程](https://docs.nsini.com/install/kpaas.html)

### 依赖

- Golang 1.12+ [安装手册](https://golang.org/dl/)
- MySQL 5.7+ (大多数据都存在mysql)
- Docker 18.x+ [安装](https://docs.docker.com/install/)
- RabbitMQ (主要用于消息队列)
- Jenkins 2.176.2+ (老版本对java适配可能会有问题，尽量使用新版本)

## 快速开始

1. 克隆

```
$ mkdir -p $GOPATH/src/github.com/kplcloud
$ cd $GOPATH/src/github.com/kplcloud
$ git clone https://github.com/kplcloud/kplcloud.git
$ cd kplcloud
```

2. 配置文件准备

    - 将连接Kubernets的kubeconfig文件放到该项目目录
    - app.cfg文件配置也放到该项目目录app.cfg配置请参考 [配置文件解析](https://docs.nsini.com/start/config.html)

3. docker-compose 启动

```
$ cd install/docker-compose
$ docker-compose up
```

4. make 启动

```
$ make run
```

## 文档

[文档](https://docs.nsini.com/)

### 视频教程

- [本地启动](https://www.bilibili.com/video/av75847198/)
- [本地连接K8S](https://www.bilibili.com/video/av75890739/)
- [创建一个应用](https://www.bilibili.com/video/av75898315/)

## 成员

- **[@icowan](https://github.com/icowan)**
- **[@yuntinghu](https://github.com/yuntinghu)**
- **[@soup-zhang](https://github.com/soup-zhang)**
- **[@xzghua](https://github.com/xzghua)**

### 支持我们

### 技术交流

- QQ群: 722578340
