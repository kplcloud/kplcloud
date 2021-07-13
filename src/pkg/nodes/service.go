/**
 * @Time : 6/15/21 4:58 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package nodes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
)

type Middleware func(Service) Service

type Service interface {
	List(ctx context.Context, nodeName string, page, pageSize int)
	Sync(ctx context.Context) (err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) List(ctx context.Context, nodeName string, page, pageSize int) {
	panic("implement me")
}

func (s *service) Sync(ctx context.Context) (err error) {
	logger := log.With(s.logger, "method", "Sync")

	nodes, err := s.k8sClient.Do().CoreV1().Nodes().List(metaV1.ListOptions{})
	if err != nil {
		_ = level.Error(logger).Log("k8sClient.Do.CoreV1.Nodes", "List", "err", err.Error())
		return
	}

	for _, node := range nodes.Items {
		fmt.Println(node.Name)
		b, _ := json.Marshal(node)
		fmt.Println(string(b))
	}

	return
}

func New(logger log.Logger, traceId string, repository repository.Repository, k8sClient kubernetes.K8sClient) Service {
	logger = log.With(logger, "service", "nodes")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
		k8sClient:  k8sClient,
	}
}
