/**
 * @Time : 2019/7/17 2:19 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : consul
 * @Software: GoLand
 */

package repository

import (
	"github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type ConsulRepository interface {
	Find(ns, name string) (*types.Consul, bool)
	FirstOrCreate(entry *api.ACLEntry, ns, name string) (consul *types.Consul, err error)
	UpdateOrCreate(entry *api.ACLEntry, ns, name string) (err error)
	Count(ns, name string) (count int, err error)
	FindOffsetLimit(ns, name string, offset, limit int) (res []*types.Consul, err error)
	Delete(ns, name string) error
}

type consul struct {
	db *gorm.DB
}

func NewConsulReporitory(db *gorm.DB) ConsulRepository {
	return &consul{db: db}
}

func (c *consul) Find(ns, name string) (*types.Consul, bool) {
	var rs types.Consul
	res := c.db.First(&rs, "namespace = ? AND name = ?", ns, name).RecordNotFound()
	return &rs, res
}

func (c *consul) FirstOrCreate(entry *api.ACLEntry, ns, name string) (consul *types.Consul, err error) {
	class := types.Consul{
		Namespace:   ns,
		Name:        name,
		Type:        entry.Type,
		Rules:       entry.Rules,
		Token:       entry.ID,
		CreateIndex: int64(entry.CreateIndex),
		ModifyIndex: int64(entry.ModifyIndex),
	}
	err = c.db.FirstOrCreate(&class, types.Consul{
		Namespace: ns,
		Name:      name,
	}).Error
	return &class, err
}

func (c *consul) UpdateOrCreate(entry *api.ACLEntry, ns, name string) (err error) {
	var temp types.Consul
	consul, notExistState := c.Find(ns, name)
	consul.Rules = entry.Rules
	consul.Type = entry.Type
	consul.Token = entry.ID
	consul.CreateIndex = int64(entry.CreateIndex)
	consul.ModifyIndex = int64(entry.ModifyIndex)

	if notExistState == true {
		return c.db.Model(&temp).Where("namespace = ? AND name = ?", ns, name).Update(consul).Error
	} else {
		consul.Name = name
		consul.Namespace = ns
		return c.db.Save(consul).Error
	}
}

func (c *consul) Count(ns, name string) (count int, err error) {
	query := c.db.Model(&types.Consul{})
	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err = query.Count(&count).Error
	return
}

func (c *consul) FindOffsetLimit(ns, name string, offset, limit int) (res []*types.Consul, err error) {
	var list []*types.Consul
	query := c.db.Model(&types.Consul{})
	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err = query.Offset(offset).Limit(limit).Find(&list).Error
	return list, err
}

func (c *consul) Delete(ns, name string) error {
	return c.db.Where("namespace = ? AND name = ?", ns, name).Delete(types.Consul{}).Error
}
