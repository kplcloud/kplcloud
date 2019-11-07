/**
 * @Time : 2019-07-09 15:55
 * @Author : solacowa@gmail.com
 * @File : build
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

const (
	Building      = "BUILDING"
	BuildSuccess  = "SUCCESS"
	BuildFailure  = "FAILURE"
	BuildAborted  = "ABORTED"
	BuildRoolback = "ROLLBACK"
	NoBuild       = "NOBUILD"
)

type StatisticsRequest struct {
	Namespace   string
	Name        string
	ProjectName []string
	STime       string
	ETime       string
	BuildID     int
	GitType     string
}

type Ress struct {
	types.Build
	Count int
}

type BuildRepository interface {
	FirstByTag(ns, name, version string) (res *types.Build, err error)
	FindBuildByBuildId(ns, name string, buildId int) (res types.Build, err error)
	FindById(ns, name string, id int64) (res types.Build, err error)
	CreateBuild(rep *types.Build) (*types.Build, error)
	Update(b *types.Build) error
	FindOffsetLimit(ns, name string, offset, limit int) (res []*types.Build, err error)
	Count(ns, name string) (count int64, err error)
	FindStatisticsOffsetLimit(req StatisticsRequest, offset, limit int) (res []*types.Build, err error)
	CountStatistics(req StatisticsRequest) (count int64, err error)
	GetGroupByBuilds(req StatisticsRequest, groupBy string) (ress []Ress, err error)
	Delete(ns, name string) error
	CountByStatus(ns, name, buildStatus string) (total int64, err error)
}

type build struct {
	db *gorm.DB
}

func NewBuildRepository(db *gorm.DB) BuildRepository {
	return &build{db: db}
}

func (c *build) CountByStatus(ns, name, buildStatus string) (total int64, err error) {
	query := c.db.Model(&types.Build{}).Where("namespace = ? AND name = ?", ns, name)
	if buildStatus != "" {
		query = query.Where("status = ? ", buildStatus)
	}
	err = query.Count(&total).Error
	return
}

func (c *build) Count(ns, name string) (count int64, err error) {
	err = c.db.Model(&types.Build{}).
		Where("namespace = ? AND name = ?", ns, name).Count(&count).Error
	return
}

func (c *build) FindById(ns, name string, id int64) (res types.Build, err error) {
	err = c.db.First(&res, "namespace = ? AND name = ? AND id = ?", ns, name, id).Error
	return
}

func (c *build) FindOffsetLimit(ns, name string, offset, limit int) (res []*types.Build, err error) {
	query := c.db.Where("namespace = ?", ns)
	if name != "" {
		query = query.Where("name = ?", name)
	}
	err = query.Preload("Member", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,username,email")
	}).
		Order(gorm.Expr("id DESC")).
		Offset(offset).Limit(limit).Find(&res).Error
	return
}

func (c *build) FirstByTag(ns, name, version string) (res *types.Build, err error) {
	var resp types.Build
	if err = c.db.First(&resp, "namespace = ? AND name = ? AND version = ? AND status = ? ", ns, name, version, "BUILDING").Error; err != nil {
		return
	}

	return &resp, nil
}

func (c *build) FindBuildByBuildId(ns, name string, buildId int) (res types.Build, err error) {
	err = c.db.First(&res, "namespace = ? AND name = ? AND build_id = ?", ns, name, buildId).Error
	return
}

func (c *build) CreateBuild(rep *types.Build) (*types.Build, error) {
	if err := c.db.Save(rep).Error; err != nil {
		return nil, err
	}
	return rep, nil
}

func (c *build) Update(b *types.Build) error {
	return c.db.Model(b).Where("id = ?", b.ID).Update(b).Error
}

func (c *build) FindStatisticsOffsetLimit(req StatisticsRequest, offset, limit int) (res []*types.Build, err error) {
	query := c.db.Table("builds")
	if req.Namespace != "" {
		query = query.Where("namespace = ?", req.Namespace)
	}
	if req.Name != "" {
		query = query.Where("name like ?", "%"+req.Name+"%")
	}
	if len(req.ProjectName) > 0 {
		query = query.Where("name in (?)", req.ProjectName)
	}
	if req.BuildID != 0 {
		query = query.Where("build_id = ?", req.BuildID)
	}
	if req.GitType != "" {
		query = query.Where("git_type = ?", req.GitType)
	}
	err = query.Group("namespace,name").Order("id desc").Offset(offset).Limit(limit).Find(&res).Error
	return
}

func (c *build) CountStatistics(req StatisticsRequest) (count int64, err error) {
	query := c.db.Table("builds")
	if req.Namespace != "" {
		query = query.Where("namespace = ?", req.Namespace)
	}
	if req.Name != "" {
		query = query.Where("name like ?", "%"+req.Name+"%")
	}
	if len(req.ProjectName) > 0 {
		query = query.Where("name in (?)", req.ProjectName)
	}
	if req.BuildID != 0 {
		query = query.Where("build_id = ?", req.BuildID)
	}
	if req.GitType != "" {
		query = query.Where("git_type = ?", req.GitType)
	}
	if req.STime != "" {
		query = query.Where("created_at >= ?", req.STime)
	}
	if req.ETime != "" {
		query = query.Where("created_at <= ?", req.ETime)
	}
	err = query.Group("namespace,name").Count(&count).Error
	return
}

func (c *build) GetGroupByBuilds(req StatisticsRequest, groupBy string) (ress []Ress, err error) {
	query := c.db.Table("builds").Select("*, count(*) as count")
	if req.Namespace != "" {
		query = query.Where("namespace = ?", req.Namespace)
	}
	if req.Name != "" {
		query = query.Where("name like ?", "%"+req.Name+"%")
	}
	if len(req.ProjectName) > 0 {
		query = query.Where("name in (?)", req.ProjectName)
	}
	if req.BuildID != 0 {
		query = query.Where("build_id = ?", req.BuildID)
	}
	if req.GitType != "" {
		query = query.Where("git_type = ?", req.GitType)
	}
	if req.STime != "" {
		query = query.Where("created_at >= ?", req.STime)
	}
	if req.ETime != "" {
		query = query.Where("created_at <= ?", req.ETime)
	}
	if groupBy != "" {
		query = query.Group(groupBy)
	}
	err = query.Find(&ress).Error
	return
}

func (c *build) Delete(ns, name string) error {
	return c.db.Where("name = ? AND namespace = ?", name, ns).Delete(&types.Build{}).Error
}
