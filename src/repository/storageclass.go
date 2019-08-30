/**
 * @Time : 2019-06-26 10:13
 * @Author : solacowa@gmail.com
 * @File : storageclass
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type PersistentVolumeReclaimPolicy string

const (
	// PersistentVolumeReclaimRecycle means the volume will be recycled back into the pool of unbound persistent volumes on release from its claim.
	// The volume plugin must support Recycling.
	PersistentVolumeReclaimRecycle PersistentVolumeReclaimPolicy = "Recycle"
	// PersistentVolumeReclaimDelete means the volume will be deleted from Kubernetes on release from its claim.
	// The volume plugin must support Deletion.
	PersistentVolumeReclaimDelete PersistentVolumeReclaimPolicy = "Delete"
	// PersistentVolumeReclaimRetain means the volume will be left in its current phase (Released) for manual reclamation by the administrator.
	// The default policy is Retain.
	PersistentVolumeReclaimRetain PersistentVolumeReclaimPolicy = "Retain"
)

func (c PersistentVolumeReclaimPolicy) String() string {
	return string(c)
}

type VolumeBindingMode string

const (
	// VolumeBindingImmediate indicates that PersistentVolumeClaims should be
	// immediately provisioned and bound.  This is the default mode.
	VolumeBindingImmediate VolumeBindingMode = "Immediate"

	// VolumeBindingWaitForFirstConsumer indicates that PersistentVolumeClaims
	// should not be provisioned and bound until the first Pod is created that
	// references the PeristentVolumeClaim.  The volume provisioning and
	// binding will occur during Pod scheduing.
	VolumeBindingWaitForFirstConsumer VolumeBindingMode = "WaitForFirstConsumer"
)

func (c VolumeBindingMode) String() string {
	return string(c)
}

type StorageClassRepository interface {
	Find(name string) (res *types.StorageClass, err error)
	FirstOrCreate(name, provisioner string, reclaimPolicy PersistentVolumeReclaimPolicy, volumeBindingMode VolumeBindingMode, detail []byte) error
	Create(storage *types.StorageClass) error
	Delete(name string) error
	FindOffsetLimit(offset, limit int) (res []*types.StorageClass, err error)
}

type storageClass struct {
	db *gorm.DB
}

func NewStorageClassRepository(db *gorm.DB) StorageClassRepository {
	return &storageClass{db: db}
}

func (c *storageClass) Find(name string) (res *types.StorageClass, err error) {
	var rs types.StorageClass
	err = c.db.First(&rs, "name = ?", name).Error
	return &rs, nil
}

func (c *storageClass) FirstOrCreate(name, provisioner string, reclaimPolicy PersistentVolumeReclaimPolicy, volumeBindingMode VolumeBindingMode, detail []byte) error {
	class := types.StorageClass{
		Name:              name,
		Provisioner:       provisioner,
		ReclaimPolicy:     reclaimPolicy.String(),
		VolumeBindingMode: volumeBindingMode.String(),
		Detail:            string(detail),
	}
	return c.db.FirstOrCreate(&class, types.StorageClass{
		Name: name,
	}).Error
}

func (c *storageClass) Create(storage *types.StorageClass) error {
	return c.db.Save(storage).Error
}

func (c *storageClass) Delete(name string) error {
	class := types.StorageClass{
		Name: name,
	}
	return c.db.Delete(class, "name = ?", name).Error
}

func (c *storageClass) FindOffsetLimit(offset, limit int) (res []*types.StorageClass, err error) {
	//var list []*types.StorageClass
	err = c.db.Order(gorm.Expr("id DESC")).Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}
