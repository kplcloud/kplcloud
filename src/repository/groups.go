/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-06-24
 * Time: 17:43
 */
package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type GroupsRepository interface {
	Create(g *types.Groups) error
	Find(id int64) (*types.Groups, error)
	GroupExistsById(id int64) (*types.Groups, bool)
	CreateGroupAndRelation(g *types.Groups, member *types.Member) error
	GetGroupByName(name string) (*types.Groups, error)
	GroupExistsByName(name string) (*types.Groups, bool)
	GetGroupByDisplayName(displayName string) (*types.Groups, error)
	GroupExistsByDisplayName(displayName string) (*types.Groups, bool)
	GroupDisplayNameExists(displayName string, id int64) (*types.Groups, bool)
	GroupNameExists(name string, id int64) (*types.Groups, bool)
	AllGroupsCount(displayName string, ns string) (int64, error)
	GroupsPaginate(displayName string, ns string, offset int, limit int) ([]types.Groups, error)
	UpdateGroup(g *types.Groups, groupId int64) error
	UpdateGroupAndRelation(g *types.Groups, groupId int64, member *types.Member) error
	DestroyAndRelation(group *types.Groups) error
	GroupAddProject(group *types.Groups, project *types.Project) error
	GroupDelProject(group *types.Groups, project *types.Project) error
	GroupAddCronjob(group *types.Groups, cronjob *types.Cronjob) error
	GroupDelCronjob(group *types.Groups, cronjob *types.Cronjob) error
	GroupAddMember(group *types.Groups, member *types.Member) error
	GroupDelMember(group *types.Groups, member *types.Member) error
	UserMyGroupList(name string, ns string, memberId int64, isAdmin bool) ([]*types.Groups, error)
	RelDetail(groupId int64) (*types.Groups, error)
	IsInGroup(groupId int64, memberId int64) (bool, error)
	CheckPermissionForMidCronJob(cronJobId int64, groupIds []int64) (notFound bool, err error)
	CheckPermissionForMidProject(projectId int64, groupIds []int64) (notFound bool, err error)
	GetMemberIdsByProjectNameAndNs(pName, pNs string) (members []types.Member)
	GetIndexProjectByMemberIdAndNs(memberId int64, ns string) (projectLists []*types.Project, err error)
	GroupNameIsExists(name string) (notFound bool)
	GroupDisplayNameIsExists(displayName string) (notFound bool)
}

type group struct {
	db *gorm.DB
}

func NewGroupsRepository(db *gorm.DB) GroupsRepository {
	return &group{db: db}
}

func (c *group) Create(g *types.Groups) error {
	return c.db.Create(g).Error
}

func (c *group) GroupExistsById(id int64) (*types.Groups, bool) {
	var g types.Groups
	res := c.db.First(&g, id).RecordNotFound()
	return &g, res
}

func (c *group) Find(id int64) (g *types.Groups, err error) {
	var gs types.Groups
	err = c.db.First(&gs, id).Error
	return &gs, err
}

func (c *group) CreateGroupAndRelation(group *types.Groups, member *types.Member) error {
	tx := c.db.Begin()

	err := tx.Create(group).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&group).Association("Members").Append(types.Member{ID: member.ID}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *group) GetGroupByName(name string) (res *types.Groups, err error) {
	var g types.Groups
	err = c.db.Where("name = ?", name).Attrs(types.Groups{ID: 0}).FirstOrInit(&g).Error
	return &g, err
}

func (c *group) GroupExistsByName(name string) (*types.Groups, bool) {
	var g types.Groups
	res := c.db.Where("name = ?", name).First(&g).RecordNotFound()
	return &g, res
}

func (c *group) GroupExistsByDisplayName(displayName string) (*types.Groups, bool) {
	var g types.Groups
	res := c.db.Where("display_name = ?", displayName).First(&g).RecordNotFound()
	return &g, res
}

func (c *group) GetGroupByDisplayName(displayName string) (res *types.Groups, err error) {
	g := new(types.Groups)
	err = c.db.Where("display_name = ?", displayName).Attrs(types.Groups{ID: 0}).FirstOrInit(&g).Error
	return g, err
}

