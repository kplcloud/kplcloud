package types

import (
	"gopkg.in/guregu/null.v3"
)

type NoticeMember struct {
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	IsRead    int       `gorm:"column:is_read" json:"is_read"`
	MemberID  int64     `gorm:"column:member_id" json:"member_id"`
	NoticeID  int64     `gorm:"column:notice_id" json:"notice_id"`
}

// TableName sets the insert table name for this struct type
func (n *NoticeMember) TableName() string {
	return "notice_member"
}
