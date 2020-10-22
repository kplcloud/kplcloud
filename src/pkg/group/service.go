package group

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/paginator"
)

var (
	ErrNamespaceNotAllowed   = errors.New("没有该业务线的权限")
	ErrGroupCreateFailed     = errors.New("组创建失败")
	ErrGetUserInfoError      = errors.New("获取用户信息失败")
	ErrGetGroupInfoFailed    = errors.New("获取组信息失败")
	ErrGroupNameEnExists     = errors.New("组英文名已存在")
	ErrGroupNameExists       = errors.New("组名已存在")
	ErrGroupCount            = errors.New("获取组数目失败")
	ErrGroupPaginate         = errors.New("获取一些组信息失败")
	ErrMemberNotExist        = errors.New("未找到该用户信息")
	ErrAdminUpdateGroup      = errors.New("管理员更新组信息失败")
	ErrAdminDestroyGroup     = errors.New("删除组和组关联关系失败")
	ErrNsProjectList         = errors.New("获取业务线下项目信息失败")
	ErrNsCronjobList         = errors.New("获取业务线下定时任务信息失败")
	ErrGetProjectInfoFailed  = errors.New("获取项目信息失败")
	ErrAdminAddProjectFailed = errors.New("管理员添加项目失败")
	ErrPNsDiffGNs            = errors.New("项目业务线与组业务线不相同")
	ErrGetCronjobInfoFailed  = errors.New("获取定时任务信息失败")
	ErrCNsDiffGNs            = errors.New("定时任务业务线与组业务线不相同")
	ErrAdminAddCronjobFailed = errors.New("管理员添加定时任务失败")
	ErrMNsDiffGNs            = errors.New("成员没有该组业务线的权限")
	ErrAdminAddMemberFailed  = errors.New("管理员添加成员失败")
	ErrAdminDelMemberFailed  = errors.New("管理员删除成员失败")
	ErrAdminDelProjectFailed = errors.New("管理员删除项目失败")
	ErrAdminDelCronjobFailed = errors.New("管理员删除定时任务失败")
	ErrOwnerUpdateGroup      = errors.New("组长修改组信息失败")
	ErrOwnerAddMemberFailed  = errors.New("组长添加成员失败")
	ErrOwnerDelMemberFailed  = errors.New("组长删除成员失败")
	ErrOwnerAddProjectFailed = errors.New("组长添加项目失败")
	ErrOwnerDelProjectFailed = errors.New("组长添加项目失败")
	ErrOwnerAddCronjobFailed = errors.New("组长添加定时任务失败")
	ErrOwnerDelCronjobFailed = errors.New("组长删除定时任务失败")
	ErrOwnerCanNotDelSelf    = errors.New("不能删除自己")
	ErrOwnerDelGroupFailed   = errors.New("组长删除组失败")
	ErrGetMyGroupListFailed  = errors.New("获取自己组列表信息失败")
	ErrGetMyNsListFailed     = errors.New("获取自己业务线列表信息失败")
	ErrGetGroupRelListFailed = errors.New("获取自己组详细关联关系数据失败")
)

