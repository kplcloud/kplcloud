/**
 * @Time : 3/10/21 3:27 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package sysrole

import (
	"context"
	"github.com/kplcloud/kplcloud/src/encode"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(next Service) Service

type Service interface {
	// List 角色列表
	List(ctx context.Context, page, pageSize int) (res []listResult, total int, err error)
	// Add 添加角色
	Add(ctx context.Context, alias, name, description string, enabled bool) (err error)
	// Permissions 角色所拥有的权限列表
	Permissions(ctx context.Context, id int64, page, pageSize int) (res []listResult, total int, err error)
	// Permission 设置角色权限
	Permission(ctx context.Context, id int64, permIds []int64) (err error)
	// 角色下的用户
	//Users(ctx context.Context, id int64, page, pageSize int) (res []listResult, total int, err error)

	// User 设置用户角色
	User(ctx context.Context, id int64, userIds []int64) (err error)
	// Update 更新角色信息
	Update(ctx context.Context, id int64, alias, name, description string, enabled bool) (err error)
	// Delete 删除角色
	Delete(ctx context.Context, id int64) (err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) Delete(ctx context.Context, id int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Update")
	if err = s.repository.SysRole().Delete(ctx, id); err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "Delete", "err", err.Error())
		return encode.ErrSysRoleUserDelete.Error()
	}
	return
}

func (s *service) Update(ctx context.Context, id int64, alias, name, description string, enabled bool) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Update")

	role, err := s.repository.SysRole().Find(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "Find", "err", err.Error())
		return encode.ErrSysRoleNotfound.Error()
	}
	role.Alias = alias
	role.Name = name
	role.Description = description
	role.Enabled = enabled
	if err = s.repository.SysRole().Save(ctx, &role); err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "Save", "err", err.Error())
		return encode.ErrSysRoleSave.Error()
	}
	return
}

func (s *service) User(ctx context.Context, id int64, userIds []int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "User")

	role, err := s.repository.SysRole().Find(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "Find", "err", err.Error())
		return encode.ErrSysRoleNotfound.Error()
	}

	sysUsers, err := s.repository.SysUser().FindByIds(ctx, userIds, "SysRoles")
	if err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "FindByIds", "err", err.Error())
		return encode.ErrSysRoleUserNotfound.Error()
	}

	for _, sysUser := range sysUsers {
		roles := sysUser.SysRoles
		roles = append(roles, role)
		if err = s.repository.SysUser().AddRoles(ctx, &sysUser, roles); err != nil {
			_ = level.Error(logger).Log("repository.SysUser", "AddRoles", "err", err.Error())
			return encode.ErrSysRoleUser.Error()
		}
	}

	return
}

func (s *service) Permission(ctx context.Context, id int64, permIds []int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Permission")

	role, err := s.repository.SysRole().Find(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "Find", "err", err.Error())
		return encode.ErrSysRoleNotfound.Error()
	}

	perms, err := s.repository.SysPermission().FindByIds(ctx, permIds)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysPermission", "FindByIds", "err", err.Error())
		return err
	}

	if err = s.repository.SysRole().AddPermissions(ctx, &role, perms); err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "Save", "err", err.Error())
		return encode.ErrSysRoleSave.Error()
	}

	// rds

	return
}

func (s *service) Permissions(ctx context.Context, id int64, page, pageSize int) (res []listResult, total int, err error) {
	panic("implement me")
}

func (s *service) Add(ctx context.Context, alias, name, description string, enabled bool) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "List")

	if err = s.repository.SysRole().Save(ctx, &types.SysRole{
		Alias:       alias,
		Name:        name,
		Enabled:     enabled,
		Description: description,
	}); err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "Save", "err", err.Error())
	}
	return
}

func (s *service) List(ctx context.Context, page, pageSize int) (res []listResult, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "List")
	list, total, err := s.repository.SysRole().List(ctx, page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "List", "err", err.Error())
		return
	}
	for _, v := range list {
		res = append(res, listResult{
			Id:          v.Id,
			Alias:       v.Alias,
			Name:        v.Name,
			Enabled:     v.Enabled,
			Description: v.Description,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}
	return
}

func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "sysrole", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}
