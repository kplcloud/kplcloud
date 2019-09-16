/**
 * @Time : 2019/6/25 4:07 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : service
 * @Software: GoLand
 */

package template

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/paginator"
)

var (
	ErrInvalidArgument          = errors.New("invalid argument")
	ErrTemplateParamsRefused    = errors.New("参数校验未通过")
	ErrTemplateUpdateIdNotFound = errors.New("要修改的模板ID不合法")
	ErrTemplateList             = errors.New("获取模板列表失败")
	ErrTemplateListCount        = errors.New("获取模板列表总数失败")
)

type Service interface {
	// 获取单个模板列表
	Get(ctx context.Context, id int) (resp *types.Template, err error)

	// 创建模板
	Post(ctx context.Context, req templateRequest) (err error)

	// 更新模板
	Update(ctx context.Context, req templateRequest) (err error)

	// 删除模版
	Delete(ctx context.Context, id int) (err error)

	// 模板列表
	List(ctx context.Context, name string, page, limit int) (res map[string]interface{}, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{
		logger:     logger,
		repository: repository,
	}
}

func (c *service) Get(ctx context.Context, id int) (resp *types.Template, err error) {
	return c.repository.Template().FindById(id)
}

func (c *service) Post(ctx context.Context, req templateRequest) (err error) {
	if req.Kind == "" || req.Detail == "" || req.Name == "" {
		return ErrTemplateParamsRefused
	}
	err = c.repository.Template().Create(req.Name, req.Kind, req.Detail)
	return err
}

func (c *service) Update(ctx context.Context, req templateRequest) (err error) {
	if req.Id <= 0 {
		return ErrTemplateUpdateIdNotFound
	}
	if req.Kind == "" || req.Detail == "" || req.Name == "" {
		return ErrTemplateParamsRefused
	}
	return c.repository.Template().Update(req.Id, req.Name, req.Kind, req.Detail)
}

func (c *service) Delete(ctx context.Context, id int) (err error) {
	return c.repository.Template().DeleteById(id)
}

func (c *service) List(ctx context.Context, name string, page, limit int) (res map[string]interface{}, err error) {
	count, err := c.repository.Template().Count(name)
	if err != nil {
		_ = level.Error(c.logger).Log("template", "Count", "err", err.Error())
		return nil, ErrTemplateListCount
	}

	p := paginator.NewPaginator(page, limit, count)

	list, err := c.repository.Template().FindOffsetLimit(name, p.Offset(), limit)
	if err != nil {
		_ = level.Error(c.logger).Log("template", "FindOffsetLimit", "err", err.Error())
		return nil, ErrTemplateList
	}
	res = map[string]interface{}{
		"list": list,
		"page": p.Result(),
	}
	return
}
