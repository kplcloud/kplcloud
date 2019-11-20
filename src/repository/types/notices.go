package types

import "gopkg.in/guregu/null.v3"

type Notices struct {
	Action          string    `gorm:"column:action" json:"action"`
	Content         string    `gorm:"column:content;type:text" json:"content"`
	CreatedAt       null.Time `gorm:"column:created_at" json:"created_at"`
	ID              int64     `gorm:"column:id;primary_key" json:"id"`
	MemberID        int       `gorm:"column:member_id" json:"member_id"`
	Name            string    `gorm:"column:name" json:"name"`
	Namespace       string    `gorm:"column:namespace" json:"namespace"`
	ProclaimReceive string    `gorm:"column:proclaim_receive" json:"proclaim_receive"`
	ProclaimType    string    `gorm:"column:proclaim_type" json:"proclaim_type"`
	Title           string    `gorm:"column:title" json:"title"`
	Type            int       `gorm:"column:type" json:"type"`
}

type NoticeAction string

const (
	NoticeActionAlarm  NoticeAction = "Alarm"
	NoticeActionDelete NoticeAction = "Delete"
	NoticeActionApply  NoticeAction = "Apply"
	NoticeActionBuild  NoticeAction = "Build"
	NoticeActionAudit  NoticeAction = "Audit"
)

// TableName sets the insert table name for this struct type
func (n *Notices) TableName() string {
	return "notices"
}
