/**
 * @Time : 2019-06-26 15:14
 * @Author : solacowa@gmail.com
 * @File : persistentvolumeclaim
 * @Software: GoLand
 */

package types

import (
	"time"
)

type PersistentVolumeClaim struct {
	Id               int64      `gorm:"column:id;primary_key" json:"id"`
	Name             string     `gorm:"column:name" json:"name"`
	Namespace        string     `gorm:"column:namespace" json:"namespace"`
	AccessModes      string     `gorm:"column:access_modes" json:"access_modes"`
	Desc             string     `gorm:"column:desc" json:"desc"`
	Detail           string     `gorm:"column:detail;type:text" json:"detail"`
	Labels           string     `gorm:"column:labels" json:"labels"`
	Selector         string     `gorm:"column:selector" json:"selector"`
	Storage          string     `gorm:"column:storage" json:"storage"`
	StorageClassName string     `gorm:"column:storage_class_name" json:"storage_class_name"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (p *PersistentVolumeClaim) TableName() string {
	return "persistent_volume_claim"
}