func (c *group) GroupDisplayNameExists(name string, id int64) (*types.Groups, bool) {
	var g types.Groups
	res := c.db.Where("display_name = ?", name).Where("id != ?", id).First(&g).RecordNotFound()
	return &g, res
}

func (c *group) GroupNameExists(name string, id int64) (*types.Groups, bool) {
	var g types.Groups
	res := c.db.Where("name = ?", name).Where("id != ?", id).First(&g).RecordNotFound()
	return &g, res
}

func (c *group) AllGroupsCount(displayName string, ns string) (cnt int64, err error) {
	g := new(types.Groups)
	query := c.db.Model(&g)
	if displayName != "" {
		query = query.Where("display_name like ?", "%"+displayName+"%")
	}

	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}
	err = query.Count(&cnt).Error
	return cnt, err
}

func (c *group) GroupsPaginate(displayName string, ns string, offset int, limit int) ([]types.Groups, error) {
	var groups []types.Groups
	query := c.db
	if displayName != "" {
		query = query.Where("display_name like ?", "%"+displayName+"%")
	}

	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}

	err := query.Preload("Ns").Preload("Member").Offset(offset).Limit(limit).Find(&groups).Error
	return groups, err
}

func (c *group) UpdateGroup(g *types.Groups, groupId int64) error {
	err := c.db.Model(&g).Where("id = ?", groupId).Update(g).Error
	return err
}

func (c *group) UpdateGroupAndRelation(g *types.Groups, groupId int64, member *types.Member) error {
	tx := c.db.Begin()
	err := c.db.Model(&g).Where("id = ?", groupId).Update(g).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&g).Association("Members").Append(types.Member{ID: member.ID}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *group) DestroyAndRelation(g *types.Groups) error {
	tx := c.db.Begin()

	if g.ID < 1 {
		return errors.New("Group can not null ")
	}

	err := tx.Model(&g).Association("Members").Clear().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&g).Association("Cronjobs").Clear().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&g).Association("Projects").Clear().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&g).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *group) GroupAddProject(group *types.Groups, project *types.Project) error {
	return c.db.Model(&group).Association("Projects").Append(types.Project{ID: project.ID}).Error
}

func (c *group) GroupDelProject(group *types.Groups, project *types.Project) error {
	return c.db.Model(&group).Association("Projects").Delete(types.Project{ID: project.ID}).Error
}

func (c *group) GroupAddCronjob(group *types.Groups, cronjob *types.Cronjob) error {
	return c.db.Model(&group).Association("Cronjobs").Append(types.Cronjob{ID: cronjob.ID}).Error
}

func (c *group) GroupDelCronjob(group *types.Groups, cronjob *types.Cronjob) error {
	return c.db.Model(&group).Association("Cronjobs").Delete(types.Cronjob{ID: cronjob.ID}).Error
}

func (c *group) GroupAddMember(group *types.Groups, member *types.Member) error {
	return c.db.Model(&group).Association("Members").Append(types.Member{ID: member.ID}).Error
}

func (c *group) GroupDelMember(group *types.Groups, member *types.Member) error {
	return c.db.Model(&group).Association("Members").Delete(types.Member{ID: member.ID}).Error
}

func (c *group) UserMyGroupList(name string, ns string, memberId int64, isAdmin bool) (gs []*types.Groups, err error) {

	// group name = "name" && group namespace = "namespace" &&  member_id = memberId

	query := c.db
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}

	if ns != "" {
		query = query.Where("namespace = ?", ns)
	}

	if !isAdmin {
		var m types.Member
		var groupIds []int64

		err := c.db.Preload("Groups").Find(&m, memberId).Error

		if err != nil {
			return nil, err
		}

		for _, v := range m.Groups {
			groupIds = append(groupIds, v.ID)
		}

		err = query.
			Preload("Member", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,email,username,state,created_at")
			}).
			Preload("Ns").
			Find(&gs, "id in (?)", groupIds).Error
	} else {
		err = query.
			Preload("Member", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,email,username,state,created_at")
			}).
			Preload("Ns").
			Find(&gs).Error
	}

	return gs, err
}

func (c *group) RelDetail(groupId int64) (*types.Groups, error) {
	var g types.Groups
	err := c.db.
		Preload("Cronjobs").
		Preload("Members").
		Preload("Projects").First(&g, groupId).Error
	return &g, err
}

