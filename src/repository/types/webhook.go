/**
 * @Time : 2019/6/27 10:19 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : webhook
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type Webhook struct {
	AppName   string    `gorm:"column:app_name" json:"app_name"`
	AutherID  int       `gorm:"column:auther_id" json:"auther_id"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Namespace string    `gorm:"column:namespace" json:"namespace"`
	Status    int       `gorm:"column:status" json:"status"`
	Target    string    `gorm:"column:target" json:"target"`
	Token     string    `gorm:"column:token" json:"token"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
	URL       string    `gorm:"column:url" json:"url"`

	Member Member   `gorm:"ForeignKey:id;AssociationForeignKey:AutherID" json:"auther"`
	Events []*Event `gorm:"many2many:events_webhooks;ForeignKey:ID;AssociationForeignKey:ID;jointable_foreignkey:webhooks_id;association_jointable_foreignkey:events_id" json:"events"`
}

// TableName sets the insert table name for this struct type
func (w *Webhook) TableName() string {
	return "webhooks"
}
