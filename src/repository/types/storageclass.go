/**
 * @Time : 2019-06-26 10:12
 * @Author : solacowa@gmail.com
 * @File : storageclass
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type StorageClass struct {
	CreatedAt         null.Time   `gorm:"column:created_at" json:"created_at"`
	Desc              null.String `gorm:"column:desc" json:"desc"`
	Detail            string      `gorm:"column:detail;size(10000)" json:"detail"`
	ID                int64       `gorm:"column:id;primary_key" json:"id"`
	Name              string      `gorm:"column:name" json:"name"`
	Provisioner       string      `gorm:"column:provisioner" json:"provisioner"`
	ReclaimPolicy     string      `gorm:"column:reclaim_policy" json:"reclaim_policy"`
	UpdatedAt         null.Time   `gorm:"column:updated_at" json:"updated_at"`
	VolumeBindingMode string      `gorm:"column:volume_binding_mode" json:"volume_binding_mode"`
}

// TableName sets the insert table name for this struct type
func (s *StorageClass) TableName() string {
	return "storageclass"
}
