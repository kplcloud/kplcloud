/**
 * @Time : 2019/7/5 6:33 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : configdata
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type ConfigDataRepository interface {
	Create(confMap *types.ConfigMapData) error
	Update(id int64, valu, path string) error
	Find(ns, name string) (list []*types.ConfigMapData, err error)
	Delete(configMapId int64) error
	FindByConfMapId(configMapId int64) (re *types.ConfigMapData, err error)
	FindByConfMapIdAndKey(configMapId int64, key string) (re *types.ConfigMapData, notFound bool)
	FindById(id int64) (re *types.ConfigMapData, err error)
	Count(configMapId int64) (count int, err error)
	FindOffsetLimit(configMapId int64, offset, limit int) (res []*types.ConfigMapData, err error)
	DeleteById(id int64) error
}

type configData struct {
	db *gorm.DB
}

func NewConfigDataRepository(db *gorm.DB) ConfigDataRepository {
	return &configData{db: db}
}

func (c *configData) FindByConfMapId(configMapId int64) (re *types.ConfigMapData, err error) {
	var res types.ConfigMapData
	err = c.db.Find(&res, "config_map_id = ?", configMapId).Error
	return &res, err
}

func (c *configData) FindByConfMapIdAndKey(configMapId int64, key string) (re *types.ConfigMapData, notFound bool) {
	var res types.ConfigMapData
	notFound = c.db.
		Where("`key` = ?", key).
		Where("config_map_id = ?", configMapId).
		First(&res).RecordNotFound()
	return &res, notFound
}

func (c *configData) FindById(id int64) (re *types.ConfigMapData, err error) {
	var res types.ConfigMapData
	err = c.db.Find(&res, "id = ?", id).Error
	return &res, err
}

func (c *configData) Create(confMap *types.ConfigMapData) error {
	return c.db.Save(confMap).Error
}

func (c *configData) Update(id int64, value, path string) error {
	var temp types.ConfigMapData
	return c.db.Model(&temp).Where("id = ?", id).Update(&types.ConfigMapData{
		Value: value,
	}).Error
}

func (c *configData) Find(ns, name string) (list []*types.ConfigMapData, err error) {
	var confMap types.ConfigMap
	var confData types.ConfigMapData
	query := c.db.Model(&confData).Joins("inner join " + confMap.TableName() + " t1 on t1.id = " + confData.TableName() + ".config_map_id")
	if ns != "" {
		query = query.Where("t1.namespace = ?", ns)
	}
	if name != "" {
		query = query.Where("t1.name = ?", name)
	}
	err = query.Preload("ConfigMap").Find(&list).Error
	return
}

func (c *configData) Delete(configMapId int64) error {
	return c.db.Where("config_map_id = ?", configMapId).Delete(types.ConfigMapData{}).Error
}

func (c *configData) DeleteById(id int64) error {
	return c.db.Where("id = ?", id).Delete(types.ConfigMapData{}).Error
}

/**
 * @Title 获取总数
 */
func (c *configData) Count(configMapId int64) (count int, err error) {
	var n types.ConfigMapData
	query := c.db.Model(&n).Where("config_map_id = ?", configMapId)
	err = query.Count(&count).Error
	return
}

/**
* @Title 获取列表
 */
func (c *configData) FindOffsetLimit(configMapId int64, offset, limit int) (res []*types.ConfigMapData, err error) {
	var list []*types.ConfigMapData
	query := c.db.Where("config_map_id = ?", configMapId)
	err = query.Offset(offset).Limit(limit).Find(&list).Error
	return list, err
}
