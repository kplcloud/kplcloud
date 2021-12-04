/**
 * @Time : 2019-06-26 14:34
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package persistentvolumeclaim

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type Middleware func(Service) Service

type Service interface {
	// Sync 同步pvc
	Sync(ctx context.Context, clusterId int64, ns string) (err error)
	// Get 获取pvc详情
	Get(ctx context.Context, clusterId int64, ns, name string) (res result, err error)
	// Delete 删除存储卷声明
	// 查看的有绑定关系，如果存在绑定要求先解除绑定关系，如果没有绑定关系按下面流程删除
	// 删除pv -> 删除k8s pvc -> 删除pvc
	Delete(ctx context.Context, clusterId int64, ns, name string) (err error)
	// Create 创建持久化存储卷
	Create(ctx context.Context, clusterId int64, ns, name, storage, storageClassName string, accessModes []string) (err error)
	// List 持久化存储卷列表
	// 选择集群之后获取存储类和空间列表，暂时不考虑存储类，只考虑空间
	List(ctx context.Context, clusterId int64, storageClass, ns string, page, pageSize int) (resp []result, total int, err error)
	// All 当前空间下所有的pvc
	All(ctx context.Context, clusterId int64) (resp map[string]interface{}, err error)
}

type service struct {
	logger     log.Logger
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
	traceId    string
}

func (s *service) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	pvcs, err := s.k8sClient.Do(ctx).CoreV1().PersistentVolumeClaims(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		_ = level.Error(logger).Log("PersistentVolumeClaims", "List", "err", err.Error())
		return encode.ErrPersistentVolumeClaimList.Wrap(err)
	}

	for _, pvc := range pvcs.Items {
		storage, e := s.repository.StorageClass(ctx).FindName(ctx, clusterId, *pvc.Spec.StorageClassName)
		if e != nil {
			_ = level.Error(logger).Log("repository.StorageClass", "FindName", "err", err.Error())
			continue
		}
		accessModels, _ := json.Marshal(pvc.Spec.AccessModes)
		labels, _ := json.Marshal(pvc.Labels)

		if er := s.repository.Pvc(ctx).Save(ctx, &types.PersistentVolumeClaim{
			Name:           pvc.Name,
			Namespace:      pvc.Namespace,
			AccessModes:    string(accessModels),
			Labels:         string(labels),
			RequestStorage: pvc.Spec.Resources.Requests.Storage().String(),
			LimitStorage:   pvc.Spec.Resources.Limits.Storage().String(),
			StorageClassId: storage.Id,
			ClusterId:      clusterId,
		}, nil); er != nil {
			_ = level.Error(logger).Log("repository.Pvc", "Save", "err", er.Error())
		}
	}

	return
}

func (s *service) Get(ctx context.Context, clusterId int64, ns, name string) (res result, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	pvc, err := s.repository.Pvc(ctx).FindByName(ctx, clusterId, ns, name)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Pvc", "FindByName", "err", err.Error())
		err = encode.ErrPersistentVolumeClaimNotfound.Error()
		return
	}
	kpvc, err := s.k8sClient.Do(ctx).CoreV1().PersistentVolumeClaims(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		_ = level.Warn(logger).Log("k8sClient.Do.CoreV1.PersistentVolumeClaims", "Get", "err", err.Error())
		err = encode.ErrPersistentVolumeClaimNotfound.Error()
		return
	}

	//pv, err := s.k8sClient.Do(ctx).CoreV1().PersistentVolumes().Get(ctx, kpvc.Spec.VolumeName, metav1.GetOptions{})
	//if err != nil {
	//	_ = level.Warn(logger).Log("k8sClient.Do.CoreV1.PersistentVolumes", "Get", "err", err.Error())
	//}
	var accessModes []string
	_ = json.Unmarshal([]byte(pvc.AccessModes), &accessModes)

	res.Namespace = pvc.Namespace
	res.Name = pvc.Name
	res.StorageClass = *kpvc.Spec.StorageClassName
	res.Status = string(kpvc.Status.Phase)
	res.CreatedAt = pvc.CreatedAt
	res.UpdatedAt = pvc.UpdatedAt
	res.RequestStorage = pvc.RequestStorage
	res.LimitStorage = pvc.LimitStorage
	res.Annotations = kpvc.Annotations
	res.Labels = kpvc.Labels
	res.VolumeName = kpvc.Spec.VolumeName
	res.AccessModes = accessModes
	res.ClusterName = pvc.Cluster.Name
	res.ClusterAlias = pvc.Cluster.Alias

	return
}

func (s *service) Delete(ctx context.Context, clusterId int64, ns, name string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	spvc, err := s.repository.Pvc(ctx).FindByName(ctx, clusterId, ns, name)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Pvc", "FindByName", "err", err)
		if gorm.IsRecordNotFoundError(err) {
			return encode.ErrPersistentVolumeClaimNotfound.Error()
		}
		return err
	}

	// TODO: 查绑定关系，如果有关系返回失败，要求先解除绑定关系
	// TODO: 删除pvc

	if err := s.repository.Pvc(ctx).Delete(ctx, spvc.Id, func() error {
		return s.k8sClient.Do(ctx).CoreV1().PersistentVolumeClaims(ns).Delete(ctx, name, metav1.DeleteOptions{})
	}); err != nil {
		_ = level.Error(logger).Log("pvc", "Delete", "err", err.Error())
		return encode.ErrPersistentVolumeClaimDelete.Wrap(err)
	}

	return

}

func (s *service) Create(ctx context.Context, clusterId int64, ns, name, storage, storageClassName string, accessModes []string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	sc, err := s.repository.StorageClass(ctx).FindName(ctx, clusterId, storageClassName)
	if err != nil {
		_ = level.Error(logger).Log("repository.StorageClass", "FindName", "err", err)
		return encode.ErrStorageClassNotfound.Error()
	}
	spvc, err := s.repository.Pvc(ctx).FindByName(ctx, clusterId, ns, name)
	if !gorm.IsRecordNotFoundError(err) {
		return encode.ErrPersistentVolumeClaimExists.Error()
	}
	var pvc *corev1.PersistentVolumeClaim
	_, err = s.repository.K8sTpl(ctx).EncodeTemplate(ctx, types.KindPersistentVolumeClaim, map[string]interface{}{
		"name":             name,
		"namespace":        ns,
		"accessModes":      accessModes,
		"storage":          storage,
		"storageClassName": sc.Name,
	}, &pvc)
	if err != nil {
		return encode.ErrPersistentVolumeClaimCreate.Wrap(err)
	}

	am, _ := json.Marshal(accessModes)
	spvc.Name = name
	spvc.Namespace = ns
	spvc.ClusterId = clusterId
	spvc.StorageClassId = sc.Id
	spvc.AccessModes = string(am)
	spvc.Remark = ""
	spvc.RequestStorage = storage

	if err = s.repository.Pvc(ctx).Save(ctx, &spvc, func() error {
		pvc, err = s.k8sClient.Do(ctx).CoreV1().PersistentVolumeClaims(ns).Create(ctx, pvc, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrap(err, "CoreV1.PersistentVolumeClaims.Create")
		}
		spvc.Status = string(pvc.Status.Phase)
		return nil
	}); err != nil {
		_ = level.Error(logger).Log("repository.Pvc", "Save", "err", err.Error())
		return encode.ErrPersistentVolumeClaimCreate.Wrap(err)
	}

	return
}

func (s *service) List(ctx context.Context, clusterId int64, storageClass, ns string, page, pageSize int) (resp []result, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	var storageClassIds []int64
	if !strings.EqualFold(storageClass, "") {
		if storage, e := s.repository.StorageClass(ctx).FindName(ctx, clusterId, storageClass); e == nil {
			storageClassIds = []int64{storage.Id}
		} else {
			_ = level.Warn(logger).Log("repository.StorageClass", "FindName", "err", err.Error())
		}
	}

	list, total, err := s.repository.Pvc(ctx).List(ctx, clusterId, storageClassIds, ns, "", page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.Pvc", "List", "err", err.Error())
		return
	}

	for _, v := range list {
		var accessModes []string
		_ = json.Unmarshal([]byte(v.AccessModes), &accessModes)
		resp = append(resp, result{
			Name:           v.Name,
			Namespace:      v.Namespace,
			StorageClass:   v.StorageClass.Name,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			AccessModes:    accessModes,
			Remark:         v.Remark,
			RequestStorage: v.RequestStorage,
			LimitStorage:   v.LimitStorage,
			Status:         v.Status,
		})
	}

	return
}

func (s *service) All(ctx context.Context, clusterId int64) (resp map[string]interface{}, err error) {
	panic("implement me")
}

func New(logger log.Logger, traceId string, client kubernetes.K8sClient, repository repository.Repository) Service {
	return &service{
		logger: logger, k8sClient: client,
		repository: repository,
		traceId:    traceId,
	}
}
