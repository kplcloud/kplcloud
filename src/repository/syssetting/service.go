/**
 * @Time: 2020/9/26 19:09
 * @Author: solacowa@gmail.com
 * @File: setting
 * @Software: GoLand
 */

package syssetting

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/icowan/kit-admin/src/repository/types"
)

type Service interface {
	Add(ctx context.Context, key, val, desc string) (err error)
	Delete(ctx context.Context, key string) (err error)
	Update(ctx context.Context, data *types.SysSetting) (err error)
	Find(ctx context.Context, key string) (res types.SysSetting, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Delete(_ context.Context, key string) (err error) {
	return s.db.Model(&types.SysSetting{}).Where("key = ?", key).Delete(&types.SysSetting{}).Error
}

func (s *service) Update(_ context.Context, data *types.SysSetting) (err error) {
	return s.db.Model(data).Where("id = ?", data.Id).Update(data).Error
}

func (s *service) Find(_ context.Context, key string) (res types.SysSetting, err error) {
	err = s.db.Model(&types.SysSetting{}).Where("`key` = ?", key).First(&res).Error
	return
}

func (s *service) Add(_ context.Context, key, val, desc string) (err error) {
	return s.db.Model(&types.SysSetting{}).Save(&types.SysSetting{
		Key:         key,
		Value:       val,
		Description: desc,
	}).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
