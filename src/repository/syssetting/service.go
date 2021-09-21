/**
 * @Time: 2020/9/26 19:09
 * @Author: solacowa@gmail.com
 * @File: setting
 * @Software: GoLand
 */

package syssetting

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Service interface {
	Add(ctx context.Context, section, key, val, desc string) (err error)
	Delete(ctx context.Context, section, key string) (err error)
	Update(ctx context.Context, data *types.SysSetting) (err error)
	Find(ctx context.Context, section, key string) (res types.SysSetting, err error)
	FindAll(ctx context.Context) (res []types.SysSetting, err error)
	List(ctx context.Context, key string, page, pageSize int) (res []types.SysSetting, total int, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) List(ctx context.Context, key string, page, pageSize int) (res []types.SysSetting, total int, err error) {
	query := s.db.Model(&types.SysSetting{})
	if !strings.EqualFold(key, "") {
		query = query.Where("key = ?", key)
	}
	err = query.Count(&total).
		Order("id,section DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&res).Error
	return
}

func (s *service) FindAll(ctx context.Context) (res []types.SysSetting, err error) {
	err = s.db.Model(&types.SysSetting{}).Find(&res).Error
	return
}

func (s *service) Delete(_ context.Context, section, key string) (err error) {
	return s.db.Model(&types.SysSetting{}).Where("key = ?", key).Delete(&types.SysSetting{}).Error
}

func (s *service) Update(_ context.Context, data *types.SysSetting) (err error) {
	return s.db.Model(data).Where("id = ?", data.Id).Update(data).Error
}

func (s *service) Find(_ context.Context, section, key string) (res types.SysSetting, err error) {
	err = s.db.Model(&types.SysSetting{}).Where("`key` = ?", key).First(&res).Error
	return
}

func (s *service) Add(ctx context.Context, section, key, val, desc string) (err error) {
	st, err := s.Find(ctx, section, key)
	if err != nil {
		st = types.SysSetting{}
	}
	st.Section = section
	st.Key = key
	st.Value = val
	st.Description = desc
	return s.db.Model(&st).Save(&st).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
