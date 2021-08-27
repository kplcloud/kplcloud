/**
 * @Time : 2019-06-25 19:26
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package storage

import (
	"context"
	"errors"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

var (
	ErrInvalidArgument               = errors.New("invalid argument")
	ErrStorageClassNotFound          = errors.New("没有找到相关存储类")
	ErrStorageClassIsExists          = errors.New("存储类已经存在")
	ErrStorageClassTemplateNotExists = errors.New("存储类模版不存在")
	ErrStorageClassTemplateEncode    = errors.New("存储类模版处理失败")
	ErrStorageClassK8sCreate         = errors.New("存储类创建失败")
	ErrStorageClassList              = errors.New("存储类列表获取失败")
	ErrStorageClassK8sGet            = errors.New("存储类获取失败")
)

type Service interface {
	// 同步storageclass
	Sync(_ context.Context) (err error)

	// 根据ns获取 storageclass
	Get(ctx context.Context, name string) (rs interface{}, err error)

	// 创建storageclass
	Post(ctx context.Context, name, provisioner, reclaimRolicy, volumeBindingMode string) (err error)

	// 删除存储类
	Delete(ctx context.Context, name string) (err error)

	// 存储类列表
	List(ctx context.Context, offset, limit int) (res []*types.StorageClass, err error)
}

type storageClass struct {
	*v1.StorageClass
	PersistentVolumeList []coreV1.PersistentVolume `json:"persistent_volume_list"`
}

type persistentVolumeListChannel struct {
	List  chan *coreV1.PersistentVolumeList
	Error chan error
}

type service struct {
	logger     log.Logger
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func NewService(logger log.Logger, client kubernetes.K8sClient, repository repository.Repository) Service {
	return &service{logger, client, repository}
}

func (c *service) Sync(_ context.Context) (err error) {
	storages, err := c.k8sClient.Do().StorageV1().
		StorageClasses().
		List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("storageclass", "list", "err", err.Error())
		return
	}
	for _, class := range storages.Items {
		class.APIVersion = "repository.StorageClass().k8s.io/v1"
		class.Kind = "StorageClass"
		b, _ := yaml.Marshal(class)
		if err = c.repository.StorageClass().FirstOrCreate(class.Name, class.Provisioner,
			repository.PersistentVolumeReclaimPolicy(*class.ReclaimPolicy),
			repository.VolumeBindingMode(*class.VolumeBindingMode), b); err != nil {
			_ = level.Error(c.logger).Log("storage", "FirstOrCreate", "err", err.Error())
		}
	}
	return
}

func (c *service) Get(ctx context.Context, name string) (rs interface{}, err error) {
	if _, err = c.repository.StorageClass().Find(name); err != nil {
		_ = level.Error(c.logger).Log("storage", "find", "err", err.Error())
		return nil, ErrStorageClassNotFound
	}

	storage, err := c.k8sClient.Do().StorageV1().StorageClasses().Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("StorageClasses", "Get", "err", err.Error())
		return nil, ErrStorageClassK8sGet
	}

	channels := c.getPersistentVolumeListChannel(1)
	persistentVolumeList := <-channels.List
	err = <-channels.Error

	if err != nil {
		return
	}

	storagePersistentVolumes := make([]coreV1.PersistentVolume, 0)
	for _, pv := range persistentVolumeList.Items {
		if strings.Compare(pv.Spec.StorageClassName, name) == 0 {
			storagePersistentVolumes = append(storagePersistentVolumes, pv)
		}
	}

	return &storageClass{
		storage,
		storagePersistentVolumes,
	}, nil
}

func (c *service) getPersistentVolumeListChannel(numReads int) persistentVolumeListChannel {
	channel := persistentVolumeListChannel{
		List:  make(chan *coreV1.PersistentVolumeList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := c.k8sClient.Do().CoreV1().PersistentVolumes().List(metav1.ListOptions{
			LabelSelector: labels.Everything().String(),
			FieldSelector: fields.Everything().String(),
		})

		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

func (c *service) Post(ctx context.Context, name, provisioner, reclaimRolicy, volumeBindingMode string) (err error) {
	storage, err := c.repository.StorageClass().Find(name)
	if err == nil && storage != nil && storage.Name != "" {
		return ErrStorageClassIsExists
	}

	tpl, err := c.repository.Template().FindByKindType(repository.StorageClassKind)
	if err != nil {
		_ = level.Error(c.logger).Log("template", "FindByKindType", "err", err.Error())
		return ErrStorageClassTemplateNotExists
	}

	enTpl, err := encode.EncodeTemplate(repository.StorageClassKind.ToString(), tpl.Detail, map[string]string{
		"name":              name,
		"provisioner":       provisioner,
		"reclaimPolicy":     reclaimRolicy,
		"volumeBindingMode": volumeBindingMode,
	})
	if err != nil {
		_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
		return ErrStorageClassTemplateEncode
	}

	var storageClass *v1.StorageClass
	if err = yaml.Unmarshal([]byte(enTpl), &storageClass); err != nil {
		_ = level.Error(c.logger).Log("yaml", "Unmarshal", "err", err.Error())
		return
	}

	if storageClass, err = c.k8sClient.Do().StorageV1().StorageClasses().Create(storageClass); err != nil {
		_ = level.Error(c.logger).Log("StorageClasses", "Create", "err", err.Error())
		return errors.New(ErrStorageClassK8sCreate.Error() + err.Error())
	}

	defer func() {
		if err != nil {
			if e := c.k8sClient.Do().StorageV1().StorageClasses().Delete(name, &metav1.DeleteOptions{}); e != nil {
				_ = level.Warn(c.logger).Log("StorageClasses", "Delete", "err", e.Error())
			}
		}
	}()

	b, _ := yaml.Marshal(storageClass)

	if err = c.repository.StorageClass().Create(&types.StorageClass{
		Name:              name,
		Provisioner:       provisioner,
		ReclaimPolicy:     repository.PersistentVolumeReclaimPolicy(reclaimRolicy).String(),
		VolumeBindingMode: repository.VolumeBindingMode(volumeBindingMode).String(),
		Detail:            string(b),
	}); err != nil {
		_ = level.Error(c.logger).Log("storage", "create", "err", err.Error())
		return ErrStorageClassK8sCreate
	}

	return
}

func (c *service) Delete(ctx context.Context, name string) (err error) {
	defer func() {
		if err == nil {
			if e := c.repository.StorageClass().Delete(name); e != nil {
				_ = level.Warn(c.logger).Log("storage", "delete", "err", e.Error())
			}
		}
	}()

	if err := c.k8sClient.Do().StorageV1().StorageClasses().Delete(name, &metav1.DeleteOptions{}); err != nil {
		_ = level.Error(c.logger).Log("StorageClasses", "Delete", "err", err.Error())
		return err
	}

	return
}

func (c *service) List(ctx context.Context, offset, limit int) (res []*types.StorageClass, err error) {
	res, err = c.repository.StorageClass().FindOffsetLimit(offset, limit)
	if err != nil {
		_ = level.Error(c.logger).Log("storage", "FindOffsetLimit", "err", err.Error())
		return nil, ErrStorageClassList
	}
	return
}
