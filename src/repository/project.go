/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-06-25
 * Time: 21:10
 */
package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type AuditState int
type PublishState int

const (
	AuditUnsubmit AuditState = iota // 未提交
	AuditSubmit                     // 提交
	AuditFail                       // 审核失败
	AuditPass                       // 审核成功

)

const (
	PublishUnSubmit PublishState = iota // 未发布
	PublishPass                         // 发布成功
	PublishFail                         // 发布失败
)

type Language string

const (
	Java   Language = "Java"
	Golang Language = "Golang"
	NodeJS Language = "NodeJS"
	Python Language = "Python"
)

func (c Language) String() string {
	return string(c)
}

type ProjectRepository interface {
	Create(p *types.Project) error
	Find(pId int64) (*types.Project, error)
	FindByNsName(ns, name string) (res *types.Project, err error)
	FindByNsNameOnly(ns, name string) (res *types.Project, err error)
	FindByNsNameExist(ns, name string) (res *types.Project, notExist bool)
	GetProjectByNs(ns string) (res []*types.Project, err error)
	GetProjectAndTemplateByNs(ns string, name string, offset, limit int) (res []*types.Project, count int64, err error)
	GetProjectByNameLikeAndNs(nameLike string, ns string) (res []*types.Project, err error)
	GetProjectByMid(mid int64) (res []*types.Project, err error)
	GetProjectByGroupId(gid int64) (res []*types.Project, err error)
	GetProjectByMidAndNs(mid int64, ns string) (res []*types.Project, err error)
	Update(project *types.Project) error
	CountLanguage() (languageCount []*LanguageCount, err error)
	UpdateProjectById(project *types.Project) error
	Delete(ns, name string) error
	//Count(name string)
	GetProjectByGroupAndPNsAndPName(pName, pNs string, groupId int64, offset, limit int) (projects []*types.Project, count int64, err error)
	GetProjectByNsLimit(ns string) (res []*types.Project, err error)
}

type LanguageCount struct {
	Language string `json:"language"`
	Total    int64  `json:"number"`
}

type project struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &project{db: db}
}

func (c *project) CountLanguage() (languageCount []*LanguageCount, err error) {
	err = c.db.Raw("SELECT coalesce(`language`,'Total') AS language ,COUNT(`id`) AS total FROM projects WHERE `audit_state` = ? GROUP BY `language` WITH ROLLUP", 3).Scan(&languageCount).Error
	return
}

func (c *project) Create(p *types.Project) error {
	return c.db.Save(p).Error
}

func (c *project) Update(project *types.Project) error {
	return c.db.Model(&types.Project{}).Where("id = ?", project.ID).
		Update(types.Project{DisplayName: project.DisplayName, Desc: project.Desc}).Error
}

func (c *project) UpdateProjectById(project *types.Project) error {
	return c.db.Model(&types.Project{}).Where("id = ?", project.ID).Update(project).Error
}

func (c *project) Find(pId int64) (*types.Project, error) {
	var pro types.Project
	err := c.db.First(&pro, pId).Error
	return &pro, err
}

func (c *project) FindByNsName(ns, name string) (res *types.Project, err error) {
	var rs types.Project
	err = c.db.Where("name = ? AND namespace = ? AND audit_state = ?", name, ns, AuditPass).
		Preload("Member").First(&rs).Error
	return &rs, err
}

func (c *project) FindByNsNameOnly(ns, name string) (res *types.Project, err error) {
	var rs types.Project
	err = c.db.Where("name = ? AND namespace = ? ", name, ns).
		Preload("Member").First(&rs).Error
	return &rs, err
}

func (c *project) FindByNsNameExist(ns, name string) (res *types.Project, notExist bool) {
	var rs types.Project
	notExist = c.db.Where("name = ? AND namespace = ?", name, ns).First(&rs).RecordNotFound()
	return &rs, notExist
}

func (c *project) GetProjectByNs(ns string) (res []*types.Project, err error) {
	err = c.db.
		Where("namespace = ?", ns).
		Where("publish_state = ?", PublishPass).
		Where("audit_state = ?", AuditPass).
		Preload("Member").
		Find(&res).Error
	return
}

