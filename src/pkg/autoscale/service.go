/**
 * @Time: 2021/12/12 17:26
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package autoscale

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Middleware func(Service) Service

type Service interface {
	Sync(ctx context.Context, clusterId int64, namespace string) (err error)
}

type service struct {
	traceId    string
	logging    log.Logger
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) Sync(ctx context.Context, clusterId int64, namespace string) (err error) {
	list, err := s.k8sClient.Do(ctx).AutoscalingV1().HorizontalPodAutoscalers(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	for _, v := range list.Items {
		fmt.Println(v.Name)
	}
	return
}

func New(traceId string, logger log.Logger, store repository.Repository, k8sClient kubernetes.K8sClient) Service {
	logger = log.With(logger, "autoscale", "service")
	return &service{
		traceId:    traceId,
		logging:    logger,
		repository: store,
		k8sClient:  k8sClient,
	}
}
