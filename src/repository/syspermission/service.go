/**
 * @Time : 5/25/21 3:51 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package syspermission

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(next Service) Service

type Service interface {
	Save(ctx context.Context, data *types.SysPermission) (err error)
	Delete(ctx context.Context, id int64) (err error)
	FindByIds(ctx context.Context, ids []int64) (res []types.SysPermission, err error)
	Find(ctx context.Context, id int64) (res types.SysPermission, err error)
	FindAll(ctx context.Context) (res []types.SysPermission, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindAll(ctx context.Context) (res []types.SysPermission, err error) {
	err = s.db.Model(&res).Find(&res).Error
	return
}

func (s *service) Find(ctx context.Context, id int64) (res types.SysPermission, err error) {
	err = s.db.Model(&res).Where("id = ?", id).First(&res).Error
	return
}

func (s *service) FindByIds(ctx context.Context, ids []int64) (res []types.SysPermission, err error) {
	err = s.db.Model(&types.SysPermission{}).Where("id IN (?)", ids).Find(&res).Error
	return
}

func (s *service) Delete(ctx context.Context, id int64) (err error) {
	var list []types.SysPermission
	if err = s.db.Model(&types.SysPermission{}).Where("parent_id = ?", id).Find(&list).Error; err != nil {
		return err
	}
	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	for _, v := range list {
		if err = tx.Where("id = ?", v.Id).Delete(v).Error; err != nil {
			break
		}
	}

	if err = tx.Where("id = ?", id).Delete(&types.SysPermission{Id: id}).Error; err != nil {
		return err
	}

	// 删除角色下的这个权限
	if err = tx.Table("sys_role_permissions").Where("permission_id = ?", id).
		Unscoped().Delete(nil).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (s *service) Save(ctx context.Context, data *types.SysPermission) (err error) {
	return s.db.Model(data).Save(data).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
