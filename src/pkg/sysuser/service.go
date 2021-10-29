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
	// clusterIds, namespaceIds, roleIds 需要过滤，不能高过自己所拥有 用name,不用ID
	Add(ctx context.Context, username, email, remark string, locked bool, clusterIds, namespaceIds, roleIds []int64) (err error)
	// Locked 锁定或解锁用户
	Locked(ctx context.Context, userId int64) (err error)
	// Delete 删除用户
	// unscoped: 是否硬删除 默认: false
	Delete(ctx context.Context, userId int64, unscoped bool) (err error)
	// Update 更新用户
	Update(ctx context.Context, userId int64, username, email, remark string, locked bool, clusterIds, roleIds []int64) (err error)
	// GetRoles 获取当前用户的角色
	// names: 从中间件获取
	// 只会返回比当前用户所拥有的最高权限等级更低的角色 TODO: 管理员取所有
	GetRoles(ctx context.Context, sysUserId int64, names []string) (res []roleResult, err error)
	// GetCluster 获取当前用户拥有的集群
	// clusterNames: 从中间件获取 // TODO: 管理员取所有
	GetCluster(ctx context.Context, sysUserId int64, clusterNames []string) (res []clusterResult, err error)
	// GetNamespaces 获取当前用户可以操作的namespaces
	// clusterNames: 前端传过来，但需要在中间件进行校验 // TODO: 管理员取所有
	GetNamespaces(ctx context.Context, sysUserId int64, clusterNames []string) (res []namespaceResult, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) GetRoles(ctx context.Context, sysUserId int64, names []string) (res []roleResult, err error) {
	list, err := s.repository.SysRole().FindByNames(ctx, names)
	if err != nil {
		return
	}
	lv := list[0].Level
	list, err = s.repository.SysRole().FindByLevel(ctx, lv, "<=")
	if err != nil {
		return
	}

	for _, v := range list {
		if !v.Enabled {
			continue
		}
		res = append(res, roleResult{
			Alias:       v.Alias,
			Name:        v.Name,
			Enabled:     v.Enabled,
			Description: v.Description,
		})
	}
	return
}

func (s *service) GetCluster(ctx context.Context, sysUserId int64, clusterNames []string) (res []clusterResult, err error) {
	list, err := s.repository.Cluster(ctx).FindByNames(ctx, clusterNames)
	if err != nil {
		return
	}
	for _, v := range list {
		res = append(res, clusterResult{
			Name:   v.Name,
			Alias:  v.Alias,
			Remark: v.Remark,
		})
	}
	return
}

func (s *service) GetNamespaces(ctx context.Context, sysUserId int64, clusterNames []string) (res []namespaceResult, err error) {
	list, err := s.repository.Cluster(ctx).FindByNames(ctx, clusterNames)
	if err != nil {
		return
	}
	var ids []int64
	for _, v := range list {
		ids = append(ids, v.Id)
	}
	namespaces, err := s.repository.Namespace(ctx).FindByIds(ctx, ids)
	if err != nil {
		return
	}

	for _, v := range namespaces {
		res = append(res, namespaceResult{
			Name:   v.Name,
			Alias:  v.Alias,
			Remark: v.Remark,
		})
	}

	return
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

func (s *service) Add(ctx context.Context, username, email, remark string, locked bool, clusterIds, namespaceIds, roleIds []int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Add")
	roles, err := s.repository.SysRole().FindByIds(ctx, roleIds)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "FindByIds", "err", err.Error())
		return
	}
	clusters, err := s.repository.Cluster(ctx).FindByIds(ctx, clusterIds)
	if err != nil {
		_ = level.Error(logger).Log("repository.Cluster", "FindByIds", "err", err.Error())
		return
	}
	namespaces, err := s.repository.Namespace(ctx).FindByIds(ctx, namespaceIds)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysRole", "FindByIds", "err", err.Error())
		return
	}

	if err = s.repository.SysUser().Save(ctx, &types.SysUser{
		Username:   username,
		LoginName:  username,
		Email:      email,
		Remark:     remark,
		Locked:     locked,
		SysRoles:   roles,
		Namespaces: namespaces,
		Clusters:   clusters,
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
