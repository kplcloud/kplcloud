/**
 * @Time : 2019/8/5 5:35 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : eventhistory
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type EventHistory struct {
	AppName   string    `gorm:"column:app_name"`
	CreatedAt null.Time `gorm:"column:created_at"`
	Date      string    `gorm:"column:date;size(10000)"`
	Event     string    `gorm:"column:event"`
	ID        int       `gorm:"column:id;primary_key"`
	Namespace string    `gorm:"column:namespace"`
	UpdatedAt null.Time `gorm:"column:updated_at"`
}

// TableName sets the insert table name for this struct type
func (e *EventHistory) TableName() string {
	return "event_history"
}
