/**
 * @Time : 2019/7/5 11:47 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : configmap
 * @Software: GoLand
 */

package types

import (
	"time"
)

//
type ConfigMap struct {
	Id              int64      `json:"id"`
	ClusterId       int64      `gorm:"column:cluster_id;index;notnull;comment:'集群ID'" json:"cluster_id"`
	Name            string     `gorm:"column:name;index;notnull;size:32;comment:'名称'" json:"name"`
	Namespace       string     `gorm:"column:namespace;index;notnull;size:32;comment:'空间'" json:"namespace"`
	Desc            string     `gorm:"column:desc;null;comment:'备注'" json:"desc"`
	ResourceVersion string     `gorm:"column:resource_version;null;comment:'版本'" json:"resource_version"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	Data []Data `gorm:"ForeignKey:id;AssociationForeignKey:target_id" json:"data"`
}

// TableName sets the insert table name for this struct type
func (c *ConfigMap) TableName() string {
	return "config_map"
}
