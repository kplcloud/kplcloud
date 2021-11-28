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
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	coreV1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/json"

	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

// Service StorageClass模块
type Service interface {
	Sync(ctx context.Context, clusterId int64) (err error)
	SyncPv(ctx context.Context, clusterId int64, storageName string) (err error)
	SyncPvc(ctx context.Context, clusterId int64, storageName string) (err error)
	// Create 创建StorageClass
	Create(ctx context.Context, clusterId int64, ns, name, provisioner string, reclaimPolicy *coreV1.PersistentVolumeReclaimPolicy, volumeBindingMode *storagev1.VolumeBindingMode, remark string) (err error)
	// CreateProvisioner 创建供应者
	CreateProvisioner(ctx context.Context, clusterId int64) (err error)
	// List 存储类列表
	List(ctx context.Context, clusterId int64, page, pageSize int) (res []listResult, total int, err error)
	// Delete 删除存储类
	// 存储类删除需要先判断pvc是否删除，否则无法删除
	Delete(ctx context.Context, clusterId int64, storageName string) (err error)
	// Update 更新存储类
	Update(ctx context.Context, clusterId int64, storageName, provisioner string, reclaimPolicy *coreV1.PersistentVolumeReclaimPolicy, volumeBindingMode *storagev1.VolumeBindingMode, remark string) (err error)
	// Recover 恢复删除的存储类
	Recover(ctx context.Context, clusterId int64, storageName string) (err error)
	// Info 获取存储类详情
	// 从数据库查基本信息 ps: 该存储类下的所有pvc 应该在pvc包，通过clusterId 和 storageName查询
	Info(ctx context.Context, clusterId int64, storageName string) (res infoResult, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) Info(ctx context.Context, clusterId int64, storageName string) (res infoResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	class, err := s.repository.StorageClass(ctx).FindName(ctx, clusterId, storageName)
	if err != nil {
		_ = level.Warn(logger).Log("repository.StorageClass", "FindName", "err", err.Error())
		err = encode.ErrStorageClassNotfound.Wrap(err)
		return
	}

	res.Remark = class.Remark
	res.Name = class.Name
	res.UpdatedAt = class.UpdatedAt
	res.CreatedAt = class.CreatedAt
	res.Provisioner = class.Provisioner
	res.VolumeMode = class.VolumeBindingMode
	res.ResourceVersion = class.ResourceVersion
	res.ClusterName = class.Cluster.Name
	res.ClusterAlias = class.Cluster.Alias
	res.ReclaimPolicy = class.ReclaimPolicy
	// 是否需要自动远程同步？
	return
}

func (s *service) Recover(ctx context.Context, clusterId int64, storageName string) (err error) {
	// 1. 查询删除的存储类
	// 2. 将数据添加到k8s
	// 3. 取消删除字段
	panic("implement me")
}

func (s *service) Delete(ctx context.Context, clusterId int64, storageName string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	class, err := s.repository.StorageClass(ctx).FindName(ctx, clusterId, storageName)
	if err != nil {
		_ = level.Warn(logger).Log("repository.StorageClass", "FindName", "err", err.Error())
		err = encode.ErrStorageClassNotfound.Wrap(err)
		return
	}

	if err = s.repository.StorageClass(ctx).Delete(ctx, class.Id, func() error {
		e := s.k8sClient.Do(ctx).StorageV1().StorageClasses().Delete(ctx, class.Name, metav1.DeleteOptions{})
		if e != nil {
			return errors.Wrap(e, "k8sClient.Do(ctx).StorageV1().StorageClasses().Delete")
		}
		return nil
	}); err != nil {
		return encode.ErrStorageClassDelete.Wrap(err)
	}

	return
}

func (s *service) Update(ctx context.Context, clusterId int64, storageName, provisioner string, reclaimPolicy *coreV1.PersistentVolumeReclaimPolicy, volumeBindingMode *storagev1.VolumeBindingMode, remark string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	class, err := s.repository.StorageClass(ctx).FindName(ctx, clusterId, storageName)
	if err != nil {
		_ = level.Warn(logger).Log("repository.StorageClass", "FindName", "err", err.Error())
		err = encode.ErrStorageClassNotfound.Wrap(err)
		return
	}

	class.Provisioner = provisioner
	class.ReclaimPolicy = string(*reclaimPolicy)
	class.VolumeBindingMode = string(*volumeBindingMode)
	class.Remark = remark

	err = s.repository.StorageClass(ctx).Save(ctx, &class, func() error {
		storageClass, e := s.k8sClient.Do(ctx).StorageV1().StorageClasses().Get(ctx, storageName, metav1.GetOptions{})
		if e != nil {
			return encode.ErrStorageClassNotfound.Wrap(errors.Wrap(e, "k8sClient.Do(ctx).StorageV1().StorageClasses().Get"))
		}
		storageClass.Provisioner = provisioner
		storageClass.ReclaimPolicy = reclaimPolicy
		storageClass.VolumeBindingMode = volumeBindingMode
		storageClass, e = s.k8sClient.Do(ctx).StorageV1().StorageClasses().Update(ctx, storageClass, metav1.UpdateOptions{})
		if e != nil {
			return encode.ErrStorageClassUpdate.Wrap(errors.Wrap(e, "k8sClient.Do(ctx).StorageV1().StorageClasses().Update"))
		}
		return nil
	})
	if err != nil {
		return encode.ErrStorageClassUpdate.Wrap(err)
	}

	return
}

func (s *service) List(ctx context.Context, clusterId int64, page, pageSize int) (res []listResult, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	list, total, err := s.repository.StorageClass(ctx).List(ctx, clusterId, page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.StorageClass", "List", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, listResult{
			Name:          v.Name,
			Provisioner:   v.Provisioner,
			VolumeMode:    v.VolumeBindingMode,
			ReclaimPolicy: v.ReclaimPolicy,
			Remark:        v.Remark,
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
		})
	}

	return
}

