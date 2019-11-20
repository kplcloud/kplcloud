/**
 * @Time : 2019-06-26 15:14
 * @Author : solacowa@gmail.com
 * @File : persistentvolumeclaim
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type PersistentVolumeClaim struct {
	AccessModes      string      `gorm:"column:access_modes" json:"access_modes"`
	CreatedAt        null.Time   `gorm:"column:created_at" json:"created_at"`
	Desc             null.String `gorm:"column:desc" json:"desc"`
	Detail           null.String `gorm:"column:detail;type:text" json:"detail"`
	ID               int64       `gorm:"column:id;primary_key" json:"id"`
	Labels           null.String `gorm:"column:labels" json:"labels"`
	Name             string      `gorm:"column:name" json:"name"`
	Namespace        string      `gorm:"column:namespace" json:"namespace"`
	Selector         null.String `gorm:"column:selector" json:"selector"`
	Storage          string      `gorm:"column:storage" json:"storage"`
	StorageClassName null.String `gorm:"column:storage_class_name" json:"storage_class_name"`
	UpdatedAt        null.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (p *PersistentVolumeClaim) TableName() string {
	return "persistentvolumeclaim"
}
