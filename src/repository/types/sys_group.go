/**
 * @Time : 3/4/21 6:05 PM
 * @Author : solacowa@gmail.com
 * @File : sys_group
 * @Software: GoLand
 */

package types

import "time"

// 系统组，与namespace类似，预留还没想好怎么用
type SysGroup struct {
	Id        int64      `gorm:"column:id;primary_key" json:"id"`
	Alias     string     `gorm:"column:notnull;alias;comment:'中文名'" json:"alias"`                // 中文名
	Name      string     `gorm:"column:name;notnull;size:24;unique;comment:'名称'" json:"name"`    // 名称
	Enabled   bool       `gorm:"column:enabled;null;default:true;comment:'是否可用'" json:"enabled"` // 是否可用
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"`          // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`          // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`                            // 删除时间
}

// TableName set table
func (*SysGroup) TableName() string {
	return "sys_group"
}
