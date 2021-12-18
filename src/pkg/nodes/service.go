/**
 * @Time : 8/11/21 11:43 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package nodes

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"

	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

// Service 集群Node节点模块
type Service interface {
	// Sync 同步节点信息
	Sync(ctx context.Context, clusterName string) (err error)
	// List 节点列表
	List(ctx context.Context, clusterId int64, query string, page, pageSize int) (res []nodeResult, total int, err error)
	// Cordon 将节点设置为可调度或不可调度
	Cordon(ctx context.Context, clusterId int64, nodeName string) (err error)
	// Drain 驱逐节点上有pods nodeName 节点名称 force 强制
	Drain(ctx context.Context, clusterId int64, nodeName string, force bool) (err error)
	// Info 节点详情
	Info(ctx context.Context, clusterId int64, nodeName string) (res infoResult, err error)
	// Delete 删除集群节点
	Delete(ctx context.Context, clusterId int64, nodeName string) (err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func (s *service) Delete(ctx context.Context, clusterId int64, nodeName string) (err error) {
	//logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	resNode, err := s.repository.Nodes(ctx).FindByName(ctx, clusterId, nodeName)
	if err != nil {
		err = errors.Wrap(err, "repository.Nodes.FindByName")
		return encode.ErrNodeNotfound.Wrap(err)
	}

	if resNode.Scheduled {
		// 请先设置为不可调度
		return encode.ErrNodeScheduled.Error()
	}
	// TODO: 请把该节点pods驱逐完 不一定能驱逐掉

	err = s.repository.Nodes(ctx).Delete(ctx, clusterId, resNode.Name, func() error {
		if err = s.k8sClient.Do(ctx).CoreV1().Nodes().Delete(ctx, nodeName, metav1.DeleteOptions{}); err != nil {
			err = errors.Wrap(err, "pkg.nodes.Delete.k8sClient.Do.CoreV2.Nodes.Delete")
			return encode.ErrNodeDelete.Wrap(err)
		}
		return nil
	})

	return
}

func (s *service) Info(ctx context.Context, clusterId int64, nodeName string) (res infoResult, err error) {
	//logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	resNode, err := s.repository.Nodes(ctx).FindByName(ctx, clusterId, nodeName)
	if err != nil {
		err = encode.ErrNodeNotfound.Wrap(errors.Wrap(err, "repository.Nodes.FindByName"))
		return
	}

	cluster, err := s.repository.Cluster(ctx).Find(ctx, clusterId)
	if err != nil {
		err = encode.ErrClusterNotfound.Wrap(errors.Wrap(err, "repository.Cluster.Find"))
		return
	}

	node, err := s.k8sClient.Do(ctx).CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		err = encode.ErrNodeCordon.Wrap(errors.Wrap(err, "k8sClient.Do.CoreV1.Nodes.Get"))
		return
	}
	res.Name = resNode.Name
	res.Remark = resNode.Remark
	res.InternalIp = resNode.InternalIp
	res.Scheduled = resNode.Scheduled
	res.Status = resNode.Status
	res.ExternalIp = resNode.ExternalIp
	res.OsImage = resNode.OsImage
	res.Labels = node.Labels
	res.Annotations = node.Annotations
	res.UsedCPU = fmt.Sprintf("%.2f", node.Status.Capacity.Cpu().AsApproximateFloat64()-node.Status.Allocatable.Cpu().AsApproximateFloat64())
	res.UsedMemory = apiresource.NewQuantity(node.Status.Capacity.Memory().Value()-node.Status.Allocatable.Memory().Value(), apiresource.BinarySI).String()
	res.CPU = node.Status.Capacity.Cpu().String()
	res.Memory = node.Status.Capacity.Memory().String()
	res.KubeProxyVersion = resNode.KubeProxyVersion
	res.KubeletVersion = resNode.KubeletVersion
	//res.SystemDisk = fmt.Sprintf("%dG", apiresource.NewQuantity(node.Status.Capacity.StorageEphemeral().Value(), apiresource.DecimalSI).ScaledValue(apiresource.Giga))
	res.SystemDisk = node.Status.Capacity.StorageEphemeral().String()
	res.Bandwidth = ""
	res.PodNum = node.Status.Capacity.Pods().Value()
	res.ClusterAlias = cluster.Alias
	res.ClusterName = cluster.Name
	res.ClusterStatus = cluster.Status
	res.ContainerVersion = resNode.ContainerVersion
	return
}

func (s *service) Drain(ctx context.Context, clusterId int64, nodeName string, force bool) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	resNode, err := s.repository.Nodes(ctx).FindByName(ctx, clusterId, nodeName)
	if err != nil {
		err = errors.Wrap(err, "repository.Nodes.FindByName")
		return encode.ErrNodeNotfound.Wrap(err)
	}
	//labelSelector, err := labels.Parse(d.PodSelector)
	//if err != nil {
	//	return encode.ErrNodeDrain.Wrap(errors.Wrap(err, "labels.Parse"))
	//}
	if resNode.Scheduled {
		// 请先设置为不可调度
		return encode.ErrNodeScheduled.Error()
	}
	pods, err := s.k8sClient.Do(ctx).CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{
		//LabelSelector: labelSelector.String(),
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": nodeName}).String(),
	})
	if err != nil {
		err = errors.Wrap(err, "k8sClient.Do.CoreV1.Pods.List")
		return encode.ErrNodeCordon.Wrap(err)
	}

	fmt.Println(resNode.Name)

	for _, pod := range pods.Items {
		err := s.k8sClient.Do(ctx).CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{
			DryRun: []string{metav1.DryRunAll},
		})
		if err != nil {
			// TODO 删除失败的记录到某个地方
			_ = level.Error(logger).Log("k8sClient.Do", "CoreV1.Pods", "Delete", pod.Name, "err", err.Error())
		}
	}

	return
}

func (s *service) Cordon(ctx context.Context, clusterId int64, nodeName string) (err error) {
	resNode, err := s.repository.Nodes(ctx).FindByName(ctx, clusterId, nodeName)
	if err != nil {
		err = errors.Wrap(err, "repository.Nodes.FindByName")
		return encode.ErrNodeNotfound.Wrap(err)
	}

	node, err := s.k8sClient.Do(ctx).CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, "k8sClient.Do.CoreV1.Nodes.Get")
		return encode.ErrNodeCordon.Wrap(err)
	}
	node.Spec.Unschedulable = !node.Spec.Unschedulable
	node, err = s.k8sClient.Do(ctx).CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		err = errors.Wrap(err, "k8sClient.Do.CoreV1.Nodes.Update")
		return encode.ErrNodeCordon.Wrap(err)
	}
	resNode.Scheduled = !node.Spec.Unschedulable
	return s.repository.Nodes(ctx).Save(ctx, &resNode)
}

func (s *service) List(ctx context.Context, clusterId int64, query string, page, pageSize int) (res []nodeResult, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	list, total, err := s.repository.Nodes(ctx).List(ctx, clusterId, query, page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.Nodes", "List", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, nodeResult{
			Name:             v.Name,
			Memory:           fmt.Sprintf("%dGi", apiresource.NewQuantity(v.Memory, apiresource.BinarySI).ScaledValue(apiresource.Giga)),
			Cpu:              v.Cpu,
			EphemeralStorage: apiresource.NewQuantity(v.EphemeralStorage, apiresource.BinarySI).String(),
			InternalIp:       v.InternalIp,
			ExternalIp:       v.ExternalIp,
			KubeletVersion:   v.KubeletVersion,
			KubeProxyVersion: v.KubeProxyVersion,
			ContainerVersion: v.ContainerVersion,
			OsImage:          v.OsImage,
			Status:           v.Status,
			Scheduled:        v.Scheduled,
			Remark:           v.Remark,
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
		cluster.NodeNum = len(nodes.Items)
		if err := s.repository.Cluster(ctx).Save(ctx, &cluster); err != nil {
			_ = level.Error(logger).Log("repository.Cluster", "Save", "err", err.Error())
		}

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
