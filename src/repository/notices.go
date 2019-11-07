package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type NoticesRepository interface {
	FindById(id int) (v *types.Notices, err error)
	Create(n *types.Notices) error
	CreateReturnId(n *types.Notices) int64
	Count(name string, noticeType int) (count int, err error)
	FindOffsetLimit(name string, noticeType int, offset, limit int) (res []*types.Notices, err error)
	CountByAction(ns, name string, action types.NoticeAction) (total int64, err error)
}

type notices struct {
	db *gorm.DB
}

func NewNoticesRepository(db *gorm.DB) NoticesRepository {
	return &notices{db: db}
}

func (c *notices) List() {

}

/**
 * @Title 创建消息
 */
func (c *notices) Create(ns *types.Notices) error {
	return c.db.Save(ns).Error
}

/**
 * @Title 创建消息并返回自增id
 */
func (c *notices) CreateReturnId(ns *types.Notices) int64 {
	return c.db.Save(ns).RowsAffected
}

/**
 * @Title 获取消息详情
 */
func (c *notices) FindById(id int) (v *types.Notices, err error) {
	var n types.Notices
	if err = c.db.First(&n, id).Error; err != nil {
		return
	}
	return &n, nil
}

/**
 * @Title 获取总数
 */
func (c *notices) Count(name string, noticeType int) (count int, err error) {
	var n types.Notices
	query := c.db.Model(&n)
	if name != "" {
		query = query.Where("title like ?", "%"+name+"%")
	}
	if noticeType > 0 {
		query = query.Where("type = ?", noticeType)
	}
	err = query.Count(&count).Error
	return
}

func (c *notices) CountByAction(ns, name string, action types.NoticeAction) (total int64, err error) {
	err = c.db.Model(&types.Notices{}).
		Where("name = ?", name).
		Where("namespace = ?", ns).
		Where("action = ?", string(action)).Count(&total).Error
	return
}

/**
* @Title 获取列表
 */
func (c *notices) FindOffsetLimit(name string, noticeType int, offset, limit int) (res []*types.Notices, err error) {
	var list []*types.Notices
	query := c.db
	if name != "" {
		query = query.Where("title like ?", "%"+name+"%")
	}
	if noticeType > 0 {
		query = query.Where("type = ?", noticeType)
	}
	err = query.Offset(offset).Order("id desc").Limit(limit).Find(&list).Error
	return list, err
}
