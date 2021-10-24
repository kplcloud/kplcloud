/**
 * @Time : 3/5/21 5:32 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package sysuser

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	// List 系统用户列表
	List(ctx context.Context, email string, page, pageSize int) (res []listResult, total int, err error)
	// Add 添加系统用户
	Add(ctx context.Context, username, email, remark string, locked bool, namespaceIds, roleIds []int64) (err error)
	// Locked 锁定或解锁用户
	Locked(ctx context.Context, userId int64) (err error)
	// Delete 删除用户
	// unscoped: 是否硬删除 默认: false
	Delete(ctx context.Context, userId int64, unscoped bool) (err error)
	// Update 更新用户
	Update(ctx context.Context, userId int64, username, email, remark string, locked bool, clusterIds, roleIds []int64) (err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) Locked(ctx context.Context, userId int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Locked")
	sysUser, err := s.repository.SysUser().Find(ctx, userId)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "Find", "err", err.Error())
		err = encode.ErrSysUserNotfound.Error()
		return
	}

	sysUser.Locked = !sysUser.Locked
	return s.repository.SysUser().Save(ctx, &sysUser)
}

func (s *service) Delete(ctx context.Context, userId int64, unscoped bool) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Delete")

	err = s.repository.SysUser().Delete(ctx, userId, unscoped)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "Delete", "err", err.Error())
		err = encode.ErrSysUserNotfound.Error()
		return
	}

	return
}

func (s *service) Update(ctx context.Context, userId int64, username, email, remark string, locked bool, clusterIds, roleIds []int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Locked")
	sysUser, err := s.repository.SysUser().Find(ctx, userId)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "Find", "err", err.Error())
		err = encode.ErrSysUserNotfound.Error()
		return
	}
	var roles []types.SysRole
	var clusters []types.Cluster
	//var ns []types.SysNamespace
	if roles, err = s.repository.SysRole().FindByIds(ctx, roleIds); err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "FindByIds", "err", err.Error())
		return err
	}

	if clusters, err = s.repository.Cluster(ctx).FindByIds(ctx, clusterIds); err != nil {
		_ = level.Error(logger).Log("repository.Cluster", "FindByIds", "err", err.Error())
		return err
	}

	sysUser.Locked = locked
	sysUser.Username = username
	sysUser.Email = email
	sysUser.SysRoles = roles
	sysUser.Clusters = clusters
	sysUser.Remark = remark
	return s.repository.SysUser().Save(ctx, &sysUser)
}

func (s *service) Add(ctx context.Context, username, email, remark string, locked bool, namespaceIds, roleIds []int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Add")
	roles, err := s.repository.SysRole().FindByIds(ctx, roleIds)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "FindByIds", "err", err.Error())
		return
	}
	//namespaces, err := s.repository.SysNamespace().FindByIds(ctx, namespaceIds)
	//if err != nil {
	//	_ = level.Error(logger).Log("repository.SysRole", "FindByIds", "err", err.Error())
	//	return
	//}

	if err = s.repository.SysUser().Save(ctx, &types.SysUser{
		Username:  username,
		LoginName: username,
		Email:     email,
		Remark:    remark,
		Locked:    locked,
		SysRoles:  roles,
		//SysNamespaces: namespaces,
	}); err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "Save", "err", err.Error())
		return err
	}

	return
}

func (s *service) List(ctx context.Context, email string, page, pageSize int) (res []listResult, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "List")

	list, total, err := s.repository.SysUser().List(ctx, nil, email, page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "List", "err", err.Error())
		return
	}

	for _, v := range list {
		var roles []string
		for _, vv := range v.SysRoles {
			roles = append(roles, vv.Alias)
		}
		res = append(res, listResult{
			Username:     v.Username,
			Email:        v.Email,
			Locked:       v.Locked,
			WechatOpenId: v.WechatOpenId,
			LastLogin:    v.LastLogin,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
			ExpiresAt:    v.ExpiresAt,
			Remark:       v.Remark,
			Namespaces:   []string{"default"},
			Roles:        roles,
		})
	}

	return
}

func New(logger log.Logger, traceId string, store repository.Repository) Service {
	logger = log.With(logger, "sysuser", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: store,
	}
}
