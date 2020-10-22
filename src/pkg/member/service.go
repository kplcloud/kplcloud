/**
 * @Time : 2019-07-17 13:42
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package member

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/casbin"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	"gopkg.in/guregu/null.v3"
	"strconv"
)

var (
	ErrMemberNamespaceGet = errors.New("Namespace 获取错误")
	ErrMemberGet          = errors.New("用户获取错误")
	ErrMemberList         = errors.New("用户列表获取错误")
	ErrMemberRoleGet      = errors.New("用户角色信息获取错误")
	ErrMemberCount        = errors.New("用户统计错误")
)

type Service interface {
	// Deprecated: please explicitly pick a version if possible.
	Namespaces(ctx context.Context) (list []map[string]string, err error)
	// Deprecated: please explicitly pick a version if possible.
	Detail(ctx context.Context, id int64) (map[string]interface{}, error)
	// Deprecated: please explicitly pick a version if possible.
	Post(ctx context.Context, username, email, password string, state int64, namespaces []string, roleIds []int64) error
	// Deprecated: please explicitly pick a version if possible.
	Update(ctx context.Context, id int64, username, email, password string, state int64, namespaces []string, roleIds []int64) error
	// Deprecated: please explicitly pick a version if possible.
	List(ctx context.Context, page, limit int, email string) (map[string]interface{}, error)
	// Deprecated: please explicitly pick a version if possible.
	Me(ctx context.Context) (map[string]interface{}, error)
	All(ctx context.Context) (interface{}, error)
}

type service struct {
	logger     log.Logger
	config     *config.Config
	casbin     casbin.Casbin
	repository repository.Repository
}

func NewService(logger log.Logger,
	config *config.Config,
	casbin casbin.Casbin,
	store repository.Repository) Service {

	return &service{logger,
		config,
		casbin,
		store,
	}
}

/**
 * @Title 当前用户信息
 */
func (c *service) Me(ctx context.Context) (map[string]interface{}, error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	email := ctx.Value(middleware.EmailContext).(string)
	namespaces := ctx.Value(middleware.NamespacesContext).([]string)

	member, err := c.repository.Member().FindById(memberId)
	if err != nil {
		_ = level.Error(c.logger).Log("memberRepository", "FindById", "err", err.Error())
		return nil, ErrMemberGet
	}

	roles, err := c.repository.Member().GetRolesByMemberId(memberId)
	if err != nil {
		_ = level.Error(c.logger).Log("memberRepository", "GetRolesByMemberId", "err", err.Error())
		return nil, ErrMemberRoleGet
	}

	count, err := c.repository.NoticeMember().CountRead(map[string]string{
		"is_read": "0",
	}, memberId)
	if err != nil {
		_ = level.Error(c.logger).Log("noticeRepository", "CountRead", "err", err.Error())
	}

	return map[string]interface{}{
		"notify_count": count,
		"id":           memberId,
		"email":        email,
		//"state":        member.State,
		"attrs":      "", // 从ldap获取 如果启用了ldap的话
		"namespaces": namespaces,
		"roles":      roles,
		"username":   member.Username,
	}, nil
}

/**
 * @Title 获取当前用户的业务线信息
 */
func (c *service) Namespaces(ctx context.Context) (list []map[string]string, err error) {
	namespaces := ctx.Value(middleware.NamespacesContext).([]string)

	nsList, err := c.repository.Namespace().FindByNames(namespaces)
	if err != nil {
		_ = level.Error(c.logger).Log("namespaceRepository", "FindByNames", "err", err.Error())
		return nil, ErrMemberNamespaceGet
	}

	for _, v := range nsList {
		list = append(list, map[string]string{
			"name":         v.Name,
			"display_name": v.DisplayName,
		})
	}

	return
}

/**
 * @Title 获取用户的详细信息
 */
func (c *service) Detail(ctx context.Context, id int64) (map[string]interface{}, error) {
	member, err := c.repository.Member().FindById(id)
	if err != nil {
		_ = level.Error(c.logger).Log("memberRepository", "FindById", "err", err.Error())
		return nil, ErrMemberGet
	}

	roles, err := c.repository.Member().GetRolesByMemberId(member.ID)
	if err != nil {
		_ = level.Error(c.logger).Log("memberRepository", "GetRolesByMemberId", "err", err.Error())
	}

	return map[string]interface{}{
		"id":         member.ID,
		"username":   member.Username,
		"email":      member.Email,
		"state":      member.State,
		"roles":      roles,
		"namespaces": member.Namespaces,
	}, nil
}

