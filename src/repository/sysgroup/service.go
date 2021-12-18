/**
 * @Time : 2021/12/17 11:15 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package sysgroup

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"strings"
)

type Call func() error

type Middleware func(Service) Service

type Service interface {
	Save(ctx context.Context, group *types.SysGroup, call Call) (err error)
	FindByName(ctx context.Context, clusterId int64, namespace, name string) (res types.SysGroup, err error)
	List(ctx context.Context, clusterId int64, groupIds []int64, namespace, name string, page, pageSize int) (res []types.SysGroup, total int, err error)
	FindIds(ctx context.Context, clusterId int64, namespace string) (ids []int64, err error)
	Delete(ctx context.Context, group *types.SysGroup, call ...Call) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Delete(ctx context.Context, group *types.SysGroup, call ...Call) (err error) {
	return s.db.Model(group).Transaction(func(tx *gorm.DB) error {
		if err = tx.Association("Users").Clear().Error; err != nil {
			return err
		}
		if err = tx.Association("Apps").Clear().Error; err != nil {
			return err
		}
		if err = tx.Where("id = ?", group.Id).Delete(group).Error; err != nil {
			return err
		}
		if len(call) > 0 {
			for _, v := range call {
				if err = v(); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *service) FindIds(ctx context.Context, clusterId int64, namespace string) (ids []int64, err error) {
	type idStruct struct {
		Id int64 `json:"id"`
	}
	var res []idStruct
	err = s.db.Model(&types.SysGroup{}).
		Where("cluster_id = ? AND namespace = ?", clusterId, namespace).Find(&res).Error

	if err == nil {
		for _, v := range res {
			ids = append(ids, v.Id)
		}
	}

	return
}

func (s *service) List(ctx context.Context, clusterId int64, groupIds []int64, namespace, name string, page, pageSize int) (res []types.SysGroup, total int, err error) {
	q := s.db.Model(&types.SysGroup{}).
		Preload("User").
		Where("cluster_id = ? AND id IN (?) AND namespace = ?", clusterId, groupIds, namespace)
	if !strings.EqualFold(name, "") {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}

	err = q.Count(&total).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&res).Error
	return
}

func (s *service) FindByName(ctx context.Context, clusterId int64, namespace, name string) (res types.SysGroup, err error) {
	err = s.db.Model(&res).Where("cluster_id = ? AND namespace = ? AND name = ?", clusterId, namespace, name).First(&res).Error
	return
}

func (s *service) Save(ctx context.Context, group *types.SysGroup, call Call) (err error) {
	return s.db.Model(group).Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(group).Error; err != nil {
			return err
		}
		if call != nil {
			return call()
		}
		return nil
	})
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
