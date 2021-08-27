/**
 * @Time : 2021/8/27 2:56 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package cronjob

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kit/kit/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
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

	list, err := s.k8sClient.Do(ctx).BatchV1beta1().CronJobs(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		b, _ := json.Marshal(item)
		fmt.Println(string(b))
	}
	return
}

func New(logger log.Logger, traceId string, repository repository.Repository, client kubernetes.K8sClient) Service {
	logger = log.With(logger, "cronjob", "service")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
		k8sClient:  client,
	}
}
