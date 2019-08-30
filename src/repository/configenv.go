/**
 * Created by GoLand.
 * Email: xzghua@gmail.com
 * Date: 2019-07-17
 * Time: 15:07
 */
package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type ConfigEnvRepository interface {
	GetConfigEnvByNameNs(name, namespace string) ([]*types.ConfigEnv, error)
	GetConfigEnvCountByNameNs(name, ns string) (cnt int64, err error)
	GetConfigEnvPaginate(name, ns string, offset int, limit int) ([]types.ConfigEnv, error)
	CreateConfEnv(name, ns, envKey, envVar, EnvDesc string) error
	FindById(id int64) (types.ConfigEnv, bool)
	Update(id int64, confEnv types.ConfigEnv) error
	Delete(id int64) error
}

type configEnv struct {
	db *gorm.DB
}

func NewConfigEnvRepository(db *gorm.DB) ConfigEnvRepository {
	return &configEnv{db: db}
}

func (c *configEnv) GetConfigEnvByNameNs(name, namespace string) (configEnvs []*types.ConfigEnv, err error) {
	err = c.db.Where("name = ?", name).Where("namespace = ?", namespace).Find(&configEnvs).Error
	return configEnvs, err
}

func (c *configEnv) GetConfigEnvCountByNameNs(name, ns string) (cnt int64, err error) {
	confEnv := new(types.ConfigEnv)
	query := c.db.Model(&confEnv)
	if name != "" {
		query = query.Where("name = ?", name)
	}

	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}
	err = query.Count(&cnt).Error
	return cnt, err
}

func (c *configEnv) GetConfigEnvPaginate(name, ns string, offset int, limit int) ([]types.ConfigEnv, error) {
	var confEnv []types.ConfigEnv
	query := c.db
	if name != "" {
		query = query.Where("name = ?", name)
	}

	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}
	err := query.Offset(offset).Limit(limit).Find(&confEnv).Error
	return confEnv, err
}

func (c *configEnv) CreateConfEnv(name, ns, envKey, envVar, EnvDesc string) error {
	confEnv := types.ConfigEnv{
		Name:      name,
		Namespace: ns,
		EnvKey:    envKey,
		EnvVar:    envVar,
		EnvDesc:   EnvDesc,
	}
	err := c.db.Create(&confEnv).Error
	return err
}

func (c *configEnv) FindById(id int64) (confEnv types.ConfigEnv, res bool) {
	res = c.db.Find(&confEnv, "id = ?", id).RecordNotFound()
	return
}

func (c *configEnv) Update(id int64, confEnv types.ConfigEnv) error {
	return c.db.Where("id = ?", id).Save(&confEnv).Error
}

func (c *configEnv) Delete(id int64) error {
	var confEnv types.ConfigEnv
	return c.db.Delete(&confEnv, "id = ?", id).Error
}
