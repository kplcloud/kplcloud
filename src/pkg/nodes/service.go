/**
 * @Time : 8/11/21 11:43 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package nodes

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	Sync(ctx context.Context, clusterName string) (err error)
	List(ctx context.Context, clusterId int64, page, pageSize int) (res []nodeResult, total int, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func (s *service) List(ctx context.Context, clusterId int64, page, pageSize int) (res []nodeResult, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	list, total, err := s.repository.Nodes(ctx).List(ctx, clusterId, page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.Nodes", "List", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, nodeResult{
			Name:   v.Name,
			Memory: v.Memory,
			Cpu:    v.Cpu,
		})
	}

	return
}

func (s *service) Sync(ctx context.Context, clusterName string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	cluster, err := s.repository.Cluster(ctx).FindByName(ctx, clusterName)
	if err != nil {
		_ = level.Error(logger).Log("repository.Cluster", "FindByName", "err", err.Error())
		err = encode.ErrClusterNotfound.Error()
		return
	}

	if nodes, err := s.k8sClient.Do(ctx).CoreV1().Nodes().List(ctx, metav1.ListOptions{}); err == nil {
		for _, node := range nodes.Items {
			cpu, _ := node.Status.Capacity.Cpu().AsInt64()
			memory, _ := node.Status.Capacity.Memory().AsInt64()
			storage, _ := node.Status.Capacity.Storage().AsInt64()
			var internalIp, externalIp, status string
			for _, v := range node.Status.Addresses {
				if v.Type == v1.NodeInternalIP {
					internalIp = v.Address
					continue
				}
				if v.Type == v1.NodeExternalIP {
					externalIp = v.Address
					continue
				}
			}
			for _, v := range node.Status.Conditions {
				if v.Type == v1.NodeReady {
					status = string(v.Status)
				}
			}

			n, err := s.repository.Nodes(ctx).FindByName(ctx, cluster.Id, node.Name)
			if err != nil {
				if !gorm.IsRecordNotFoundError(err) {
					_ = level.Error(logger).Log("repository.Nodes", "FindByName", "err", err.Error())
					err = encode.ErrClusterNotfound.Error()
					return err
				}
				n = types.Nodes{}
			}
			n.ClusterId = cluster.Id
			n.Name = node.Name
			n.Memory = memory
			n.Cpu = cpu
			n.EphemeralStorage = storage
			n.InternalIp = internalIp
			n.ExternalIp = externalIp
			n.KubeletVersion = node.Status.NodeInfo.KubeletVersion
			n.KubeProxyVersion = node.Status.NodeInfo.KubeProxyVersion
			n.ContainerVersion = node.Status.NodeInfo.ContainerRuntimeVersion
			n.OsImage = node.Status.NodeInfo.OSImage
			n.Status = status
			n.Scheduled = !node.Spec.Unschedulable

			if err = s.repository.Nodes(ctx).Save(ctx, &n); err != nil {
				_ = level.Error(logger).Log("repository.Nodes", "Save", "err", err.Error())
			}
		}
	} else {
		_ = level.Error(logger).Log("k8sClient.Do.CoreV1.Nodes", "List", "err", err.Error())
	}

	return
}

func New(logger log.Logger, traceId string, client kubernetes.K8sClient, repository repository.Repository) Service {
	logger = log.With(logger, "nodes", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		k8sClient:  client,
		repository: repository,
	}
}
