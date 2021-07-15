/**
 * @Time : 3/4/21 6:05 PM
 * @Author : solacowa@gmail.com
 * @File : sys_namespace
 * @Software: GoLand
 */

package types

import "time"

// 系统空间，比如各个项目或各个业务线的区分
// 没有分配空间的用户 在本空间的用户无法操作其他空间的东西，有权限也不行
type SysNamespace struct {
	Id        int64      `gorm:"column:id;primary_key" json:"id"`
	Alias     string     `gorm:"column:notnull;alias;comment:'中文名'" json:"alias"`                // 中文名
	Name      string     `gorm:"column:name;notnull;size:24;unique;comment:'名称'" json:"name"`    // 名称
	Enabled   bool       `gorm:"column:enabled;null;default:true;comment:'是否可用'" json:"enabled"` // 是否可用
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"`          // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`          // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`                            // 删除时间
}

// TableName set table
func (*SysNamespace) TableName() string {
	return "sys_namespace"
}