func (s *service) CreateProvisioner(ctx context.Context, clusterId int64) (err error) {
	panic("implement me")
}

func (s *service) Create(ctx context.Context, clusterId int64, ns, name, provisioner string, reclaimPolicy *coreV1.PersistentVolumeReclaimPolicy, volumeBindingMode *storagev1.VolumeBindingMode, remark string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	_, err = s.repository.StorageClass(ctx).FindName(ctx, clusterId, name)
	if err == nil {
		return encode.ErrStorageClassExists.Error()
	}
	if !gorm.IsRecordNotFoundError(err) {
		return encode.ErrStorageClassExists.Wrap(err)
	}
	storage := &types.StorageClass{}
	storage.ClusterId = clusterId
	storage.ReclaimPolicy = string(*reclaimPolicy)
	storage.VolumeBindingMode = string(*volumeBindingMode)
	storage.Provisioner = provisioner
	storage.Name = name
	storage.Remark = remark
	if err = s.repository.StorageClass(ctx).Save(ctx, storage, func() error {
		// TODO: 考虑使用模版
		create, err := s.k8sClient.Do(ctx).StorageV1().StorageClasses().Create(ctx, &storagev1.StorageClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Provisioner:       provisioner,
			ReclaimPolicy:     reclaimPolicy,
			VolumeBindingMode: volumeBindingMode,
		}, metav1.CreateOptions{})
		if err != nil {
			return encode.ErrStorageClassCreate.Wrap(err)
		}
		b, _ := json.Marshal(create)
		storage.ResourceVersion = create.ResourceVersion
		storage.Detail = string(b)
		return nil
	}); err != nil {
		_ = level.Error(logger).Log("repository.StorageClass", "Save", "err", err.Error())
		return encode.ErrStorageClassCreate.Wrap(err)
	}

	return nil
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

func (s *service) SyncPvc(ctx context.Context, clusterId int64, storageName string) (err error) {
	//logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	list, err := s.k8sClient.Do(ctx).CoreV1().PersistentVolumeClaims(coreV1.NamespaceAll).List(ctx, metav1.ListOptions{
		//FieldSelector: fields.SelectorFromSet(fields.Set{
		//	"spec.storageClassName": storageName,
		//}).String(),
		LabelSelector: fields.SelectorFromSet(fields.Set{
			"spec.storageClassName": storageName,
		}).String(),
	})
	fmt.Println(fields.SelectorFromSet(fields.Set{
		"spec.storageClassName": storageName,
	}).String())
	if err != nil {
		return err
	}
	for _, v := range list.Items {
		fmt.Println(v.Name)
	}
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
		if err := s.repository.StorageClass(ctx).Save(ctx, storage, nil); err != nil {
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
