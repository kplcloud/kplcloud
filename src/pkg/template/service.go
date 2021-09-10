/**
 * @Time : 2021/9/10 7:07 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package template

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	// Add 添加模版
	Add(ctx context.Context, kind, alias, rules, content string) (err error)
	// Delete 删除模版
	Delete(ctx context.Context, kind string) (err error)
	// Update 更新模版
	Update(ctx context.Context, kind, alias, rules, content string) (err error)
	// List 模版列表
	List(ctx context.Context, searchValue string, page, pageSize int64) (res []infoResult, total int, err error)
	// Info 模版详情
	Info(ctx context.Context, kind string) (res infoResult, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) Delete(ctx context.Context, kind string) (err error) {
	panic("implement me")
}

func (s *service) Update(ctx context.Context, kind, alias, rules, content string) (err error) {
	panic("implement me")
}

func (s *service) List(ctx context.Context, searchValue string, page, pageSize int64) (res []infoResult, total int, err error) {
	panic("implement me")
}

func (s *service) Info(ctx context.Context, kind string) (res infoResult, err error) {
	panic("implement me")
}

func (s *service) Add(ctx context.Context, kind, alias, rules, content string) (err error) {
	if err = s.repository.K8sTpl(ctx).Save(ctx, &types.K8sTemplate{
		Kind:    types.Kind(kind),
		Alias:   alias,
		Rules:   rules,
		Content: content,
	}); err != nil {
		return encode.ErrTemplateSave.Wrap(err)
	}
	return nil
}

func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "service", "template")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}
