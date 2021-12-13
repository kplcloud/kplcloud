/**
 * @Time: 2020/4/23 14:12
 * @Author: solacowa@gmail.com
 * @File: repository
 * @Software: GoLand
 */

package repository

import (
	"context"
	"github.com/kplcloud/kplcloud/src/repository/hpa"
	"github.com/kplcloud/kplcloud/src/repository/pvc"

	"github.com/go-kit/kit/log"
	kitcache "github.com/icowan/kit-cache"
	redisclient "github.com/icowan/redis-client"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"

	"github.com/kplcloud/kplcloud/src/repository/application"
	"github.com/kplcloud/kplcloud/src/repository/audit"
	"github.com/kplcloud/kplcloud/src/repository/cluster"
	"github.com/kplcloud/kplcloud/src/repository/configmap"
	"github.com/kplcloud/kplcloud/src/repository/k8stpl"
	"github.com/kplcloud/kplcloud/src/repository/namespace"
	"github.com/kplcloud/kplcloud/src/repository/nodes"
	"github.com/kplcloud/kplcloud/src/repository/registry"
	"github.com/kplcloud/kplcloud/src/repository/secrets"
	"github.com/kplcloud/kplcloud/src/repository/storageclass"
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
	ConfigMap(ctx context.Context) configmap.Service
	Secrets(ctx context.Context) secrets.Service
	StorageClass(ctx context.Context) storageclass.Service
	K8sTpl(ctx context.Context) k8stpl.Service
	Registry(ctx context.Context) registry.Service
	Audit(ctx context.Context) audit.Service
	Application(ctx context.Context) application.Service
	Pvc(ctx context.Context) pvc.Service
	HPA(ctx context.Context) hpa.Service

	SysSetting() syssetting.Service
	SysUser() sysuser.Service
	SysNamespace() sysnamespace.Service
	SysRole() sysrole.Service
	SysPermission() syspermission.Service

	// old
	Groups() GroupsRepository
	CronJob() CronjobRepository
	Project() ProjectRepository
}

type repository struct {
	clusterSvc      cluster.Service
	nodesSvc        nodes.Service
	namespaceSvc    namespace.Service
	configMapSvc    configmap.Service
	secretSvc       secrets.Service
	storageClassSvc storageclass.Service
	k8sTpl          k8stpl.Service
	registrySvc     registry.Service
	auditSvc        audit.Service
	appSvc          application.Service
	pvcSvc          pvc.Service
	hpaSvc          hpa.Service

	sysSetting    syssetting.Service
	sysUser       sysuser.Service
	sysNamespace  sysnamespace.Service
	sysRole       sysrole.Service
	sysPermission syspermission.Service
}

func (r *repository) HPA(ctx context.Context) hpa.Service {
	return r.hpaSvc
}

func (r *repository) Pvc(ctx context.Context) pvc.Service {
	return r.pvcSvc
}

func (r *repository) Application(ctx context.Context) application.Service {
	return r.appSvc
}

func (r *repository) Audit(ctx context.Context) audit.Service {
	return r.auditSvc
}

func (r *repository) Registry(ctx context.Context) registry.Service {
	return r.registrySvc
}

func (r *repository) K8sTpl(ctx context.Context) k8stpl.Service {
	return r.k8sTpl
}

func (r *repository) StorageClass(ctx context.Context) storageclass.Service {
	return r.storageClassSvc
}

func (r *repository) Secrets(ctx context.Context) secrets.Service {
	return r.secretSvc
}

func (r *repository) ConfigMap(ctx context.Context) configmap.Service {
	return r.configMapSvc
}

func (r *repository) Project() ProjectRepository {
	panic("implement me")
}

func (r *repository) CronJob() CronjobRepository {
	panic("implement me")
}

func (r *repository) Groups() GroupsRepository {
	panic("implement me")
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
	sysSetting = syssetting.NewLogging(logger, traceId)(sysSetting)

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
	configMapSvc := configmap.New(db)
	configMapSvc = configmap.NewLogging(logger, traceId)(configMapSvc)
	secretSvc := secrets.New(db)
	secretSvc = secrets.NewLogging(logger, traceId)(secretSvc)
	storageClassSvc := storageclass.New(db)
	storageClassSvc = storageclass.NewLogging(logger, traceId)(storageClassSvc)
	k8sTplSvc := k8stpl.New(db)
	k8sTplSvc = k8stpl.NewLogging(logger, traceId)(k8sTplSvc)
	registrySvc := registry.New(db)
	registrySvc = registry.NewLogging(logger, traceId)(registrySvc)
	auditSvc := audit.New(db)
	auditSvc = audit.NewLogging(logger, traceId)(auditSvc)
	appSvc := application.New(db)
	appSvc = application.NewLogging(logger, traceId)(appSvc)
	pvcSvc := pvc.New(db)
	pvcSvc = pvc.NewLogging(logger, traceId)(pvcSvc)
	hpaSvc := hpa.New(db)
	hpaSvc = hpa.NewLogging(logger, traceId)(hpaSvc)

	if tracer != nil {
		sysSetting = syssetting.NewTracing(tracer)(sysSetting)
		sysUser = sysuser.NewTracing(tracer)(sysUser)
		sysNamespace = sysnamespace.NewTracing(tracer)(sysNamespace)
		sysRole = sysrole.NewTracing(tracer)(sysRole)
		//sysPermission = sysrole.NewTracing(tracer)(sysPermission)

		clusterSvc = cluster.NewTracing(tracer)(clusterSvc)
		nodesSvc = nodes.NewTracing(tracer)(nodesSvc)
		namespaceSvc = namespace.NewTracing(tracer)(namespaceSvc)
		configMapSvc = configmap.NewTracing(tracer)(configMapSvc)
		secretSvc = secrets.NewTracing(tracer)(secretSvc)
		storageClassSvc = storageclass.NewTracing(tracer)(storageClassSvc)
		registrySvc = registry.NewTracing(tracer)(registrySvc)
		k8sTplSvc = k8stpl.NewTracing(tracer)(k8sTplSvc)
		auditSvc = audit.NewTracing(tracer)(auditSvc)
		appSvc = application.NewTracing(tracer)(appSvc)
		pvcSvc = pvc.NewTracing(tracer)(pvcSvc)
		hpaSvc = hpa.NewTracing(tracer)(hpaSvc)
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

		storageClassSvc: storageClassSvc,
		k8sTpl:          k8sTplSvc,
		clusterSvc:      clusterSvc,
		nodesSvc:        nodesSvc,
		namespaceSvc:    namespaceSvc,
		configMapSvc:    configMapSvc,
		secretSvc:       secretSvc,
		registrySvc:     registrySvc,
		auditSvc:        auditSvc,
		appSvc:          appSvc,
		pvcSvc:          pvcSvc,
		hpaSvc:          hpaSvc,
	}
}
