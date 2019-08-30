/**
 * @Time : 2019/7/5 11:50 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : configmap
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type ConfigMapRepository interface {
	Create(confMap *types.ConfigMap) (rs *types.ConfigMap, err error)
	Find(ns, name string) (*types.ConfigMap, bool)
	FindById(id int64) (*types.ConfigMap, bool)
	Count(ns, name string) (count int, err error)
	FindOffsetLimit(ns, name string, offset, limit int) (res []*types.ConfigMap, err error)
	Update(ns, name, desc string) error
	Delete(id int64) error
	DeleteByNsName(ns, name string) error
}

type configMap struct {
	db *gorm.DB
}

func NewConfigMapRepository(db *gorm.DB) ConfigMapRepository {
	return &configMap{db: db}
}

func (c *configMap) Create(confMap *types.ConfigMap) (rs *types.ConfigMap, err error) {
	err = c.db.Save(confMap).Error
	return confMap, err
}

func (c *configMap) Find(ns, name string) (*types.ConfigMap, bool) {
	var rs types.ConfigMap
	res := c.db.First(&rs, "namespace = ? AND name = ?", ns, name).
		Preload("ConfigData").RecordNotFound()
	return &rs, res
}

func (c *configMap) FindById(id int64) (*types.ConfigMap, bool) {
	var rs types.ConfigMap
	res := c.db.First(&rs, id).RecordNotFound()
	return &rs, res
}

func (c *configMap) Count(ns, name string) (count int, err error) {
	query := c.db.Model(&types.ConfigMap{})
	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err = query.Count(&count).Error
	return
}

func (c *configMap) FindOffsetLimit(ns, name string, offset, limit int) (res []*types.ConfigMap, err error) {
	var list []*types.ConfigMap
	query := c.db.Model(&types.ConfigMap{})
	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err = query.Offset(offset).Limit(limit).Find(&list).Error
	return list, err
}

func (c *configMap) Update(ns, name, desc string) error {
	var temp types.ConfigMap
	return c.db.Model(&temp).Where("namespace = ? AND name = ?", ns, name).Update(&types.ConfigMap{
		Namespace: ns, Name: name, Desc: desc,
	}).Error
}

func (c *configMap) Delete(id int64) error {
	return c.db.Where("id = ?", id).Delete(types.ConfigMap{}).Error
}

func (c *configMap) DeleteByNsName(ns, name string) error {
	if configMap, notExist := c.Find(ns, name); notExist == false {
		c.db.Delete(types.ConfigData{ConfigMap: *configMap})
		return c.db.Delete(types.ConfigMap{ID: configMap.ID}).Error
	}
	return nil
}
