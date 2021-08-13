package namespace

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	// 同步namespace
	Sync(ctx context.Context, clusterId int64) error
}

type service struct {
	traceId    string
	logger     log.Logger
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func (s *service) Sync(ctx context.Context, clusterId int64) error {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	if namespaces, err := s.k8sClient.Do(ctx).CoreV1().Namespaces().List(ctx, metav1.ListOptions{}); err == nil {
		for _, v := range namespaces.Items {
			ns, err := s.repository.Namespace(ctx).FindByName(ctx, clusterId, v.Name)
			if err != nil {
				if !gorm.IsRecordNotFoundError(err) {
					_ = level.Error(logger).Log("repository.Namespace", "FindByName", "err", err.Error())
					err = encode.ErrNamespaceNotfound.Error()
					return err
				}
				ns = types.Namespace{}
			}
			ns.ClusterId = clusterId
			ns.Name = v.Name
			ns.Alias = v.Name
			ns.Status = string(v.Status.Phase)
			if err = s.repository.Namespace(ctx).Save(ctx, &ns); err != nil {
				_ = level.Error(logger).Log("repository.Namespace", "Save", "err", err.Error())
			}
		}
	} else {
		_ = level.Error(logger).Log("k8sClient.Do.CoreV1.Namespaces", "List", "err", err.Error())
	}

	return nil
}

func New(logger log.Logger, traceId string, client kubernetes.K8sClient, store repository.Repository) Service {
	return &service{
		traceId:    traceId,
		logger:     logger,
		k8sClient:  client,
		repository: store,
	}
}
