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
)

type Middleware func(Service) Service

type Service interface {
	// UserInfo 获取用户详情包括角色权限、空间
	UserInfo(ctx context.Context, userId int64) (res userInfoResult, err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
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
