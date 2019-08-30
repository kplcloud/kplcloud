/**
 * @Time : 2019-07-04 15:19
 * @Author : soupzhb@gmail.com
 * @File : role.go
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type Role struct {
	CreatedAt   null.Time    `gorm:"column:created_at" json:"created_at"`
	ID          int64        `gorm:"column:id;primary_key" json:"id"`
	Name        string       `gorm:"column:name" json:"name"`
	State       int          `gorm:"column:state" json:"state"`
	UpdatedAt   null.Time    `gorm:"column:updated_at" json:"updated_at"`
	Description string       `gorm:"column:description" json:"description"`
	Level       int          `gorm:"column:level" json:"level"`
	Permissions []Permission `gorm:"many2many:roles_permissionss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:permission_id;jointable_foreignkey:role_id;"`
	//DisplayName     string       `gorm:"column:display_name"`
	//Code      string    `gorm:"column:code"`
}

// TableName sets the insert table name for this struct type
func (r *Role) TableName() string {
	return "roles"
}

type Level int

const (
	LevelAdmin Level = 100
	LevelOps   Level = 200
)
