/**
 * @Time : 8/9/21 6:20 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package cluster

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	Add(ctx context.Context, name, alias, data string) (err error)
	//List(ctx context.Context, name string, page, pageSize int) (res )
}

type service struct {
	k8sClient  kubernetes.K8sClient
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) Add(ctx context.Context, name, alias, data string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	cluster := types.Cluster{
		Name:       name,
		Alias:      alias,
		Status:     1,
		ConfigData: data,
	}

	if err = s.repository.Cluster(ctx).Save(ctx, &cluster, func(tx *gorm.DB) error {
		if err = s.k8sClient.Connect(ctx, name, data); err != nil {
			_ = level.Error(logger).Log("k8sClient.Connect", "err", err.Error())
			return encode.ErrClusterConnect.Error()
		}
		return nil
	}); err != nil {
		_ = level.Error(logger).Log("repository.Cluster", "Save", "err", err.Error())
		return encode.ErrClusterAdd.Error()
	}

	return
}

func New(logger log.Logger, traceId string, repository repository.Repository, k8sClient kubernetes.K8sClient) Service {
	logger = log.With(logger, "cluster", "service")
	return &service{
		k8sClient:  k8sClient,
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}