type Service interface {
	// 创建组
	Post(ctx context.Context, gr gRequest) error

	// 获取所有组列表
	GetAll(ctx context.Context, request getAllRequest) (map[string]interface{}, error)

	// 超管添加组
	AdminAddGroup(ctx context.Context, gr gRequest) error

	// 超管修改组信息
	AdminUpdateGroup(ctx context.Context, groupId int64, gr gRequest) error

	// 获取相关用户列表
	GetMemberByEmailLike(ctx context.Context, email string, ns string) ([]types.Member, error)

	// 超管删除组
	AdminDestroy(ctx context.Context, groupId int64) error

	// 业务线下项目列表
	NamespaceProjectList(ctx context.Context, nq nsListRequest) ([]*types.Project, error)

	// 业务线下定时任务列表
	NamespaceCronjobList(ctx context.Context, nq nsListRequest) ([]*types.Cronjob, error)

	// 超管添加项目
	AdminAddProject(ctx context.Context, aq adminDoProjectRequest) error

	// 超管添加定时任务
	AdminAddCronjob(ctx context.Context, ac adminDoCronjobRequest) error

	// 超管添加组用户
	AdminAddMember(ctx context.Context, am adminDoMemberRequest) error

	// 超管删除组用户
	AdminDelMember(ctx context.Context, am adminDoMemberRequest) error

	// 超管删除组项目
	AdminDelProject(ctx context.Context, aq adminDoProjectRequest) error

	// 超管删除组定时任务
	AdminDelCronjob(ctx context.Context, ac adminDoCronjobRequest) error

	// 组长添加组
	OwnerAddGroup(ctx context.Context, oq ownerDoGroup) error

	// 组长修改组
	OwnerUpdateGroup(ctx context.Context, oq ownerDoGroup) error

	// 组长添加组成员
	OwnerAddMember(ctx context.Context, om ownerDoMember) error

	// 组长删除组成员
	OwnerDelMember(ctx context.Context, om ownerDoMember) error

	// 组长添加组项目
	OwnerAddProject(ctx context.Context, op ownerDoProject) error

	// 组长删除组项目
	OwnerDelProject(ctx context.Context, op ownerDoProject) error

	// 组长添加组定时任务
	OwnerAddCronjob(ctx context.Context, oc ownerDoCronjob) error

	// 组长删除组定时任务
	OwnerDelCronjob(ctx context.Context, oc ownerDoCronjob) error

	// 组长删除组
	OwnerDelGroup(ctx context.Context, oq ownerDoGroup) error

	// 用户下组列表
	UserMyList(ctx context.Context, umq userGroupListRequest) ([]*types.Groups, error)

	// 用户下业务线列表
	NsMyList(ctx context.Context) ([]types.Namespace, error)

	// 组详细信息
	RelDetail(ctx context.Context, groupId int64) (*types.Groups, error)

	// 组名是否存在
	GroupNameExists(ctx context.Context, name string) (bool, error)

	// 组别名是否存在
	GroupDisplayNameExists(ctx context.Context, displayName string) (bool, error)
}

type service struct {
	logger     log.Logger
	config     *config.Config
	repository repository.Repository
}

/**
 * 有没有用,有待考证
 */
func (c *service) Post(ctx context.Context, gr gRequest) error {
	userNSArr := ctx.Value(middleware.NamespacesContext).([]string)

	var isTrue bool
	for _, v := range userNSArr {
		if gr.NameSpace == v {
			isTrue = true
		}
	}

	if !isTrue {
		_ = level.Error(c.logger).Log("group", "post", "err", "namespace is not allowed")
		return ErrNamespaceNotAllowed
	}

	g, err := c.repository.Groups().GetGroupByName(gr.Name)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "get group info by name_en", "err", err.Error())
		return ErrGetGroupInfoFailed
	}
	if g.ID > 0 {
		_ = level.Error(c.logger).Log("group", "get group info by name_en", "err", "Group name_en exists ")
		return ErrGroupNameEnExists
	}

	memberId := ctx.Value(middleware.UserIdContext).(int64)
	memberInfo, err := c.repository.Member().GetInfoById(memberId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "get user info by id", "err", err.Error())
		return ErrGetUserInfoError
	}

	if err := c.repository.Groups().CreateGroupAndRelation(&types.Groups{
		Name:        gr.Name,
		DisplayName: gr.DisplayName,
		Namespace:   gr.NameSpace,
		MemberId:    memberId,
	}, memberInfo); err != nil {
		_ = level.Error(c.logger).Log("group", "create", "err", err.Error())
		return ErrGroupCreateFailed
	}

	_ = level.Info(c.logger).Log("group", "create", "message", "succeed to create group")

	return nil
}