/**
 * @Title 创建用户
 */
func (c *service) Post(ctx context.Context, username, email, password string, state int64, namespaces []string, roleIds []int64) error {
	nsList, err := c.repository.Namespace().FindByNames(namespaces)
	if err != nil {
		_ = level.Error(c.logger).Log("namespaceRepository", "FindByNames", "err", err.Error())
		return ErrMemberNamespaceGet
	}

	roles, err := c.repository.Role().FindByIds(roleIds)
	if err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "FindByIds", "err", err.Error())
		return ErrMemberRoleGet
	}

	var roleList []types.Role
	for _, v := range roles {
		roleList = append(roleList, *v)
	}
	var namespaceList []types.Namespace
	for _, v := range nsList {
		namespaceList = append(namespaceList, *v)
	}

	member := &types.Member{
		Email:      email,
		Username:   username,
		Password:   null.StringFrom(encode.EncodePassword(password, c.config.GetString("server", "app_key"))),
		Roles:      roleList,
		Namespaces: namespaceList,
	}

	return c.repository.Member().CreateMember(member)
}

/**
 * @Title 更新用户
 */
func (c *service) Update(ctx context.Context, id int64, username, email, password string, state int64, namespaces []string, roleIds []int64) error {
	member, err := c.repository.Member().FindById(id)
	if err != nil {
		_ = level.Error(c.logger).Log("memberRepository", "FindById", "err", err.Error())
		return ErrMemberGet
	}

	nsList, err := c.repository.Namespace().FindByNames(namespaces)
	if err != nil {
		_ = level.Error(c.logger).Log("namespaceRepository", "FindByNames", "err", err.Error())
		return ErrMemberNamespaceGet
	}

	roles, err := c.repository.Role().FindByIds(roleIds)
	if err != nil {
		_ = level.Error(c.logger).Log("roleRepository", "FindByIds", "err", err.Error())
		return ErrMemberRoleGet
	}

	var roleList []types.Role
	for _, v := range roles {
		roleList = append(roleList, *v)
	}

	var namespaceList []types.Namespace
	for _, v := range nsList {
		namespaceList = append(namespaceList, *v)
	}

	member.Username = username
	member.Email = email
	if password != "" {
		member.Password = null.StringFrom(encode.EncodePassword(password, c.config.GetString("server", "app_key")))
	}
	member.State = state
	member.Namespaces = namespaceList
	member.Roles = roleList

	c.casbin.GetEnforcer().DeleteRolesForUser(strconv.Itoa(int(member.ID)))
	for _, role := range roleList {
		if _, err = c.casbin.GetEnforcer().AddGroupingPolicySafe(strconv.Itoa(int(member.ID)), strconv.Itoa(int(role.ID))); err != nil {
			_ = level.Error(c.logger).Log("GetEnforcer", "AddGroupingPolicySafe", "err", err.Error())
		}
	}
	_ = c.casbin.GetEnforcer().LoadPolicy()

	return c.repository.Member().UpdateMember(member)
}

/**
 * @Title 用户列表
 */
func (c *service) List(ctx context.Context, page, limit int, email string) (map[string]interface{}, error) {
	count, err := c.repository.Member().Count(email)
	if err != nil {
		_ = level.Error(c.logger).Log("memberRepository", "Count", "err", err.Error())
		return nil, ErrMemberCount
	}
	p := paginator.NewPaginator(page, limit, int(count))

	members, err := c.repository.Member().FindOffsetLimit(p.Offset(), limit, email)
	if err != nil {
		_ = level.Error(c.logger).Log("memberRepository", "FindOffsetLimit", "err", err.Error())
		return nil, ErrMemberList
	}

	return map[string]interface{}{
		"list": members,
		"page": p.Result(),
	}, nil
}

/**
 * @Title 获取所有用户列表
 */
func (c *service) All(ctx context.Context) (interface{}, error) {
	res, err := c.repository.Member().GetMembersAll()
	return res, err
}
