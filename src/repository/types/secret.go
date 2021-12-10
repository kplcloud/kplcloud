/**
 * @Time : 8/19/21 1:52 PM
 * @Author : solacowa@gmail.com
 * @File : secret
 * @Software: GoLand
 */

package types

import "time"

type Secret struct {
	Id              int64      `gorm:"column:id;rimary_key" json:"id"`
	ClusterId       int64      `gorm:"column:cluster_id;notnull;index;comment:'集群ID'" json:"cluster_id"`
	Name            string     `gorm:"column:name;notnull;size:32;index;comment:'名称'" json:"name"`
	Namespace       string     `gorm:"column:namespace;size:32;notnull;index;comment:'空间'" json:"namespace"`
	ResourceVersion string     `gorm:"column:resource_version;null;comment:'版本'" json:"resource_version"`
	Remark          string     `gorm:"column:remark;size:500;null;comment:'备注'" json:"remark"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	Data []Data `gorm:"foreignKey:target_id;references:id" json:"data"`
}

// TableName set table
func (*Secret) TableName() string {
	return "secrets"
}
