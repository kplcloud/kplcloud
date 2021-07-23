/**
 * @Time : 5/7/21 9:43 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package syspermission

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Service interface {
	// 限制列表
	List(ctx context.Context, page, pageSize int) (res []result, total int, err error)

	// 展示所有的权限
	All(ctx context.Context) (res []result, err error)

	// 添加权限
	// 别名，图标，路径，method，备注，上级ID，是否是菜单
	Add(ctx context.Context, name, alias, icon, path, method, desc string, parentId int64, menu bool) (err error)

	// 删除权限
	Delete(ctx context.Context, id int64) (err error)

	// 更新权限
	// 权限ID，别名，图标，路径，method，备注，上级ID，是否是菜单
	Update(ctx context.Context, id int64, name, alias, icon, path, method, desc string, parentId int64, menu bool) (err error)

	// 移动Permission
	Drag(ctx context.Context, dragKey, dropKey int64) (res []result, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) Drag(ctx context.Context, dragKey, dropKey int64) (res []result, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Drag")
	dragPerm, err := s.repository.SysPermission().Find(ctx, dragKey)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysPermission", "Find", "err", err.Error())
		return
	}
	dropPerm, err := s.repository.SysPermission().Find(ctx, dropKey)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysPermission", "Find", "err", err.Error())
		return
	}

	dragPerm.ParentId = dropPerm.Id

	if err = s.repository.SysPermission().Save(ctx, &dragPerm); err != nil {
		_ = level.Error(logger).Log("repository.SysPermission", "Update", "err", err.Error())
		return
	}

	return s.All(ctx)
}

func (s *service) All(ctx context.Context) (res []result, err error) {
	s.repository.SysPermission()
	return
}

func (s *service) Add(ctx context.Context, name, alias, icon, path, method, desc string, parentId int64, menu bool) (err error) {
	return s.repository.SysPermission().Save(ctx, &types.SysPermission{
		Icon:        icon,
		Menu:        menu,
		Method:      method,
		Alias:       alias,
		Name:        name,
		ParentId:    parentId,
		Path:        path,
		Description: desc,
	})
}

func (s *service) Update(ctx context.Context, id int64, name, alias, icon, path, method, desc string, parentId int64, menu bool) (err error) {
	return s.repository.SysPermission().Save(ctx, &types.SysPermission{
		Id:          id,
		Icon:        icon,
		Menu:        menu,
		Method:      method,
		Alias:       alias,
		Name:        name,
		ParentId:    parentId,
		Path:        path,
		Description: desc,
	})
}

func (s *service) Delete(ctx context.Context, id int64) (err error) {
	// 如果它的上级有在用这个的话无法删除？还是说整个都删除掉？
	// 还是整个删除掉吧
	return s.repository.SysPermission().Delete(ctx, id)
}

func (s *service) List(ctx context.Context, page, pageSize int) (res []result, total int, err error) {
	panic("implement me")
}

func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "syspermission", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}