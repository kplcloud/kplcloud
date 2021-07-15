/**
 * @Time : 3/4/21 6:04 PM
 * @Author : solacowa@gmail.com
 * @File : role
 * @Software: GoLand
 */

package types

import "time"

type SysRole struct {
	Id          int64      `json:"id"`
	Alias       string     `gorm:"column:alias;notnull;comment:'名称'" json:"alias"`
	Name        string     `gorm:"column:name;notnull;size:24;unique;comment:'标识'" json:"name"`
	Enabled     bool       `gorm:"column:enabled;null;default:true;comment:'是否可用'" json:"enabled"`
	Description string     `gorm:"column:description;null;comment:'描述'" json:"description"`
	Level       int        `gorm:"column:level;null;comment:'等级'" json:"level"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	SysPermissions []SysPermission `gorm:"many2many:sys_role_permissions;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:permission_id;jointable_foreignkey:role_id;"`
}

// TableName sets the insert table name for this struct type
func (r *SysRole) TableName() string {
	return "sys_role"
}
