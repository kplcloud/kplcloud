/**
 * @Time : 3/4/21 5:25 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package sysuser

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Service interface {
	// Save 保存用户
	Save(ctx context.Context, user *types.SysUser) (err error)
	// FindByEmail 根据邮箱查用户
	FindByEmail(ctx context.Context, email string) (res types.SysUser, err error)
	// List 系统用户列表
	List(ctx context.Context, namespaceIds []int64, email string, page, pageSize int) (res []types.SysUser, total int, err error)
	// Find 查询用户详情
	Find(ctx context.Context, userId int64, preload ...string) (res types.SysUser, err error)
	// AddRoles 给用户添加角色
	AddRoles(ctx context.Context, user *types.SysUser, roles []types.SysRole) (err error)
	// FindByIds 根据ID查用户
	FindByIds(ctx context.Context, ids []int64, preload ...string) (res []types.SysUser, err error)
	// FindByRoleId 根据角色查询用户列表
	FindByRoleId(ctx context.Context, roleId int64, page, pageSize int) (res []types.SysUser, total int, err error)
	// Delete 删除用户
	Delete(ctx context.Context, userId int64, unscoped bool) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Delete(ctx context.Context, userId int64, unscoped bool) (err error) {
	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if unscoped {
		if err = tx.Table("sys_user_namespaces").Where("sys_user_id = ?", userId).
			Unscoped().Delete(nil).Error; err != nil {
			return err
		}
		if err = tx.Table("sys_user_roles").Where("sys_user_id = ?", userId).
			Unscoped().Delete(nil).Error; err != nil {
			return err
		}
		return tx.Model(&types.SysUser{Id: userId}).Where("id = ?", userId).
			Unscoped().Delete(&types.SysUser{}).Error
	}

	err = tx.Model(&types.SysUser{Id: userId}).Where("id = ?", userId).Delete(&types.SysUser{}).Error

	return
}

func (s *service) FindByRoleId(ctx context.Context, roleId int64, page, pageSize int) (res []types.SysUser, total int, err error) {
	panic("implement me")
}

func (s *service) FindByIds(ctx context.Context, ids []int64, preload ...string) (res []types.SysUser, err error) {
	query := s.db.Model(&types.SysUser{}).Where("id IN (?)", ids)
	if len(preload) > 0 {
		for _, v := range preload {
			query = query.Preload(v)
		}
	}
	err = query.Find(&res).Error
	return
}

func (s *service) AddRoles(ctx context.Context, user *types.SysUser, roles []types.SysRole) (err error) {
	tx := s.db.Model(user).Begin()
	if err = tx.Association("SysRoles").Clear().Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Association("SysRoles").Append(roles).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *service) Find(ctx context.Context, userId int64, preload ...string) (res types.SysUser, err error) {
	query := s.db.Model(&types.SysUser{}).Where("id = ?", userId)
	if len(preload) > 0 {
		for _, v := range preload {
			query = query.Preload(v)
		}
	}
	err = query.First(&res).Error
	return
}

func (s *service) List(ctx context.Context, namespaceIds []int64, email string, page, pageSize int) (res []types.SysUser, total int, err error) {
	query := s.db.Model(&types.SysUser{}).Preload("SysRoles")

	if namespaceIds != nil {
		query = query.Preload("SysNamespaces", func(db *gorm.DB) *gorm.DB {
			return db.Where("id IN (?)", namespaceIds)
		})
	}
	if !strings.EqualFold(email, "") {
		query = query.Where("email LIKE '%?%'", email)
	}

	err = query.Count(&total).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error

	return
}

func (s *service) FindByEmail(ctx context.Context, email string) (res types.SysUser, err error) {
	query := s.db.Model(&types.SysUser{}).
		Preload("SysRoles.SysPermissions").
		Preload("SysNamespaces")
	if strings.Contains(email, "@") {
		query = query.Where("email = ?", email)
	} else {
		query = query.Where("login_name = ?", email)
	}
	err = query.Find(&res).Error
	return
}

func (s *service) Save(ctx context.Context, user *types.SysUser) (err error) {
	user.SysRoles = nil
	return s.db.Model(user).Save(user).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
