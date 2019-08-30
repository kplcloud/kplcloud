package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"gopkg.in/guregu/null.v3"
	"strconv"
)

type MyNotice struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	Type            int       `json:"type"`
	Action          string    `json:"action"`
	Content         string    `json:"content"`
	Name            string    `json:"name"`
	Namespace       string    `json:"namespace"`
	MemberId        int64     `json:"member_id"`
	IsRead          int       `json:"is_read"`
	ProclaimType    string    `json:"proclaim_type"`
	ProclaimReceive string    `json:"proclaim_receive"`
	CreatedAt       null.Time `json:"created_at"`
}

type NoticeMemberRepository interface {
	Create(nm *types.NoticeMember) error
	HasRead(noticeId int64) error
	CountMessage(param map[string]string, memberId int64) (count int, err error)
	FindMessageLimit(param map[string]string, memberId int64, offset, limit int) (res []*MyNotice, err error)
	CountRead(param map[string]string, memberId int64) (count int, err error)
	InsertMulti(nml []*types.NoticeMember) (num int, err error)
	ClearAll(noticeType int, memberId int64) (err error)
	Detail(noticeMemberId int, memberId int64) (res *MyNotice, err error)
}

type noticeMember struct {
	db *gorm.DB
}

func NewNoticeMemberRepository(db *gorm.DB) NoticeMemberRepository {
	return &noticeMember{db: db}
}

/**
 * @Title 创建用户关联消息
 */
func (c *noticeMember) Create(nm *types.NoticeMember) error {
	return c.db.Save(nm).Error
}

/**
 * @Title 更新用户消息已读
 */
func (c *noticeMember) HasRead(noticeMemberId int64) error {
	return c.db.Table("notice_member").Where("id = ?", noticeMemberId).Update("is_read", 1).Error
}

/**
 * @Title 获取用户消息总数
 */
func (c *noticeMember) CountMessage(param map[string]string, memberId int64) (count int, err error) {
	query := c.db.Table("notice_member")
	query = query.Joins("left join notices ON notice_member.notice_id=notices.id")
	query = query.Where("notice_member.member_id = ?", memberId)

	if param["title"] != "" {
		query = query.Where("notices.title like ?", "%"+param["title"]+"%")
	}
	if param["type"] != "" {
		query = query.Where("notices.type = ?", param["type"])
	}

	err = query.Count(&count).Error
	return
}

/**
* @Title 获取用户消息列表
 */
func (c *noticeMember) FindMessageLimit(param map[string]string, memberId int64, offset, limit int) (list []*MyNotice, err error) {
	query := c.db.Table("notice_member")
	query = query.Select("notice_member.id, notice_member.is_read, notices.title, notices.content, notices.type, notices.action, notices.name, notices.namespace, notices.member_id, notices.proclaim_type, notices.proclaim_receive, notices.created_at").Joins("left join notices ON notice_member.notice_id=notices.id")
	query = query.Where("notice_member.member_id = ?", memberId)

	if param["title"] != "" {
		query = query.Where("notices.title like ?", "%"+param["title"]+"%")
	}
	if param["type"] != "" {
		query = query.Where("notices.type = ?", param["type"])
	}
	if param["is_read"] != "" {
		query = query.Where("notice_member.is_read = ?", param["is_read"])
	}

	err = query.Order("id desc").Offset(offset).Limit(limit).Scan(&list).Error
	return
}

/**
 * @Title 获取用户消息已阅读/未阅读总数
 */
func (c *noticeMember) CountRead(param map[string]string, memberId int64) (count int, err error) {
	query := c.db.Table("notice_member")
	query = query.Joins("left join notices ON notice_member.notice_id=notices.id")
	query = query.Where("notice_member.member_id = ?", memberId)

	if param["title"] != "" {
		query = query.Where("notices.title like ?", "%"+param["title"]+"%")
	}
	if param["type"] != "" {
		query = query.Where("notices.type = ?", param["type"])
	}
	if param["is_read"] != "" {
		query = query.Where("notice_member.is_read = ?", param["is_read"])
	}

	err = query.Count(&count).Error
	return
}

/**
 * @Title 批量写入
 */
func (c *noticeMember) InsertMulti(nml []*types.NoticeMember) (num int, err error) {
	if len(nml) > 0 {
		sql1 := "insert into notice_member (member_id,notice_id) VALUES "
		str := ""
		for _, v := range nml {
			str += "(" + strconv.FormatInt(v.MemberID, 10) + "," + strconv.FormatInt(v.NoticeID, 10) + "),"
		}
		sql2 := str[0 : len(str)-1]
		sql := sql1 + sql2
		err = c.db.Exec(sql).Error
		return len(nml), err
	}
	return
}

/**
 * @Title 标记为已读
 */
func (c *noticeMember) ClearAll(noticeType int, memberId int64) (err error) {
	return c.db.Exec("update notice_member a left join notices b on a.notice_id=b.id set a.is_read=1 where a.is_read=0 and a.member_id=? and b.type=?", memberId, noticeType).Error
}

func (c *noticeMember) Detail(noticeMemberId int, memberId int64) (res *MyNotice, err error) {
	var data MyNotice
	query := c.db.Table("notice_member")
	query = query.Select("notice_member.id, notice_member.is_read, notices.title, notices.content, notices.type, notices.action, notices.name, notices.namespace, notices.member_id, notices.proclaim_type, notices.proclaim_receive, notices.created_at").Joins("left join notices ON notice_member.notice_id=notices.id")
	query = query.Where("notice_member.member_id = ? and notice_member.id = ?", memberId, noticeMemberId)
	query.Scan(&data)
	return &data, nil
}
