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

type PersistentVolume struct {
}

// PersistentVolumeClaim 存储卷声明信息
type PersistentVolumeClaim struct {
	Id             int64      `gorm:"column:id;primary_key" json:"id"`
	Name           string     `gorm:"column:name;index;notnull;comment:'名称'" json:"name"`
	Namespace      string     `gorm:"column:namespace;size:64;notnull;index;comment:'空间'" json:"namespace"`
	AccessModes    string     `gorm:"column:access_modes;null;comment:'访问模式'" json:"access_modes"`
	Remark         string     `gorm:"column:remark;null;comment:'备注'" json:"desc"`
	Labels         string     `gorm:"column:labels;null;size:1000;comment:'label信息'" json:"labels"`
	RequestStorage string     `gorm:"column:request_storage;notnull;size:255;comment:'请求容量'" json:"request_storage"`
	LimitStorage   string     `gorm:"column:limit_storage;notnull;size:255;comment:'最大容量'" json:"limit_storage"`
	StorageClassId int64      `gorm:"column:storage_class_id;index;notnull;comment:'存储类ID'" json:"storage_class_id"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	StorageClass StorageClass `gorm:"foreignkey:Id;references:StorageClassId"`
}

// TableName sets the insert table name for this struct type
func (p *PersistentVolumeClaim) TableName() string {
	return "persistent_volume_claim"
}
