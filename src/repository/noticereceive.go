/**
 * @Time : 2019-07-08 11:37
 * @Author : soupzhb@gmail.com
 * @File : noticereceive.go
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type NoticeReceiveRepository interface {
	Create(nr *types.NoticeReceive) error
	Update(nr *types.NoticeReceive) error
	FindListByMid(memberId int64) (list []*types.NoticeReceive, err error)
	GetNoticeReceiveByMidAction(memberId int64, action string) (nr *types.NoticeReceive, err error)
	GetNoticeReceiveByAction(action string) (nr []*types.NoticeReceive, err error)
}

type noticeReceive struct {
	db *gorm.DB
}

func NewNoticeReceiveRepository(db *gorm.DB) NoticeReceiveRepository {
	return &noticeReceive{db: db}
}

func (c *noticeReceive) Create(nr *types.NoticeReceive) error {
	return c.db.Save(nr).Error
}

func (c *noticeReceive) Update(nr *types.NoticeReceive) error {
	return c.db.Save(nr).Error
}

func (c *noticeReceive) FindListByMid(memberId int64) (list []*types.NoticeReceive, err error) {
	err = c.db.Where("member_id = ?", memberId).Find(&list).Error
	return
}

func (c *noticeReceive) GetNoticeReceiveByMidAction(memberId int64, action string) (nr *types.NoticeReceive, err error) {
	var temp types.NoticeReceive
	if err = c.db.Where("member_id = ? and notice_action = ?", memberId, action).First(&temp).Error; err != nil {
		return
	}
	return &temp, nil
}

func (c *noticeReceive) GetNoticeReceiveByAction(action string) (list []*types.NoticeReceive, err error) {
	if err = c.db.Where("notice_action = ?", action).Find(&list).Error; err != nil {
		return
	}
	return
}
