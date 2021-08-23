/**
 * @Time : 2019/7/5 11:47 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : configdata
 * @Software: GoLand
 */

package types

import (
	"time"
)

type ConfigMapData struct {
	Id          int64      `gorm:"column:id;primary_key" json:"id"`
	ConfigMapId int64      `gorm:"column:config_map_id;notnull;index;comment:'ConfigMapId'" json:"config_map_id"`
	Key         string     `gorm:"column:key;notnull;comment:'Key'" json:"key"`
	Value       string     `gorm:"column:value;null;type:text;comment:'Value'" json:"value"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (c *ConfigMapData) TableName() string {
	return "config_map_data"
}