func (c *service) GetAll(ctx context.Context, request getAllRequest) (map[string]interface{}, error) {
	cnt, err := c.repository.Groups().AllGroupsCount(request.Name, request.Ns)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "get", "err", err.Error())
		return nil, ErrGroupCount
	}

	p := paginator.NewPaginator(request.Page, request.Limit, int(cnt))
	groups, err := c.repository.Groups().GroupsPaginate(request.Name, request.Ns, p.Offset(), request.Limit)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "get", "err", err.Error())
		return nil, ErrGroupPaginate
	}

	var listMap []interface{}
	for _, v := range groups {
		ns := newNs{
			Id:          v.Ns.ID,
			Name:        v.Ns.Name,
			DisplayName: v.Ns.DisplayName,
		}

		owner := newOwner{
			Id:       v.Member.ID,
			Email:    v.Member.Email,
			Username: v.Member.Username,
		}

		group := map[string]interface{}{
			"id":           v.ID,
			"name":         v.Name,
			"display_name": v.DisplayName,
			"namespace":    ns,
			"owner":        owner,
		}
		listMap = append(listMap, group)
	}

	var returnData = map[string]interface{}{
		"list": listMap,
		"page": map[string]interface{}{
			"total":     cnt,
			"pageTotal": p.PageTotal(),
			"pageSize":  request.Limit,
			"page":      p.Page(),
		},
	}
	return returnData, nil
}

func (c *service) AdminAddGroup(ctx context.Context, gr gRequest) error {
	if _, exists := c.repository.Groups().GroupExistsByName(gr.Name); !exists {
		_ = level.Error(c.logger).Log("group", "AdminAddGroup get group info by name", "err", "Group name exists ")
		return ErrGroupNameExists
	}

	memberInfo, err := c.repository.Member().GetInfoById(gr.MemberId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddGroup get user info by id", "err", err.Error())
		return ErrGetUserInfoError
	}
	if err := c.repository.Groups().CreateGroupAndRelation(&types.Groups{
		Name:        gr.Name,
		DisplayName: gr.DisplayName,
		Namespace:   gr.NameSpace,
		MemberId:    gr.MemberId,
	}, memberInfo); err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddGroup", "err", err.Error())
		return ErrGroupCreateFailed
	}

	_ = level.Info(c.logger).Log("group", "AdminAddGroup", "message", "succeed to create group")

	return nil
}

func (c *service) AdminUpdateGroup(ctx context.Context, groupId int64, gr gRequest) error {

	group, err := c.repository.Groups().Find(groupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminUpdateGroup get group info by id", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	//if group == nil || group.ID < 1 {
	//	_ = level.Error(c.logger).Log("group", "AdminUpdateGroup get group info by id", "err", "group not exists")
	//	return nil
	//}
	//
	if group.Name != gr.Name {
		_, exists := c.repository.Groups().GroupNameExists(gr.Name, group.ID)
		if !exists {
			_ = level.Error(c.logger).Log("group", "AdminUpdateGroup get group info by name", "err", "group name exists")
			return ErrGroupNameExists
		}
	}

	if group.Name != gr.Name {
		_, exists := c.repository.Groups().GroupNameExists(gr.Name, group.ID)
		if !exists {
			_ = level.Error(c.logger).Log("group", "AdminUpdateGroup get group info by name_en", "err", "group name_en exists")
			return ErrGroupNameEnExists
		}
	}

	if group.MemberId != gr.MemberId {
		m, err := c.repository.Member().GetInfoById(gr.MemberId)
		if err != nil {
			_ = level.Error(c.logger).Log("group", "AdminUpdateGroup get member info by id", "err", err.Error())
			return ErrGetUserInfoError
		}
		if m != nil && m.ID < 1 {
			_ = level.Error(c.logger).Log("group", "AdminUpdateGroup get member info by id", "err", "member do not exists")
			return ErrMemberNotExist
		}
		if err = c.repository.Groups().UpdateGroupAndRelation(&types.Groups{
			ID:          groupId,
			Name:        gr.Name,
			DisplayName: gr.DisplayName,
			//Namespace:   gr.NameSpace,
			MemberId: gr.MemberId,
		}, groupId, m); err != nil {
			_ = level.Error(c.logger).Log("group", "AdminUpdateGroup update group info error ", "err", err.Error())
			return ErrAdminUpdateGroup
		}
	} else {
		if err = c.repository.Groups().UpdateGroup(&types.Groups{
			ID:          groupId,
			Name:        gr.Name,
			DisplayName: gr.DisplayName,
			//Namespace:   gr.NameSpace,
			MemberId: gr.MemberId,
		}, groupId); err != nil {
			_ = level.Error(c.logger).Log("group", "AdminUpdateGroup update group info error ", "err", err.Error())
			return ErrAdminUpdateGroup
		}
	}
	_ = level.Info(c.logger).Log("group", "AdminUpdateGroup", "message", " admin succeed to update group")
	return nil
}

func (c *service) GetMemberByEmailLike(ctx context.Context, email string, ns string) ([]types.Member, error) {
	list, err := c.repository.Member().GetMembers(email, ns)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "GetMemberByEmailLike get members info by email like and ns", "err", err.Error())
		return nil, ErrGetUserInfoError
	}
	return list, nil
}

