/**
 * @Time : 2019/6/27 10:19 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : webhook
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

const (
	AppTarget    = "app"
	GlobalTarget = "global"
)

type WebhookRepository interface {
	Create(w *types.Webhook) error
	UpdateById(w *types.Webhook) error
	FindById(id int) (v *types.Webhook, err error)
	FindByName(name string) (v *types.Webhook, err error)
	Count(name, appName, ns string) (count int, err error)
	FindOffsetLimit(name, appName, ns string, offset, limit int) (res []*types.Webhook, err error)
	Delete(id int) error
	DeleteEvents(w *types.Webhook) error
	CreateEvents(w *types.Webhook, events ...*types.Event) error
}

type webhook struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) WebhookRepository {
	return &webhook{db: db}
}

func (c *webhook) Create(w *types.Webhook) error {
	return c.db.Create(&w).Error
}

func (c *webhook) UpdateById(w *types.Webhook) error {
	var hook types.Webhook
	return c.db.Model(&hook).Where("id = ?", w.ID).Save(w).Error
}

func (c *webhook) FindById(id int) (v *types.Webhook, err error) {
	hook := types.Webhook{}
	c.db.First(&hook, id)
	err = c.db.Model(&hook).Related(&hook.Events, "Events").Error
	return &hook, err
}

func (c *webhook) FindByName(name string) (v *types.Webhook, err error) {
	hook := types.Webhook{}
	err = c.db.Where("name = ?", name).First(&hook).Error
	return &hook, err
}

func (c *webhook) Count(name, appName, ns string) (count int, err error) {
	var temp types.Webhook
	query := c.db.Model(&temp)
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if appName != "" {
		query = query.Where("app_name = ?", appName)
	}
	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}
	err = query.Count(&count).Error
	return
}

func (c *webhook) FindOffsetLimit(name, appName, ns string, offset, limit int) (res []*types.Webhook, err error) {
	var list []*types.Webhook
	query := c.db
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if appName != "" {
		query = query.Where("app_name = ?", appName)
	}
	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}
	// 关联member相关信息
	err = query.Preload("Events").Preload("Member", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,username,email")
	}).Offset(offset).Limit(limit).Find(&list).Error
	return list, err
}

func (c *webhook) Delete(id int) error {
	webHook, err := c.FindById(id)
	if err != nil {
		return nil
	}
	if err := c.DeleteEvents(webHook); err != nil {
		return err
	}
	return c.db.Where("id=?", id).Delete(types.Webhook{}).Error
}

func (c *webhook) DeleteEvents(w *types.Webhook) error {
	return c.db.Model(&w).Association("Events").Clear().Error
}

func (c *webhook) CreateEvents(w *types.Webhook, events ...*types.Event) error {
	err := c.db.Model(w).Association("Events").Append(events).Error
	return err
}
