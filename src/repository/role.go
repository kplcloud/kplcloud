/**
 * @Time : 2019-07-12 16:37
 * @Author : solacowa@gmail.com
 * @File : role
 * @Software: GoLand
 */

package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type RoleRepository interface {
	FindByIds(ids []int64) (res []*types.Role, err error)
	FindById(id int64) (*types.Role, error)
	FindPermission(id int64) (*types.Role, error)
	AddRolePermission(role *types.Role, permission ...*types.Permission) error
	Create(role *types.Role) error
	Update(role *types.Role) error
	FindAll() (roles []*types.Role, err error)
	Delete(id int64) error
	DeletePermission(role *types.Role) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (c *roleRepository) FindByIds(ids []int64) (res []*types.Role, err error) {
	err = c.db.Find(&res, "id in (?)", ids).Error
	return
}

func (c *roleRepository) FindById(id int64) (*types.Role, error) {
	var role types.Role
	err := c.db.First(&role, "id = ?", id).Error
	return &role, err
}

func (c *roleRepository) FindPermission(id int64) (*types.Role, error) {
	var res types.Role
	err := c.db.Preload("Permissions").First(&res, id).Error

	return &res, err
}

func (c *roleRepository) AddRolePermission(role *types.Role, permission ...*types.Permission) error {
	err := c.db.Model(role).Association("Permissions").Append(permission).Error
	return err
}

func (c *roleRepository) Create(role *types.Role) error {
	return c.db.Create(role).Error
}

func (c *roleRepository) Update(role *types.Role) error {
	return c.db.Model(role).Where("id = ?", role.ID).Update(role).Error
}

func (c *roleRepository) FindAll() (roles []*types.Role, err error) {
	err = c.db.Find(&roles).Error
	return
}

func (c *roleRepository) Delete(id int64) error {
	role, err := c.FindPermission(id)
	if err != nil {
		return err
	}

	tx := c.db.Begin()

	if role.ID < 1 {
		return errors.New("Role can not null ")
	}

	err = tx.Model(&role).Association("Permissions").Clear().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&role).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (c *roleRepository) DeletePermission(role *types.Role) error {
	return c.db.Model(&role).Association("Permissions").Clear().Error
}
