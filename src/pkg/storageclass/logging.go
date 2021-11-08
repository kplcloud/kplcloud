/**
 * @Time : 8/11/21 4:21 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package storageclass

import (
	"context"
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) List(ctx context.Context, clusterId int64, page, pageSize int) (res []listResult, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "clusterId", clusterId, "page", page, "pageSize", pageSize,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, clusterId, page, pageSize)
}

func (s *logging) Delete(ctx context.Context, clusterId int64, storageName string) (err error) {
	panic("implement me")
}

func (s *logging) Update(ctx context.Context, clusterId int64, storageName, provisioner string, reclaimPolicy *v1.PersistentVolumeReclaimPolicy, volumeBindingMode *storagev1.VolumeBindingMode) (err error) {
	panic("implement me")
}

func (s *logging) Create(ctx context.Context, clusterId int64, ns, name, provisioner string, reclaimPolicy *v1.PersistentVolumeReclaimPolicy, volumeBindingMode *storagev1.VolumeBindingMode) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create",
			"clusterId", clusterId,
			"ns", ns,
			"name", name,
			"provisioner", provisioner,
			"reclaimPolicy", reclaimPolicy,
			"volumeBindingMode", volumeBindingMode,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, clusterId, ns, name, provisioner, reclaimPolicy, volumeBindingMode)
}

func (s *logging) CreateProvisioner(ctx context.Context, clusterId int64) (err error) {
	panic("implement me")
}

func (s *logging) SyncPv(ctx context.Context, clusterId int64, storageName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SyncPv",
			"clusterId", clusterId,
			"storageName", storageName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.SyncPv(ctx, clusterId, storageName)
}

func (s *logging) SyncPvc(ctx context.Context, clusterId int64, ns string, storageName string) (err error) {
	panic("implement me")
}

func (s *logging) Sync(ctx context.Context, clusterId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Sync",
			"clusterId", clusterId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Sync(ctx, clusterId)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "storageclass", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