func (c *group) IsInGroup(groupId int64, memberId int64) (bool, error) {
	var g types.Groups
	err := c.db.Preload("Members").First(&g, groupId).Error
	var isTrue bool
	for _, v := range g.Members {
		if v.ID == memberId {
			isTrue = true
			break
		}
	}
	return isTrue, err
}

func (c *group) CheckPermissionForMidCronJob(cronJobId int64, groupIds []int64) (notFound bool, err error) {
	var cronjob types.Cronjob
	var group types.Groups
	groupsName := group.TableName()
	err = c.db.Where("id = ?", cronJobId).Preload("Groups", groupsName+".id in (?)", groupIds).Find(&cronjob).Error
	if len(cronjob.Groups) > 0 {
		return false, err
	}
	return true, err
}

func (c *group) CheckPermissionForMidProject(projectId int64, groupIds []int64) (notFound bool, err error) {
	var project types.Project
	var group types.Groups
	groupsName := group.TableName()
	err = c.db.Where("id = ?", projectId).Preload("Groups", groupsName+".id in (?)", groupIds).Find(&project).Error
	if len(project.Groups) > 0 {
		return false, err
	}
	return true, err
}

func (c *group) GetMemberIdsByProjectNameAndNs(pName, pNs string) (members []types.Member) {
	var project types.Project
	var group types.Groups
	notFound := c.db.Where("name = ?", pName).Where("namespace = ?", pNs).Preload("Groups").Preload("Member").Find(&project).RecordNotFound()
	if project.ID == 0 {
		return members
	}
	if notFound {
		members = append(members, project.Member)
		return members
	}
	var groupsIds []int64

	for _, v := range project.Groups {
		groupsIds = append(groupsIds, v.ID)
	}

	notFound = c.db.Where("id in (?)", groupsIds).Preload("Members").Find(&group).RecordNotFound()
	if notFound {
		members = append(members, project.Member)
		return members
	}
	members = group.Members
	var isTrue bool
	for _, v := range group.Members {
		if v.ID == project.MemberID {
			isTrue = true
			break
		}
	}
	if !isTrue {
		members = append(members, project.Member)
	}
	return members
}

func (c *group) GetIndexProjectByMemberIdAndNs(memberId int64, ns string) (projectLists []*types.Project, err error) {
	var member types.Member
	var projects []*types.Project
	query := c.db.Where("member_id = ?", memberId).
		Where("namespace = ?", ns).
		Where("publish_state = ?", PublishPass).
		Where("audit_state = ?", AuditPass).
		Order("id desc")

	notFound := c.db.Where("id = ?", memberId).Preload("Groups").First(&member).RecordNotFound()
	if notFound || len(member.Groups) == 0 {
		err = query.Limit(6).Find(&projects).Error
		return projects, err
	}

	var groupIds []int64
	for _, v := range member.Groups {
		groupIds = append(groupIds, v.ID)
	}

	var group types.Groups
	notFound = c.db.Where("id in (?)", groupIds).Preload("Projects").Find(&group).RecordNotFound()
	if notFound || len(group.Projects) == 0 {
		err = query.Limit(6).Find(&projects).Error
		return projects, err
	}

	var projectIds []int64
	for _, v := range group.Projects {
		projectIds = append(projectIds, v.ID)
	}

	err = query.Find(&projects).Error
	if err != nil {
		return
	}

	for _, n := range projects {
		projectIds = append(projectIds, n.ID)
	}

	err = c.db.Where("id in (?)", projectIds).
		Where("publish_state = ?", PublishPass).
		Where("audit_state = ?", AuditPass).
		Order("id desc").
		Limit(6).
		Find(&projectLists).Error
	return
}

func (c *group) GroupNameIsExists(name string) (notFound bool) {
	var group types.Groups
	notFound = c.db.Where("name = ?", name).First(&group).RecordNotFound()
	return
}

func (c *group) GroupDisplayNameIsExists(displayName string) (notFound bool) {
	var group types.Groups
	notFound = c.db.Where("display_name = ?", displayName).First(&group).RecordNotFound()
	return
}
