/**
 * @Time : 3/10/21 3:33 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package sysrole

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/icowan/kit-admin/src/repository/types"
)

type Middleware func(next Service) Service

type Service interface {
	// 角色列表
	List(ctx context.Context, page, pageSize int) (res []types.SysRole, total int, err error)
	// 保存角色
	Save(ctx context.Context, data *types.SysRole) (err error)
	// 根据ID查询角色列表
	FindByIds(ctx context.Context, ids []int64) (res []types.SysRole, err error)
	// 根据ID查询角色信息
	Find(ctx context.Context, id int64) (res types.SysRole, err error)
	// 更新角色权限
	AddPermissions(ctx context.Context, role *types.SysRole, permissions []types.SysPermission) (err error)
	// 删除角色
	Delete(ctx context.Context, id int64) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Delete(ctx context.Context, id int64) (err error) {
	// 删除角色相关的权限
	role, err := s.Find(ctx, id)
	if err != nil {
		return err
	}
	tx := s.db.Model(&role).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Association("SysPermissions").Clear().Error; err != nil {
		return err
	}

	// 删除角色相关的用户
	if err = tx.Table("sys_user_roles").Where("role_id = ?", id).
		Unscoped().Delete(nil).Error; err != nil {
		return err
	}

	if err = tx.Where("id = ? ", id).Delete(&role).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (s *service) AddPermissions(ctx context.Context, role *types.SysRole, permissions []types.SysPermission) (err error) {
	tx := s.db.Model(role).Begin()
	if err = tx.Association("SysPermissions").Clear().Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Association("SysPermissions").Append(permissions).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *service) Find(ctx context.Context, id int64) (res types.SysRole, err error) {
	err = s.db.Model(&types.SysRole{}).Where("id = ?", id).First(&res).Error
	return
}

func (s *service) FindByIds(ctx context.Context, ids []int64) (res []types.SysRole, err error) {
	err = s.db.Model(&types.SysRole{}).
		Where("id IN (?)", ids).
		Find(&res).Error
	return
}

func (s *service) Save(ctx context.Context, data *types.SysRole) (err error) {
	data.SysPermissions = nil
	return s.db.Model(data).Save(data).Error
}

func (s *service) List(ctx context.Context, page, pageSize int) (res []types.SysRole, total int, err error) {
	err = s.db.Model(&types.SysRole{}).Order("updated_at DESC").
		Count(&total).
		Offset((page - 1) * pageSize).Limit(total).Find(&res).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
