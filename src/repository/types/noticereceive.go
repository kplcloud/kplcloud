package types

import (
	"gopkg.in/guregu/null.v3"
)

type NoticeReceive struct {
	Bee          int       `gorm:"column:bee" json:"bee"`
	CreatedAt    null.Time `gorm:"column:created_at" json:"created_at"`
	Email        int       `gorm:"column:email" json:"email"`
	ID           int64     `gorm:"column:id;primary_key" json:"id"`
	MemberID     int       `gorm:"column:member_id" json:"member_id"`
	NoticeAction string    `gorm:"column:notice_action" json:"notice_action"`
	Site         int       `gorm:"column:site" json:"site"`
	Sms          int       `gorm:"column:sms" json:"sms"`
	UpdatedAt    null.Time `gorm:"column:updated_at" json:"updated_at"`
	Wechat       int       `gorm:"column:wechat" json:"wechat"`
}

// TableName sets the insert table name for this struct type
func (n *NoticeReceive) TableName() string {
	return "notice_receive"
}

