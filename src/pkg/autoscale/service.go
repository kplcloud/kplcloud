/**
 * @Time: 2021/12/12 17:26
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package autoscale

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	autoscalev1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Middleware func(Service) Service

// Service 自动水平伸缩依赖metrics-server
// 使用该功能需要安装metrics-server
// metrics-server command参考
//  name: metrics-server
//	image: 'k8s.gcr.io/metrics-server/metrics-server:v0.5.2'
//	args:
//	- '--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname'
//	- '--kubelet-use-node-status-port'
//	- '--metric-resolution=15s'
//	- '--kubelet-insecure-tls'
type Service interface {
	Sync(ctx context.Context, clusterId int64, namespace string) (err error)
	// Create 创建HPA
	// 如果 kind 的资源没有设置 request 和 limit 则创建的hpa可能会不生效
	// 如果没有则给返回无法创建，需要先设置 request和limit
	Create(ctx context.Context, clusterId int64, namespace, name, kind, appName string) (err error)
	// Delete 删除自动伸缩
	Delete(ctx context.Context, clusterId int64, namespace, name string) (err error)
	// List 列表
	List(ctx context.Context, clusterId int64, namespace string, page, pageSize int) (res []result, total int, err error)
	// Detail 获取自动伸缩详情
	Detail(ctx context.Context, clusterId int64, namespace, name string) (res result, err error)
}

type service struct {
	traceId    string
	logging    log.Logger
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) Delete(ctx context.Context, clusterId int64, namespace, name string) (err error) {
	panic("implement me")
}

func (s *service) List(ctx context.Context, clusterId int64, namespace string, page, pageSize int) (res []result, total int, err error) {
	// TODO: 先查询group是否对该应用可读权限
	s.repository.HPA(ctx).List(ctx, clusterId, namespace, nil, page, pageSize)
	return
}

func (s *service) Detail(ctx context.Context, clusterId int64, namespace, name string) (res result, err error) {
	panic("implement me")
}

func (s *service) Create(ctx context.Context, clusterId int64, namespace, name, kind, appName string) (err error) {
	var hpa *autoscalev1.HorizontalPodAutoscaler
	_, err = s.repository.K8sTpl(ctx).EncodeTemplate(ctx, types.KindHorizontalPodAutoscaler, map[string]interface{}{
		"name":         name,
		"namespace":    namespace,
		"labelAppName": types.LabelAppName,
		"kind":         "Deployment",
		"apiVersion":   "apps/v1",
		"minReplicas":  1,
		"maxReplicas":  3,
		"resourceName": "cpu",
		"target":       50,
	}, &hpa)

	s.k8sClient.Do(ctx).AutoscalingV1().HorizontalPodAutoscalers(namespace).Create(ctx, hpa, metav1.CreateOptions{})
	return err
}

func (s *service) Sync(ctx context.Context, clusterId int64, namespace string) (err error) {
	logger := log.With(s.logging, s.traceId, ctx.Value(s.traceId))
	list, err := s.k8sClient.Do(ctx).AutoscalingV1().HorizontalPodAutoscalers(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	for _, v := range list.Items {
		var minReplicas = 1
		if v.Spec.MinReplicas != nil {
			minReplicas = int(*v.Spec.MinReplicas)
		}
		var targetCpu = 50
		if v.Spec.TargetCPUUtilizationPercentage != nil {
			targetCpu = int(*v.Spec.TargetCPUUtilizationPercentage)
		}
		e := s.repository.HPA(ctx).Save(ctx, &types.HorizontalPodAutoscaler{
			ClusterId:                clusterId,
			Namespace:                namespace,
			Name:                     v.Name,
			AppName:                  v.Name,
			ApiVersion:               v.APIVersion,
			MinReplicas:              minReplicas,
			MaxReplicas:              int(v.Spec.MaxReplicas),
			ResourceName:             "cpu",
			TargetAverageUtilization: targetCpu,
			Kind:                     types.Kind(v.Spec.ScaleTargetRef.Kind),
		})
		if e != nil {
			_ = level.Error(logger).Log("repository.Hpa", "Save", "err", e.Error())
			continue
		}
	}
	return
}

func New(logger log.Logger, traceId string, store repository.Repository, k8sClient kubernetes.K8sClient) Service {
	logger = log.With(logger, "autoscale", "service")
	return &service{
		traceId:    traceId,
		logging:    logger,
		repository: store,
		k8sClient:  k8sClient,
	}
}
