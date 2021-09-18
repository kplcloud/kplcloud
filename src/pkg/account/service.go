/**
 * @Time : 2021/9/17 3:03 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package account

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	// UserInfo 获取用户详情包括角色权限、空间
	UserInfo(ctx context.Context, userId int64) (res userInfoResult, err error)
	// Menus 返回用户菜单
	Menus(ctx context.Context, userId int64) (res []userMenuResult, err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
}

func (s *service) Menus(ctx context.Context, userId int64) (res []userMenuResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	userInfo, err := s.repository.SysUser().Find(ctx, userId, "SysRoles.SysPermissions")
	if err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "Find", "err", err.Error())
		return
	}

	tmpPerm := map[string]types.SysPermission{}
	var menus []userMenuResult
	for _, v := range userInfo.SysRoles {
		if !v.Enabled {
			continue
		}
		for _, p := range v.SysPermissions {
			if !p.Menu {
				continue
			}
			if _, ok := tmpPerm[p.Name]; ok {
				continue
			}
			tmpPerm[p.Name] = p
			menus = append(menus, userMenuResult{
				Id:       p.Id,
				ParentId: p.ParentId,
				Icon:     p.Icon,
				Key:      p.Name,
				Text:     p.Alias,
				Link:     p.Path,
				Alias:    p.Alias,
			})
		}
	}

	res = permissionMenus(menus, 0)

	return
}

// GetMenu 获取菜单
func permissionMenus(permissions []userMenuResult, parentId int64) (menus []userMenuResult) {
	for _, v := range permissions {
		if v.ParentId == parentId {
			child := permissionMenus(permissions, v.Id)
			node := v
			node.Items = child
			menus = append(menus, node)
		}
	}
	return menus
}

func (s *service) UserInfo(ctx context.Context, userId int64) (res userInfoResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	userInfo, err := s.repository.SysUser().Find(ctx, userId, "SysRoles", "SysRoles.SysPermissions", "Clusters")
	if err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "Find", "err", err.Error())
		return
	}

	var permissions, roles, clusters []string
	for _, v := range userInfo.SysRoles {
		if !v.Enabled {
			continue
		}
		roles = append(roles, v.Name)
		for _, p := range v.SysPermissions {
			permissions = append(permissions, p.Name)
		}
	}
	for _, v := range userInfo.Clusters {
		clusters = append(clusters, v.Name)
	}
	res.Username = userInfo.Username
	res.Permissions = permissions
	res.Roles = roles
	res.Clusters = clusters

	return
}

func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "account", "service")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
	}
}