func (c *service) AdminDestroy(ctx context.Context, groupId int64) error {
	group, err := c.repository.Groups().Find(groupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDestroy get group info by group_id", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	err = c.repository.Groups().DestroyAndRelation(group)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDestroy delete group and relation has error", "err", err.Error())
		return ErrAdminDestroyGroup
	}

	return nil
}

func (c *service) NamespaceProjectList(ctx context.Context, nq nsListRequest) (res []*types.Project, err error) {
	//var res []*types.Project
	if nq.Name != "" {
		res, err = c.repository.Project().GetProjectByNameLikeAndNs(nq.Name, nq.Ns)
	} else {
		res, err = c.repository.Project().GetProjectByNs(nq.Ns)
	}
	if err != nil {
		_ = level.Error(c.logger).Log("group", "NamespaceProjectList get namespace and nameLike project list has error", "err", err.Error())
		return nil, ErrNsProjectList
	}
	return res, nil
}

func (c *service) NamespaceCronjobList(ctx context.Context, nq nsListRequest) (res []*types.Cronjob, err error) {
	if nq.Name != "" {
		res, err = c.repository.CronJob().GetCronjobByNameLikeAndNs(nq.Name, nq.Ns)
	} else {
		res, err = c.repository.CronJob().GetCronjobByNs(nq.Ns)
	}
	if err != nil {
		_ = level.Error(c.logger).Log("group", "NamespaceCronjobList get namespace and nameLike cronjob list has error", "err", err.Error())
		return nil, ErrNsCronjobList
	}
	return
}

func (c *service) AdminAddProject(ctx context.Context, aq adminDoProjectRequest) (err error) {
	group, err := c.repository.Groups().Find(aq.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddProject get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	project, err := c.repository.Project().Find(aq.ProjectId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddProject get project info by project_id has error ", "err", err.Error())
		return ErrGetProjectInfoFailed
	}

	// if the project namespace different from group namespace
	// return error
	if project.Namespace != group.Namespace {
		_ = level.Error(c.logger).Log("group", "AdminAddProject project namespace not the same of group namespace ", "err", " project namespace different from the group namespace ")
		return ErrPNsDiffGNs
	}

	err = c.repository.Groups().GroupAddProject(group, project)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddProject add project failed ", "err", err.Error())
		return ErrAdminAddProjectFailed
	}
	return nil
}

func (c *service) AdminAddCronjob(ctx context.Context, ac adminDoCronjobRequest) error {
	group, err := c.repository.Groups().Find(ac.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddCronjob get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	cronjob, err := c.repository.CronJob().Find(ac.CronjobId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddCronjob get cronjob info by cronjob_id has error ", "err", err.Error())
		return ErrGetCronjobInfoFailed
	}

	// if the cronjob namespace different from group namespace
	// return error
	if cronjob.Namespace != group.Namespace {
		_ = level.Error(c.logger).Log("group", "AdminAddProject cronjob namespace not the same of group namespace ", "err", " cronjob namespace different from the group namespace ")
		return ErrCNsDiffGNs
	}

	err = c.repository.Groups().GroupAddCronjob(group, cronjob)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddCronjob add cronjob failed ", "err", err.Error())
		return ErrAdminAddCronjobFailed
	}
	return nil
}

