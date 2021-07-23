/**
 * @Time: 2020/9/26 19:10
 * @Author: solacowa@gmail.com
 * @File: setting
 * @Software: GoLand
 */

package types

import "time"

// 系统一些配置,keyvaleu的方式
type SysSetting struct {
	Id          int64      `gorm:"column:id;rimary_key" json:"id"`
	Section     string     `gorm:"column:section;index;notnull;size:24;comment:'Section'" json:"section"`
	Key         string     `gorm:"column:key;index;notnull;size:128;comment:'标识'" json:"key"`           // 标识
	Value       string     `gorm:"column:value;null;size:5000;comment:'值'" json:"value"`                // 值
	Enable      bool       `gorm:"column:enable;null;default:false;comment:'是否启用'" json:"enable"`       // 是否启用
	Description string     `gorm:"column:description;notnull;size:500;comment:'备注'" json:"description"` // 备注
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`                                 // 创建时间
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`                                 // 更新时间
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`                                 // 删除时间
}

// TableName set table
func (*SysSetting) TableName() string {
	return "sys_setting"
}
