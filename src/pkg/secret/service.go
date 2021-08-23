/**
 * @Time : 8/19/21 1:36 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package secret

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Middleware func(Service) Service

type Service interface {
	Sync(ctx context.Context, clusterId int64, ns string) (err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	var items *corev1.SecretList
	if items, err = s.k8sClient.Do(ctx).CoreV1().Secrets(ns).List(ctx, metav1.ListOptions{}); err != nil {
		_ = level.Error(logger).Log("k8sClient.Do.AppsV1.Secrets", "List", "err", err.Error())
		return encode.ErrDeploymentSyncList.Wrap(err)
	}

	for _, v := range items.Items {
		var data []types.Data
		for key, val := range v.Data {
			data = append(data, types.Data{
				Style: types.DataStyleSecret,
				Key:   key,
				Value: string(val),
			})
		}

		if err = s.repository.Secrets(ctx).Save(ctx, &types.Secret{
			ClusterId:       clusterId,
			Name:            v.Name,
			Namespace:       v.Namespace,
			ResourceVersion: v.ResourceVersion,
		}, data); err != nil {
			_ = level.Error(logger).Log("repository.Secrets", "Save", "err", err.Error())
		}
	}

	return nil
}

func New(logger log.Logger, traceId string, repository repository.Repository, client kubernetes.K8sClient) Service {
	logger = log.With(logger, "secret", "service")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
		k8sClient:  client,
	}
}
