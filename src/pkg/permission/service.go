/**
 * @Time : 2019-07-11 17:19
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package permission

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/casbin"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"gopkg.in/guregu/null.v3"
	"strconv"
)

var (
	ErrPermissionDelete   = errors.New("删除失败,请重试")
	ErrPermissionGet      = errors.New("查无此记录")
	ErrPermissionCreate   = errors.New("添加失败,请重试")
	ErrPermissionUpdate   = errors.New("更新失败,请重试")
	ErrPermissionExists   = errors.New("该地址可能已存在")
	ErrPermissionDragGet  = errors.New("源目标可能不存在")
	ErrPermissionDropGet  = errors.New("目标可能不存在")
	ErrPermissionRole     = errors.New("当前用户角色获取错误")
	ErrPermissionMenusGet = errors.New("获取菜单错误,请联系管理员配置")
)

type Service interface {
	// 删除
	Delete(ctx context.Context, id int64) ([]*types.Permission, error)

	// 更新
	Update(ctx context.Context, id int64, icon, keyType string, menu bool, name, path, method string) ([]*types.Permission, error)

	// 创建Permission
	Post(ctx context.Context, name, path, method, icon string, isMenu bool, parentId int64) error

	// 移动Permission
	Drag(ctx context.Context, dragKey, dropKey int64) (res []*types.Permission, err error)

	// 获取当前用户的菜单
	Menu(ctx context.Context) (res []*types.Permission, err error)

	// 所有列表
	List(ctx context.Context) (res []*types.Permission, err error)
}

type service struct {
	logger     log.Logger
	casbin     casbin.Casbin
	repository repository.Repository
}

func NewService(logger log.Logger,
	casbin casbin.Casbin, repository repository.Repository) Service {
	return &service{logger,
		casbin,
		repository}
}

func (c *service) Delete(ctx context.Context, id int64) ([]*types.Permission, error) {
	if err := c.repository.Permission().Delete(id); err != nil {
		_ = level.Error(c.logger).Log("permissionRepository", "Delete", "err", err.Error())
		return nil, ErrPermissionDelete
	}

	return c.repository.Permission().FindAll()
}

func (c *service) Update(ctx context.Context, id int64, icon, keyType string, menu bool, name, path, method string) ([]*types.Permission, error) {

	permission := &types.Permission{
		Path:   path,
		Name:   name,
		Method: null.StringFrom(method),
		Menu:   null.BoolFromPtr(&menu),
		Icon:   null.StringFrom(icon),
	}

	if keyType == "1" {
		permis, err := c.repository.Permission().FindById(id)
		if err != nil {
			_ = level.Error(c.logger).Log("permissionRepository", "FindById", "err", err.Error())
			return nil, ErrPermissionGet
		}
		permission.ParentID = permis.ParentID
		if err = c.repository.Permission().Create(permission); err != nil {
			_ = level.Error(c.logger).Log("permissionRepository", "Create", "err", err.Error())
			return nil, ErrPermissionCreate
		}
	} else if keyType == "2" {
		//permission.ParentID =
		permission.ParentID = null.IntFromPtr(&id)
		if err := c.repository.Permission().Create(permission); err != nil {
			_ = level.Error(c.logger).Log("permissionRepository", "Create", "err", err.Error())
			return nil, ErrPermissionCreate
		}
	} else if keyType == "4" {
		permission, err := c.repository.Permission().FindById(id)
		if err != nil {
			_ = level.Error(c.logger).Log("permissionRepository", "FindById", "err", err.Error())
			return nil, ErrPermissionGet
		}
		permission.Path = path
		permission.Name = name
		permission.Method = null.StringFrom(method)
		permission.Menu = null.BoolFromPtr(&menu)
		permission.Icon = null.StringFrom(icon)
		if err = c.repository.Permission().Update(permission); err != nil {
			_ = level.Error(c.logger).Log("permissionRepository", "Update", "err", err.Error())
			return nil, ErrPermissionUpdate
		}
	}

	go func() {
		if _, err := c.casbin.GetEnforcer().AddPolicySafe("1", permission.Path, permission.Method.String); err != nil {
			_ = level.Error(c.logger).Log("enforcer", "AddPolicySafe", "err", err.Error())
		}
		if err := c.repository.Role().AddRolePermission(&types.Role{ID: 1}, permission); err != nil {
			_ = level.Error(c.logger).Log("roleRepository", "AddRolePermission", "err", err.Error())
		}
	}()

	return c.repository.Permission().FindAll()
}

func (c *service) Post(ctx context.Context, name, path, method, icon string, isMenu bool, parentId int64) error {
	permission, err := c.repository.Permission().FindByPathAndMethod(path, method)
	if err == nil && permission != nil {
		return ErrPermissionExists
	}

	if err = c.repository.Permission().Create(&types.Permission{
		Path:     path,
		Method:   null.StringFrom(method),
		Icon:     null.StringFrom(icon),
		Name:     name,
		ParentID: null.IntFromPtr(&parentId),
		Menu:     null.BoolFromPtr(&isMenu),
	}); err != nil {
		_ = level.Error(c.logger).Log("permissionRepository", "Create", "err", err.Error())
		return ErrPermissionCreate
	}

	if !c.casbin.GetEnforcer().AddPolicy("1", path, method) {
		_ = level.Error(c.logger).Log("GetEnforcer", "AddPolicy", "bool", false)
	}

	return nil
}

func (c *service) Drag(ctx context.Context, dragKey, dropKey int64) (res []*types.Permission, err error) {
	dragPerm, err := c.repository.Permission().FindById(dragKey)
	if err != nil {
		_ = level.Error(c.logger).Log("permissionRepository", "FindById", "err", err.Error())
		return nil, ErrPermissionDragGet
	}
	dropPerm, err := c.repository.Permission().FindById(dropKey)
	if err != nil {
		_ = level.Error(c.logger).Log("permissionRepository", "FindById", "err", err.Error())
		return nil, ErrPermissionDropGet
	}

	dragPerm.ParentID = null.IntFromPtr(&dropPerm.ID)

	if err = c.repository.Permission().Update(dragPerm); err != nil {
		_ = level.Error(c.logger).Log("permissionRepository", "Update", "err", err.Error())
		return nil, ErrPermissionUpdate
	}

	return c.repository.Permission().FindAll()
}

func (c *service) Menu(ctx context.Context) (res []*types.Permission, err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)

	menus, err := c.repository.Permission().FindMenus()
	if err != nil {
		_ = level.Error(c.logger).Log("permissionRepository", "FindMenus", "err", err.Error())
		return nil, ErrPermissionMenusGet
	}
	roles, err := c.casbin.GetEnforcer().GetRolesForUser(strconv.Itoa(int(userId)))
	if err != nil {
		_ = level.Error(c.logger).Log("GetEnforcer", "GetRolesForUser", "err", err.Error())
		return nil, ErrPermissionRole
	}

	// 有一种空间换时间的方案可以试试

	var perms []string
	permMap := make(map[string]bool)
	for _, role := range roles {
		for _, val := range c.casbin.GetEnforcer().GetPermissionsForUser(role) {
			if len(val) > 1 {
				permMap[val[1]] = true
			}
		}
	}

	for key, _ := range permMap {
		perms = append(perms, key)
	}

	var returnMenus MenuSort

	for _, menu := range menus {
		for perm, _ := range permMap {
			if perm == menu.Path {
				returnMenus = append(returnMenus, menu)
			}
		}
	}

	return returnMenus, nil
}

func (c *service) List(ctx context.Context) (res []*types.Permission, err error) {
	return c.repository.Permission().FindAll()
}

type MenuSort []*types.Permission

func (m MenuSort) Len() int {
	return len(m)
}
func (m MenuSort) Less(i, j int) bool {
	return m[i].ID < m[j].ID
}
func (m MenuSort) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
