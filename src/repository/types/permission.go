/**
 * @Time : 2019-07-11 17:37
 * @Author : solacowa@gmail.com
 * @File : permission
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type Permission struct {
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
	Icon      null.String `gorm:"column:icon" json:"icon"`
	ID        int64       `gorm:"column:id;primary_key" json:"id"`
	Menu      null.Bool   `gorm:"column:menu" json:"menu"`
	Method    null.String `gorm:"column:method" json:"method"`
	Name      string      `gorm:"column:name" json:"name"`
	ParentID  null.Int    `gorm:"column:parent_id" json:"parent_id"`
	Path      string      `gorm:"column:path" json:"path"`
	UpdatedAt null.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (p *Permission) TableName() string {
	return "permissions"
}