func (c *service) AdminAddMember(ctx context.Context, am adminDoMemberRequest) error {
	group, err := c.repository.Groups().Find(am.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddMember get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	member, err := c.repository.Member().FindById(am.MemberId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddMember get member info by member_id has error ", "err", err.Error())
		return ErrGetUserInfoError
	}

	var isTrue bool
	for _, v := range member.Namespaces {
		if v.Name == group.Namespace {
			isTrue = true
		}
	}
	// if the member namespace different from group namespace
	// return error
	if !isTrue {
		_ = level.Error(c.logger).Log("group", "AdminAddMember member namespace not the same of group namespace ", "err", " member namespace different from the group namespace ")
		return ErrMNsDiffGNs
	}

	err = c.repository.Groups().GroupAddMember(group, member)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminAddCronjob add member failed ", "err", err.Error())
		return ErrAdminAddMemberFailed
	}
	return nil
}

func (c *service) AdminDelMember(ctx context.Context, am adminDoMemberRequest) error {
	group, err := c.repository.Groups().Find(am.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDelMember get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	member, err := c.repository.Member().FindById(am.MemberId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDelMember get member info by member_id has error ", "err", err.Error())
		return ErrGetUserInfoError
	}

	err = c.repository.Groups().GroupDelMember(group, member)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDelMember del member failed ", "err", err.Error())
		return ErrAdminDelMemberFailed
	}
	return nil
}

func (c *service) AdminDelProject(ctx context.Context, aq adminDoProjectRequest) error {
	group, err := c.repository.Groups().Find(aq.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDelProject get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	project, err := c.repository.Project().Find(aq.ProjectId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDelProject get project info by project_id has error ", "err", err.Error())
		return ErrGetProjectInfoFailed
	}

	// if the project namespace different from group namespace
	// return error
	if project.Namespace != group.Namespace {
		_ = level.Error(c.logger).Log("group", "AdminDelProject project namespace not the same of group namespace ", "err", " project namespace different from the group namespace ")
		return ErrPNsDiffGNs
	}

	err = c.repository.Groups().GroupDelProject(group, project)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDelProject del project failed ", "err", err.Error())
		return ErrAdminDelProjectFailed
	}
	return nil
}

func (c *service) AdminDelCronjob(ctx context.Context, ac adminDoCronjobRequest) error {
	group, err := c.repository.Groups().Find(ac.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDelCronjob get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	cronjob, err := c.repository.CronJob().Find(ac.CronjobId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDelCronjob get cronjob info by cronjob_id has error ", "err", err.Error())
		return ErrGetCronjobInfoFailed
	}

	// if the cronjob namespace different from group namespace
	// return error
	if cronjob.Namespace != group.Namespace {
		_ = level.Error(c.logger).Log("group", "AdminDelCronjob cronjob namespace not the same of group namespace ", "err", " cronjob namespace different from the group namespace ")
		return ErrCNsDiffGNs
	}

	err = c.repository.Groups().GroupDelCronjob(group, cronjob)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "AdminDelCronjob del cronjob failed ", "err", err.Error())
		return ErrAdminDelCronjobFailed
	}
	return nil
}

func (c *service) OwnerAddGroup(ctx context.Context, oq ownerDoGroup) error {
	_, exists := c.repository.Groups().GroupExistsByName(oq.Name)
	if !exists {
		_ = level.Error(c.logger).Log("group", "OwnerAddGroup get group info by name", "err", "Group name exists ")
		return ErrGroupNameExists
	}

	memberId := ctx.Value(middleware.UserIdContext).(int64)

	//判断 归属人 有没有这个ns权限
	nss := ctx.Value(middleware.NamespacesContext).([]string)

	var isTrue bool
	for _, v := range nss {
		if v == oq.Namespace {
			isTrue = true
		}
	}

	if !isTrue {
		_ = level.Error(c.logger).Log("group", "OwnerAddGroup", "err", "namespace is not allowed")
		return ErrNamespaceNotAllowed
	}

	memberInfo, err := c.repository.Member().GetInfoById(memberId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddGroup get user info by id", "err", err.Error())
		return ErrGetUserInfoError
	}

	if err := c.repository.Groups().CreateGroupAndRelation(&types.Groups{
		Name:        oq.Name,
		DisplayName: oq.DisplayName,
		Namespace:   oq.Namespace,
		MemberId:    memberId,
	}, memberInfo); err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddGroup", "err", err.Error())
		return ErrGroupCreateFailed
	}

	_ = level.Info(c.logger).Log("group", "OwnerAddGroup", "message", "succeed to create group")

	return nil
}

