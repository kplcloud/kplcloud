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

type CronjobRepository interface {
	Create(*types.Cronjob) (*types.Cronjob, error)
	Find(cId int64) (*types.Cronjob, error)
	Update(cj *types.Cronjob, id int64) error
	Delete(id int64) error
	GetCronjobByNs(ns string) ([]*types.Cronjob, error)
	GetCronjobByNameLikeAndNs(nameLike string, ns string) ([]*types.Cronjob, error)
	GetCronJobByNameAndNs(name string, ns string) (*types.Cronjob, bool)
	CronJobCountWithGroup(name string, ns string, group int64) (int64, error)
	CronJobPaginateWithGroup(name string, ns string, group int64, offset int, limit int) ([]*types.Cronjob, error)
}

type cronjob struct {
	db *gorm.DB
}

func NewCronjobRepository(db *gorm.DB) CronjobRepository {
	return &cronjob{db: db}
}

func (c *cronjob) Find(cId int64) (*types.Cronjob, error) {
	var cron types.Cronjob
	err := c.db.First(&cron, cId).Error
	return &cron, err
}

func (c *cronjob) GetCronjobByNs(ns string) (res []*types.Cronjob, err error) {
	err = c.db.Where("namespace = ?", ns).Find(&res).Error
	return
}

func (c *cronjob) GetCronjobByNameLikeAndNs(nameLike string, ns string) (res []*types.Cronjob, err error) {
	err = c.db.Where("namespace = ?", ns).
		Where("name like ?", "%"+nameLike+"%").
		Find(&res).Error
	return
}

func (c *cronjob) GetCronJobByNameAndNs(name string, ns string) (*types.Cronjob, bool) {
	var cronjob types.Cronjob
	isExists := c.db.Where("name = ?", name).Where("namespace = ?", ns).First(&cronjob).RecordNotFound()
	return &cronjob, isExists
}

func (c *cronjob) Create(cj *types.Cronjob) (*types.Cronjob, error) {
	err := c.db.Create(&cj).Error
	return cj, err
}

func (c *cronjob) Update(cj *types.Cronjob, id int64) error {
	return c.db.Where("id = ?", id).Save(&cj).Error
}

func (c *cronjob) Delete(id int64) (err error) {
	return c.db.Delete(&types.Cronjob{ID: id}).Error
}

func (c *cronjob) CronJobCountWithGroup(name string, ns string, group int64) (cnt int64, err error) {
	var cjs []*types.Cronjob
	var cj types.Cronjob
	query := c.db

	if group > 0 {
		var gc types.GroupsCronjobs
		var g types.Groups
		query = query.Joins("inner join "+gc.TableName()+" T1 ON T1.`cronjobs_id` = "+cj.TableName()+".`id` "+
			" inner join "+g.TableName()+" T2 on T2.`id` = T1.`groups_id` ").Where("T2.id = ?", group)
	}

	if name != "" {
		query = query.Where(cj.TableName()+".name = ?", name)
	}

	if ns != "" {
		query = query.Where(cj.TableName()+".namespace = ?", ns)
	}

	err = query.Find(&cjs).Count(&cnt).Error
	return cnt, err
}

func (c *cronjob) CronJobPaginateWithGroup(name string, ns string, group int64, offset int, limit int) (cronjobs []*types.Cronjob, err error) {
	var cjs []*types.Cronjob
	var cj types.Cronjob
	query := c.db

	//SELECT T0.`id`, T0.`name`, T0.`namespace`, T0.`schedule`, T0.`git_type`, T0.`git_path`, T0.`image`, T0.`suspend`, T0.`active`, T0.`last_schedule`, T0.`conf_map_name`, T0.`args`, T0.`log_path`, T0.`add_type`, T0.`created_at`, T0.`updated_at`, T0.`member_id`
	// FROM `cronjobs` T0
	// INNER JOIN `groups_cronjobss` T1
	// ON T1.`cronjobs_id` = T0.`id`
	// INNER JOIN `groups` T2
	// ON
	// T2.`id` = T1.`groups_id`
	// WHERE
	// T0.`namespace` = ?
	// AND
	// T2.`id` = ?
	// ORDER BY T0.`id` DESC
	// LIMIT 10

	if group > 0 {
		var gc types.GroupsCronjobs
		var g types.Groups
		query = query.Joins("inner join "+gc.TableName()+" T1 ON T1.`cronjobs_id` = "+cj.TableName()+".`id` "+
			" inner join "+g.TableName()+" T2 on T2.`id` = T1.`groups_id` ").Where("T2.id = ?", group)
	}

	if name != "" {
		query = query.Where(cj.TableName()+".name = ?", name)
	}

	if ns != "" {
		query = query.Where(cj.TableName()+".namespace = ?", ns)
	}

	err = query.Offset(offset).Limit(limit).Order(cj.TableName() + ".id").Find(&cjs).Error

	return cjs, err
}
