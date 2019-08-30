/**
 * @Time : 2019-07-11 17:38
 * @Author : solacowa@gmail.com
 * @File : permission
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type PermissionRepository interface {
	FindById(id int64) (*types.Permission, error)
	FindAll() ([]*types.Permission, error)
	Delete(id int64) error
	Create(p *types.Permission) error
	Update(p *types.Permission) error
	FindByPathAndMethod(path, method string) (*types.Permission, error)
	FindByPerm(perm []string) ([]*types.Permission, error)
	FindMenus() (res []*types.Permission, err error)
	FindByIds(ids []int64) (res []*types.Permission, err error)
	//FindByRoleId(roleId int64) (*types.Role, error)
}

type permission struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permission{db}
}

func (c *permission) FindByPathAndMethod(path, method string) (*types.Permission, error) {
	var res types.Permission
	err := c.db.First(&res, "path = ? AND method = ?", path, method).Error
	return &res, err
}

func (c *permission) Create(p *types.Permission) error {
	return c.db.Save(p).Error
}

func (c *permission) Update(p *types.Permission) error {
	return c.db.Model(p).Where("id = ?", p.ID).Update(p).Error
}

func (c *permission) FindById(id int64) (*types.Permission, error) {
	var res types.Permission
	if err := c.db.Find(&res, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *permission) Delete(id int64) error {
	perm := types.Permission{ID: id}
	return c.db.Delete(perm, "id = ?", id).Error
}

func (c *permission) FindAll() ([]*types.Permission, error) {
	var all []*types.Permission
	if err := c.db.Find(&all).Error; err != nil {
		return nil, err
	}

	return all, nil
}

func (c *permission) FindByPerm(perm []string) ([]*types.Permission, error) {

	return nil, nil
}

func (c *permission) FindMenus() (res []*types.Permission, err error) {
	err = c.db.Find(&res, "menu = ?", true).Error
	return
}

func (c *permission) FindByIds(ids []int64) (res []*types.Permission, err error) {
	err = c.db.Find(&res, "id in (?)", ids).Error
	return
}