func (c *service) OwnerUpdateGroup(ctx context.Context, oq ownerDoGroup) error {
	group, err := c.repository.Groups().Find(oq.Id)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerUpdateGroup get group info by id", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	if group.Name != oq.Name {
		_, exists := c.repository.Groups().GroupNameExists(oq.Name, group.ID)
		if !exists {
			_ = level.Error(c.logger).Log("group", "OwnerUpdateGroup get group info by name", "err", "group name exists")
			return ErrGroupNameExists
		}
	}

	if group.Name != oq.Name {
		_, exists := c.repository.Groups().GroupNameExists(oq.Name, group.ID)
		if !exists {
			_ = level.Error(c.logger).Log("group", "OwnerUpdateGroup get group info by name_en", "err", "group name_en exists")
			return ErrGroupNameEnExists
		}
	}

	if err = c.repository.Groups().UpdateGroup(&types.Groups{
		Name:        oq.Name,
		DisplayName: oq.DisplayName,
		//Namespace:   oq.Namespace, // 普通组所有者不能修改组的业务线的,需要操作的去找管理员去管理员操作组的界面操作
	}, oq.Id); err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerUpdateGroup update group info error ", "err", err.Error())
		return ErrOwnerUpdateGroup
	}

	_ = level.Info(c.logger).Log("group", "OwnerUpdateGroup", "message", " owner succeed to update group")
	return nil
}

func (c *service) OwnerDelGroup(ctx context.Context, oq ownerDoGroup) error {
	group, err := c.repository.Groups().Find(oq.Id)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerUpdateGroup get group info by id", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	err = c.repository.Groups().DestroyAndRelation(group)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelGroup delete group and relation has error", "err", err.Error())
		return ErrOwnerDelGroupFailed
	}

	return nil
}

func (c *service) OwnerAddMember(ctx context.Context, om ownerDoMember) error {
	group, err := c.repository.Groups().Find(om.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddMember get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	member, err := c.repository.Member().FindById(om.MemberId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddMember get member info by member_id has error ", "err", err.Error())
		return ErrGetUserInfoError
	}

	var isTrue bool
	for _, v := range member.Namespaces {
		if v.Name == group.Namespace {
			isTrue = true
		}
	}
	// if the member namespace different from group namespace
	// return error
	if !isTrue {
		_ = level.Error(c.logger).Log("group", "OwnerAddMember member namespace not the same of group namespace ", "err", " member namespace different from the group namespace ")
		return ErrMNsDiffGNs
	}

	err = c.repository.Groups().GroupAddMember(group, member)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddMember add member failed ", "err", err.Error())
		return ErrOwnerAddMemberFailed
	}
	return nil
}

func (c *service) OwnerDelMember(ctx context.Context, om ownerDoMember) error {
	group, err := c.repository.Groups().Find(om.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelMember get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	userId := ctx.Value(middleware.UserIdContext).(int64)
	if userId == om.MemberId {
		_ = level.Error(c.logger).Log("group", "OwnerDelMember can not del self ", "err", "can not del self ")
		return ErrOwnerCanNotDelSelf
	}
	member, err := c.repository.Member().FindById(om.MemberId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelMember get member info by member_id has error ", "err", err.Error())
		return ErrGetUserInfoError
	}

	err = c.repository.Groups().GroupDelMember(group, member)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelMember del member failed ", "err", err.Error())
		return ErrOwnerDelMemberFailed
	}
	return nil
}

func (c *service) OwnerAddProject(ctx context.Context, op ownerDoProject) error {
	group, err := c.repository.Groups().Find(op.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddProject get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	project, err := c.repository.Project().Find(op.ProjectId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddProject get project info by project_id has error ", "err", err.Error())
		return ErrGetProjectInfoFailed
	}

	// if the project namespace different from group namespace
	// return error
	if project.Namespace != group.Namespace {
		_ = level.Error(c.logger).Log("group", "OwnerAddProject project namespace not the same of group namespace ", "err", " project namespace different from the group namespace ")
		return ErrPNsDiffGNs
	}

	err = c.repository.Groups().GroupAddProject(group, project)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddProject add project failed ", "err", err.Error())
		return ErrOwnerAddProjectFailed
	}
	return nil
}

