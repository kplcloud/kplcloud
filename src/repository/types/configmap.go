/**
 * @Time : 2019/7/5 11:47 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : configmap
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type ConfigMap struct {
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	Desc      string    `gorm:"column:desc" json:"desc"`
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Namespace string    `gorm:"column:namespace" json:"namespace"`
	Type      null.Int  `gorm:"column:type" json:"type"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`

	ConfigData []ConfigData `gorm:"ForeignKey:id;AssociationForeignKey:config_map_id"`
}

// TableName sets the insert table name for this struct type
func (c *ConfigMap) TableName() string {
	return "config_map"
}
