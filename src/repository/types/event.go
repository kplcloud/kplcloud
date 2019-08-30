/**
 * @Time : 2019/6/27 10:14 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : event
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type Event struct {
	Description null.String `gorm:"column:description"json:"description"`
	ID          int         `gorm:"column:id;primary_key"json:"id"`
	Name        null.String `gorm:"column:name"json:"name"`
	WebHook     []*Webhook  `gorm:"many2many:events_webhooks;ForeignKey:ID;AssociationForeignKey:ID;jointable_foreignkey:events_id;association_jointable_foreignkey:webhooks_id" json:"webhook"`
}

// TableName sets the insert table name for this struct type
func (e *Event) TableName() string {
	return "events"
}
