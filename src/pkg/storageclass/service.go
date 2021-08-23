/**
 * @Time : 2021/8/23 10:07 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package storageclass

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"

	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	Sync(ctx context.Context, clusterId int64) (err error)
	SyncPv(ctx context.Context, clusterId int64, storageName string) (err error)
	SyncPvc(ctx context.Context, clusterId int64, ns string, storageName string) (err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

type persistentVolumeListChannel struct {
	List  chan *coreV1.PersistentVolumeList
	Error chan error
}

func (s *service) SyncPv(ctx context.Context, clusterId int64, storageName string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	find, err := s.repository.StorageClass(ctx).FindName(ctx, clusterId, storageName)
	if err != nil {
		_ = level.Error(logger).Log("repository.StorageClass", "Find", "err", err.Error())
		return encode.ErrStorageClassNotfound.Error()
	}
	fmt.Println(find.Name)

	list, err := s.k8sClient.Do(ctx).CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{
		//FieldSelector: fmt.Sprintf("spec.storageClassName=%s", find.Name),
		//FieldSelector: fmt.Sprintf("spec.claimRef.name=%s", "newlender-gfs"),
	})

	fmt.Println(fmt.Sprintf("spec.storageClassName=%s", find.Name))
	if err != nil {
		_ = level.Error(logger).Log("k8sClient.Do", "StorageV1", "StorageClasses", "List", "err", err.Error())
		return encode.ErrStorageClassSyncPv.Wrap(err)
	}

	for _, v := range list.Items {
		b, _ := json.Marshal(v)
		fmt.Println(string(b))
	}

	return
}

func (s *service) SyncPvc(ctx context.Context, clusterId int64, ns string, storageName string) (err error) {
	//list, err := s.k8sClient.Do(ctx).CoreV1().PersistentVolumeClaims(ns).List(ctx, metav1.ListOptions{})
	//if err != nil {
	//	return err
	//}
	return
}

func (s *service) Sync(ctx context.Context, clusterId int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	list, err := s.k8sClient.Do(ctx).StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		_ = level.Error(logger).Log("k8sClient.Do", "StorageV1", "StorageClasses", "List", "err", err.Error())
		return encode.ErrStorageClassSync.Wrap(err)
	}

	for _, item := range list.Items {
		//b, _ := yaml.Marshal(item)
		b, _ := json.Marshal(item)
		storage := &types.StorageClass{
			ClusterId:         clusterId,
			Name:              item.Name,
			Provisioner:       item.Provisioner,
			ReclaimPolicy:     string(*item.ReclaimPolicy),
			VolumeBindingMode: string(*item.VolumeBindingMode),
			ResourceVersion:   item.ResourceVersion,
			Detail:            string(b),
		}
		err := s.repository.StorageClass(ctx).FirstInsert(ctx, storage)
		if err != nil {
			_ = level.Error(logger).Log("repository.StorageClass", "FirstInsert", "err", err.Error())
			continue
		}
		if err := s.repository.StorageClass(ctx).Save(ctx, storage); err != nil {
			_ = level.Error(logger).Log("repository.StorageClass", "Save", "err", err.Error())
		}
	}

	return nil
}

func New(logger log.Logger, traceId string, repository repository.Repository, k8sClient kubernetes.K8sClient) Service {
	logger = log.With(logger, "storageclass", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
		k8sClient:  k8sClient,
	}
}
