/**
 * @Time : 2019-07-16 17:51
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package role

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/casbin"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"gopkg.in/guregu/null.v3"
	"strconv"
)

var (
	ErrRoleGet              = errors.New("获取角色错误,可能不存在")
	ErrRoleCreate           = errors.New("角色创建错误")
	ErrRoleUpdate           = errors.New("角色更新错误")
	ErrRolePermissionGet    = errors.New("获取角色的权限列表错误")
	ErrRolePermission       = errors.New("角色授权错误")
	ErrRolePermissionDelete = errors.New("角色授权关系清除错误")
	// ErrRoleDelete           = errors.New("角色删除错误")
)

type Service interface {
	// 获取已经选择的PermID
	PermissionSelected(ctx context.Context, id int64) (ids []int64, err error)

	// 角色详情
	Detail(ctx context.Context, id int64) (role *types.Role, err error)

	// 创建角色
	Post(ctx context.Context, name, desc string, level int) error

	// 更新角色信息
	Update(ctx context.Context, id int64, name, desc string, level int) error

	// 取得所有角色
	All(ctx context.Context) ([]*types.Role, error)

	// 删除角色
	Delete(ctx context.Context, id int64) error

	// 给角色分配权限
	RolePermission(ctx context.Context, id int64, permIds []int64) error
}

type service struct {
	logger     log.Logger
	casbin     casbin.Casbin
	repository repository.Repository
}

func NewService(logger log.Logger, casbin casbin.Casbin, repository repository.Repository) Service {
	return &service{logger, casbin, repository}
}

func (c *service) PermissionSelected(ctx context.Context, id int64) (ids []int64, err error) {
	perms, err := c.repository.Role().FindPermission(id)
	if err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "FindPermission", "err", err.Error())
		return nil, ErrRolePermissionGet
	}

	for _, val := range perms.Permissions {
		ids = append(ids, val.ID)
	}

	return
}

func (c *service) Detail(ctx context.Context, id int64) (role *types.Role, err error) {
	role, err = c.repository.Role().FindById(id)
	if err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "FindById", "err", err.Error())
		return nil, ErrRoleGet
	}
	return
}

func (c *service) Post(ctx context.Context, name, desc string, lv int) error {
	if err := c.repository.Role().Create(&types.Role{
		Name:        name,
		Description: desc,
		Level:       lv,
	}); err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "Create", "err", err.Error())
		return ErrRoleCreate
	}
	return nil
}

func (c *service) Update(ctx context.Context, id int64, name, desc string, lv int) error {
	role, err := c.repository.Role().FindById(id)
	if err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "FindById", "err", err.Error())
		return ErrRoleGet
	}

	role.Name = name
	role.Description = desc
	role.Level = lv
	if err := c.repository.Role().Update(role); err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "Update", "err", err.Error())
		return ErrRoleUpdate
	}
	return nil
}

func (c *service) All(ctx context.Context) ([]*types.Role, error) {
	return c.repository.Role().FindAll()
}

func (c *service) Delete(ctx context.Context, id int64) error {
	return c.repository.Role().Delete(id)
}

func (c *service) RolePermission(ctx context.Context, id int64, permIds []int64) error {
	role, err := c.repository.Role().FindById(id)
	if err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "FindById", "err", err.Error())
		return ErrRoleGet
	}

	if err = c.repository.Role().DeletePermission(role); err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "DeletePermission", "err", err.Error())
		return ErrRolePermissionDelete
	}

	perms, err := c.repository.Permission().FindByIds(permIds)
	if err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "FindByIds", "err", err.Error())
		return ErrRolePermissionGet
	}

	if err = c.repository.Role().AddRolePermission(role, perms...); err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "AddRolePermission", "err", err.Error())
		return ErrRolePermission
	}

	menuIds := make(map[int64]int64)

	c.casbin.GetEnforcer().DeletePermissionsForUser(strconv.Itoa(int(role.ID)))

	for _, perm := range perms {
		if perm.ParentID != null.IntFrom(0) {
			menuIds[perm.ParentID.Int64] = perm.ParentID.Int64
		}
	}
	var ids []int64
	for _, id := range menuIds {
		ids = append(ids, id)
	}

	res, _ := c.repository.Permission().FindByIds(ids)
	perms = append(perms, res...)
	for _, perm := range perms {
		if _, err = c.casbin.GetEnforcer().AddPolicySafe(strconv.Itoa(int(role.ID)), perm.Path, perm.Method.String); err != nil {
			_ = level.Warn(c.logger).Log("GetEnforcer", "AddPolicySafe", "err", err.Error())
		}
	}

	return nil
}
