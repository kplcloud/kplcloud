/**
 * @Time: 2021/10/23 16:45
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package application

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/kplcloud/kplcloud/src/repository"
)

type Middleware func(Service) Service

// Service 与应用相关的操作也就是核心模块
// 此模块只对应用进行操作，如需对Deployment进行操作使用deployment模块
// 此模块应该会继承deployment、service、ingress、pvc、configmap、build、git等等模块，在此模块中会直接调用以上模块
// 这样的话该模块方法会非常多，合理么？先这么设计，以后会怎么样再说...
// 创建服务要求加到组里，如果没有组给弹窗创建一个，应用创建成功之后自动将该应用加入到组里
type Service interface {
	// List 获取应用列表
	// 集群权限及组权限验证在中间件完成
	// 中间件取得ClusterId、GroupIds 并通过GroupIds和NamespaceIds过滤出所有的appName TODO: Group 关系: groupId nsId appId
	// namespace 可以为空，如果前端没选择namespace, 则取所有namespace下的应用(前题条件是这些namespace在Groups里)
	// query 可以为空 过滤 标签或应用名称或别名 TODO: 是否需要加一个标签筛选？
	List(ctx context.Context, clusterId int64, namespace string, name []string, page, pageSize int) (res []appResult, total int, err error)
	// Sync 同步空间下的所有Deployment
	Sync(ctx context.Context, clusterId int64, namespace string) (err error)

	//Add(ctx context.Context, clusterId int64, ns, name, alias, remark string) (err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
}

func (s *service) Sync(ctx context.Context, clusterId int64, namespace string) (err error) {
	panic("implement me")
}

func (s *service) List(ctx context.Context, clusterId int64, namespace string, name []string, page, pageSize int) (res []appResult, total int, err error) {
	// logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	s.repository.Application(ctx).List(ctx, clusterId, namespace, nil, page, pageSize)

	return
}

func New(traceId string, logger log.Logger, repository repository.Repository) Service {
	logger = log.With(logger, "application", "service")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
	}
}
