/**
 * @Time: 2021/8/18 23:03
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package configmap

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
)

type Middleware func(Service) Service

type Service interface {
	Sync(ctx context.Context, ns string) (err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) Sync(ctx context.Context, ns string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	list, err := s.k8sClient.Do(ctx).CoreV1().ConfigMaps(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		_ = level.Error(logger).Log("k8sClient.Do.CoreV1.ConfigMaps", "List", "err", err.Error())
		return encode.ErrConfigMapSyncList.Wrap(err)
	}

	for _, v := range list.Items {
		b, _ := json.Marshal(v)
		fmt.Println(string(b))
	}

	return
}

func New(logger log.Logger, traceId string, repository repository.Repository, client kubernetes.K8sClient) Service {
	logger = log.With(logger, "configmap", "service")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
		k8sClient:  client,
	}
}
