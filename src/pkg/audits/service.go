/**
 * @Time : 2021/10/27 5:17 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package audits

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/kplcloud/kplcloud/src/repository"
)

type Middleware func(Service) Service

// Service 审核日志模块
// 所有非Get请求的日志都会在这展示
// 可根据类型、路由、空间、集群、服务等过滤
// 只有查询功能不提供其他操作
type Service interface {
	// List 获取审计日志列表
	// TODO: 应该会有很多查询条件，具体的后面再定
	List(ctx context.Context, query string, page, pageSize int) (res []auditResult, total int, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) List(ctx context.Context, query string, page, pageSize int) (res []auditResult, total int, err error) {
	list, total, err := s.repository.Audit(ctx).List(ctx, query, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for _, v := range list {
		res = append(res, auditResult{
			Username:       v.User.Username,
			Method:         v.Permission.Method,
			Remark:         v.Remark,
			PermissionName: v.Permission.Name,
			Request:        v.Request,
			Response:       v.Response,
			Headers:        v.Headers,
			TimeSince:      v.TimeSince,
			Status:         string(v.Status),
			Url:            v.Url,
			TraceId:        v.TraceId,
			CreatedAt:      v.CreatedAt,
			Name:           v.Name,
			Namespace:      v.Namespace,
			Cluster:        v.Cluster.Alias,
		})
	}
	return
}

func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "pkg.audit", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}
