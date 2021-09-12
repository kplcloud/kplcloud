package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type Service interface {
	Sync(ctx context.Context, clusterId int64, ns string) (err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	list, err := s.k8sClient.Do(ctx).CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		ports, _ := json.Marshal(item.Spec.Ports)
		selector, _ := json.Marshal(item.Spec.Selector)
		if item.Spec.Selector != nil {
			for _, val := range item.Spec.Selector {
				if strings.EqualFold(val, item.Name) {
					get, err := s.k8sClient.Do(ctx).CoreV1().Endpoints(ns).Get(ctx, item.Name, metav1.GetOptions{})
					if err != nil {
						return err
					}
					fmt.Println(get)
					break
				}
			}
		}

		svc := types.Service{
			ClusterId:   clusterId,
			Namespace:   item.Namespace,
			Name:        item.Name,
			Ports:       string(ports),
			Selector:    string(selector),
			ServiceType: string(item.Spec.Type),
		}
		fmt.Println(svc)
	}

	return
}

func New(logger log.Logger, traceId string, repository repository.Repository, client kubernetes.K8sClient) Service {
	logger = log.With(logger, "service", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
		k8sClient:  client,
	}
}