func (c *project) GetProjectByNameLikeAndNs(nameLike string, ns string) (res []*types.Project, err error) {
	err = c.db.Where("namespace = ?", ns).
		Where("publish_state = ?", PublishPass).
		Where("audit_state = ?", AuditPass).
		Where("name like ?", "%"+nameLike+"%").
		Preload("Member").
		Find(&res).Error
	return
}

func (c *project) GetProjectByMid(mid int64) (res []*types.Project, err error) {
	err = c.db.
		Where("member_id = ?", mid).
		Preload("Member").
		Find(&res).Error
	return
}

func (c *project) GetProjectByGroupId(gid int64) (res []*types.Project, err error) {
	err = c.db.
		Where("group_id = ?", gid).
		Preload("Groups").
		Find(&res).Error
	return
}

func (c *project) GetProjectByMidAndNs(mid int64, ns string) (res []*types.Project, err error) {
	var gms []*types.GroupsMemberss
	c.db.Table("groups_memberss").Where("members_id = ?", mid).Find(&gms)

	if len(gms) == 0 {
		err = c.db.Where("namespace = ? and audit_state = ? and member_id = ?", ns, 3, mid).Find(&res).Error
		return
	}

	//组内所有项目
	var gid []int64
	for _, v := range gms {
		gid = append(gid, v.ID)
	}

	query := c.db.Table("projects")
	query = query.Select("projects.*")
	query = query.Joins("left join groups_projectss ON groups_projectss.projects_id = projects.id")
	query = query.Joins("left join groups ON groups.id = groups_projectss.groups_id")
	query = query.Where("projects.namespace = ?", ns)
	query = query.Where("projects.audit_state = ?", 3)
	query = query.Where("projects.member_id = ?", mid)
	query = query.Where("groups.id in (?)", gid)
	err = query.Order("id desc").Scan(&res).Error

	return
}

func (c *project) GetProjectAndTemplateByNs(ns string, name string, offset, limit int) (res []*types.Project, count int64, err error) {
	query := c.db.Order(gorm.Expr("id DESC")).Where("namespace = ?", ns)
	if name != "" {
		query = query.Where("name like ? ", "%"+name+"%")
	}

	err = query.Preload("ProjectTemplates").
		Offset(offset).Limit(limit).
		Preload("Member").Find(&res).Error

	queryCount := c.db.Model(&types.Project{}).Where("namespace = ? ", ns)
	if name != "" {
		queryCount = queryCount.Where("name like ?", "%"+name+"%")
	}
	err = queryCount.Count(&count).Error

	return
}

func (c *project) Delete(ns, name string) error {
	return c.db.Where("namespace = ? AND name = ?", ns, name).Delete(types.Project{}).Error
}

func (c *project) GetProjectByGroupAndPNsAndPName(pName, pNs string, groupId int64, offset, limit int) (projects []*types.Project, count int64, err error) {
	var group types.Groups
	err = c.db.Where("id = ?", groupId).Preload("Projects").Find(&group).Error
	if err != nil {
		return
	}
	var projectIds []int64
	for _, v := range group.Projects {
		projectIds = append(projectIds, v.ID)
	}

	//err = c.db.Where("name = ?", pName).Where("namespace = ?", pNs).Where("id in (?)", projectIds).Offset(offset).Limit(limit).Order("id asc").Find(&projects).Error
	//if err != nil {
	//	return nil, err
	//}

	query := c.db.Order(gorm.Expr("id DESC")).Where("namespace = ?", pNs).Where("id in (?)", projectIds)
	if pName != "" {
		query = query.Where("name like ? ", "%"+pName+"%")
	}

	err = query.Preload("ProjectTemplates").
		Offset(offset).Limit(limit).
		Preload("Member").Find(&projects).Error

	queryCount := c.db.Model(&types.Project{}).Where("namespace = ? ", pNs).Where("id in (?)", projectIds)
	if pName != "" {
		queryCount = queryCount.Where("name like ?", "%"+pName+"%")
	}
	err = queryCount.Count(&count).Error

	return
}

func (c *project) GetProjectByNsLimit(ns string) (res []*types.Project, err error) {
	err = c.db.
		Where("namespace = ?", ns).
		Where("publish_state = ?", PublishPass).
		Where("audit_state = ?", AuditPass).
		Preload("Member").
		Limit(6).
		Order("id desc").
		Find(&res).Error
	return
}