func (c *service) OwnerDelProject(ctx context.Context, op ownerDoProject) error {
	group, err := c.repository.Groups().Find(op.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelProject get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	project, err := c.repository.Project().Find(op.ProjectId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelProject get project info by project_id has error ", "err", err.Error())
		return ErrGetProjectInfoFailed
	}

	err = c.repository.Groups().GroupDelProject(group, project)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelProject del project failed ", "err", err.Error())
		return ErrOwnerDelProjectFailed
	}
	return nil
}

func (c *service) OwnerAddCronjob(ctx context.Context, oc ownerDoCronjob) error {
	group, err := c.repository.Groups().Find(oc.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddCronjob get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	cronjob, err := c.repository.CronJob().Find(oc.CronjobId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddCronjob get cronjob info by cronjob_id has error ", "err", err.Error())
		return ErrGetCronjobInfoFailed
	}

	// if the cronjob namespace different from group namespace
	// return error
	if cronjob.Namespace != group.Namespace {
		_ = level.Error(c.logger).Log("group", "OwnerAddCronjob cronjob namespace not the same of group namespace ", "err", " cronjob namespace different from the group namespace ")
		return ErrCNsDiffGNs
	}

	err = c.repository.Groups().GroupAddCronjob(group, cronjob)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerAddCronjob add cronjob failed ", "err", err.Error())
		return ErrOwnerAddCronjobFailed
	}
	return nil
}

func (c *service) OwnerDelCronjob(ctx context.Context, oc ownerDoCronjob) error {
	group, err := c.repository.Groups().Find(oc.GroupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelCronjob get group info by group_id has error ", "err", err.Error())
		return ErrGetGroupInfoFailed
	}

	cronjob, err := c.repository.CronJob().Find(oc.CronjobId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelCronjob get cronjob info by cronjob_id has error ", "err", err.Error())
		return ErrGetCronjobInfoFailed
	}

	err = c.repository.Groups().GroupDelCronjob(group, cronjob)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "OwnerDelCronjob del cronjob failed ", "err", err.Error())
		return ErrOwnerDelCronjobFailed
	}
	return nil
}

func (c *service) UserMyList(ctx context.Context, umq userGroupListRequest) ([]*types.Groups, error) {
	isAdmin := ctx.Value(middleware.IsAdmin).(bool)
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	res, err := c.repository.Groups().UserMyGroupList(umq.Name, umq.Ns, memberId, isAdmin)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "UserMyList get my group list failed ", "err", err.Error())
		return nil, ErrGetMyGroupListFailed
	}
	return res, nil
}

func (c *service) NsMyList(ctx context.Context) ([]types.Namespace, error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	res, err := c.repository.Namespace().UserMyNsList(memberId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "NsMyList get my group list failed ", "err", err.Error())
		return nil, ErrGetMyNsListFailed
	}
	return res, nil
}

func (c *service) RelDetail(ctx context.Context, groupId int64) (*types.Groups, error) {
	g, err := c.repository.Groups().RelDetail(groupId)
	if err != nil {
		_ = level.Error(c.logger).Log("group", "RelDetail get  group detail relation failed ", "err", err.Error())
		return nil, ErrGetGroupRelListFailed
	}

	for k, v := range g.Members {
		g.Members[k] = types.Member{
			ID:        v.ID,
			Email:     v.Email,
			Username:  v.Username,
			State:     v.State,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
	}

	return g, nil
}

func (c *service) GroupNameExists(ctx context.Context, name string) (bool, error) {
	notFound := c.repository.Groups().GroupNameIsExists(name)
	return notFound, nil
}

func (c *service) GroupDisplayNameExists(ctx context.Context, displayName string) (bool, error) {
	notFound := c.repository.Groups().GroupDisplayNameIsExists(displayName)
	return notFound, nil
}

func NewService(logger log.Logger, config *config.Config,
	repository repository.Repository) Service {
	return &service{logger, config,
		repository}
}
