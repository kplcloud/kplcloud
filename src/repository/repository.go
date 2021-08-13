/**
 * @Time: 2020/4/23 14:12
 * @Author: solacowa@gmail.com
 * @File: repository
 * @Software: GoLand
 */

package repository

import (
	"context"
	"github.com/go-kit/kit/log"
	kitcache "github.com/icowan/kit-cache"
	redisclient "github.com/icowan/redis-client"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/cluster"
	"github.com/kplcloud/kplcloud/src/repository/namespace"
	"github.com/kplcloud/kplcloud/src/repository/nodes"
	"github.com/opentracing/opentracing-go"

	"github.com/kplcloud/kplcloud/src/repository/sysnamespace"
	"github.com/kplcloud/kplcloud/src/repository/syspermission"
	"github.com/kplcloud/kplcloud/src/repository/sysrole"
	"github.com/kplcloud/kplcloud/src/repository/syssetting"
	"github.com/kplcloud/kplcloud/src/repository/sysuser"
)

type Repository interface {
	Cluster(ctx context.Context) cluster.Service
	Nodes(ctx context.Context) nodes.Service
	Namespace(ctx context.Context) namespace.Service

	SysSetting() syssetting.Service
	SysUser() sysuser.Service
	SysNamespace() sysnamespace.Service
	SysRole() sysrole.Service
	SysPermission() syspermission.Service
}

type repository struct {
	clusterSvc   cluster.Service
	nodesSvc     nodes.Service
	namespaceSvc namespace.Service

	sysSetting    syssetting.Service
	sysUser       sysuser.Service
	sysNamespace  sysnamespace.Service
	sysRole       sysrole.Service
	sysPermission syspermission.Service
}

func (r *repository) Namespace(ctx context.Context) namespace.Service {
	return r.namespaceSvc
}

func (r *repository) Nodes(ctx context.Context) nodes.Service {
	return r.nodesSvc
}

func (r *repository) Cluster(ctx context.Context) cluster.Service {
	return r.clusterSvc
}

func (r *repository) SysPermission() syspermission.Service {
	return r.sysPermission
}

func (r *repository) SysRole() sysrole.Service {
	return r.sysRole
}

func (r *repository) SysNamespace() sysnamespace.Service {
	return r.sysNamespace
}

func (r *repository) SysUser() sysuser.Service {
	return r.sysUser
}

func (r *repository) SysSetting() syssetting.Service {
	return r.sysSetting
}

func New(db *gorm.DB, logger log.Logger, traceId string, tracer opentracing.Tracer, redis redisclient.RedisClient, kcache kitcache.Service) Repository {
	// 平台系统相关仓库
	sysSetting := syssetting.New(db)
	//sysSetting = syssetting.NewLogging(logger, traceId)(sysSetting)

	sysNamespace := sysnamespace.New(db)
	sysNamespace = sysnamespace.NewLogging(logger, traceId)(sysNamespace)

	sysUser := sysuser.New(db)
	sysUser = sysuser.NewLogging(logger, traceId)(sysUser)

	sysRole := sysrole.New(db)
	sysRole = sysrole.NewLogging(logger, traceId)(sysRole)

	sysPermission := syspermission.New(db)
	sysPermission = syspermission.NewLogging(logger, traceId)(sysPermission)

	clusterSvc := cluster.New(db)
	clusterSvc = cluster.NewLogging(logger, traceId)(clusterSvc)
	nodesSvc := nodes.New(db)
	nodesSvc = nodes.NewLogging(logger, traceId)(nodesSvc)
	namespaceSvc := namespace.New(db)
	namespaceSvc = namespace.NewLogging(logger, traceId)(namespaceSvc)

	if tracer != nil {
		//sysSetting = syssetting.NewTracing(tracer)(sysSetting)
		sysUser = sysuser.NewTracing(tracer)(sysUser)
		sysNamespace = sysnamespace.NewTracing(tracer)(sysNamespace)
		sysRole = sysrole.NewTracing(tracer)(sysRole)
		//sysPermission = sysrole.NewTracing(tracer)(sysPermission)

		clusterSvc = cluster.NewTracing(tracer)(clusterSvc)
		nodesSvc = nodes.NewTracing(tracer)(nodesSvc)
		namespaceSvc = namespace.NewTracing(tracer)(namespaceSvc)
	}

	if redis != nil {
	}

	if kcache != nil {
		clusterSvc = cluster.NewCache(logger, traceId, kcache)(clusterSvc)
	}

	return &repository{
		sysSetting:    sysSetting,
		sysUser:       sysUser,
		sysNamespace:  sysNamespace,
		sysRole:       sysRole,
		sysPermission: sysPermission,

		clusterSvc:   clusterSvc,
		nodesSvc:     nodesSvc,
		namespaceSvc: namespaceSvc,
	}
}
