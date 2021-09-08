/**
 * @Time : 2019-06-26 10:12
 * @Author : solacowa@gmail.com
 * @File : storageclass
 * @Software: GoLand
 */

package types

import (
	"time"
)

type StorageClass struct {
	Id                int64      `gorm:"column:id;primary_key" json:"id"`
	ClusterId         int64      `gorm:"column:cluster_id;index;notnull;comment:'集群ID'" json:"cluster_id"`
	Name              string     `gorm:"column:name;index;size:32;notnull;comment:'名称'" json:"name"`
	Provisioner       string     `gorm:"column:provisioner;index;size:32;notnull;comment:'供应商'" json:"provisioner"`
	ReclaimPolicy     string     `gorm:"column:reclaim_policy;notnull;comment:'回收政策'" json:"reclaim_policy"`
	VolumeBindingMode string     `gorm:"column:volume_binding_mode;notnull;comment:'卷绑定模式'" json:"volume_binding_mode"`
	ResourceVersion   string     `gorm:"column:resource_version;null;comment:'版本'" json:"resource_version"`
	Detail            string     `gorm:"column:detail;type:text;size(10000);comment:'详情'" json:"detail"`
	Remark            string     `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt         *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (s *StorageClass) TableName() string {
	return "storage_class"
}
