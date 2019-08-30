/**
 * @Time : 2019-06-26 15:14
 * @Author : solacowa@gmail.com
 * @File : persistentvolumeclaim
 * @Software: GoLand
 */

package repository

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"gopkg.in/guregu/null.v3"
)

type PvcRepository interface {
	Find(ns, name string) (rs *types.PersistentVolumeClaim, err error)
	FindBy(ns string, offset, limit int) (res []*types.PersistentVolumeClaim, err error)
	Count(ns string) (int64, error)
	Delete(ns, name string) (err error)
	FindAll() (res []*types.PersistentVolumeClaim, err error)
	FirstOrCreate(ns, name, accessModes, storage, storageClassName, detail, selector string, labels map[string]string) error
}

type pvc struct {
	db *gorm.DB
}

func NewPvcRepository(db *gorm.DB) PvcRepository {
	return &pvc{db: db}
}

func (c *pvc) FindAll() (res []*types.PersistentVolumeClaim, err error) {
	err = c.db.Find(&res).Error
	return
}

func (c *pvc) FirstOrCreate(ns, name, accessModes, storage, storageClassName, detail, selector string, labels map[string]string) error {
	b, _ := json.Marshal(labels)
	class := types.PersistentVolumeClaim{
		Name:             name,
		Namespace:        ns,
		Detail:           null.StringFrom(detail),
		Labels:           null.StringFrom(string(b)),
		Selector:         null.StringFrom(selector),
		Storage:          storage,
		StorageClassName: null.StringFrom(storageClassName),
		AccessModes:      accessModes,
	}
	return c.db.FirstOrCreate(&class, types.PersistentVolumeClaim{
		Namespace: ns,
		Name:      name,
	}).Error
}

func (c *pvc) Find(ns, name string) (rs *types.PersistentVolumeClaim, err error) {
	var res types.PersistentVolumeClaim
	err = c.db.First(&res, "name = ? AND namespace = ?", name, ns).Error
	return &res, nil
}

func (c *pvc) Delete(ns, name string) (err error) {
	pvc := types.PersistentVolumeClaim{
		Name:      name,
		Namespace: ns,
	}
	return c.db.Delete(&pvc, "name = ? AND namespace = ?", name, ns).Error
}

func (c *pvc) FindBy(ns string, offset, limit int) (res []*types.PersistentVolumeClaim, err error) {
	err = c.db.Order(gorm.Expr("id DESC")).
		Offset(offset).Limit(limit).
		Find(&res, "namespace = ?", ns).Error
	return
}

func (c *pvc) Count(ns string) (int64, error) {
	var count int64
	if err := c.db.Model(&types.PersistentVolumeClaim{}).Where("namespace = ? ", ns).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
