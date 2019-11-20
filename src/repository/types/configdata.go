/**
 * @Time : 2019/7/5 11:47 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : configdata
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type ConfigData struct {
	ConfigMapID int64     `gorm:"column:config_map_id" json:"config_map_id"`
	CreatedAt   null.Time `gorm:"column:created_at" json:"created_at"`
	ID          int64     `gorm:"column:id;primary_key" json:"id"`
	Key         string    `gorm:"column:key" json:"key"`
	UpdatedAt   null.Time `gorm:"column:updated_at" json:"updated_at"`
	Value       string    `gorm:"column:value;type:text" json:"value"`
	Path        string    `gorm:"column:path" json:"path"`

	ConfigMap ConfigMap `gorm:"ForeignKey:id;AssociationForeignKey:ConfigMapID"`
}

// TableName sets the insert table name for this struct type
func (c *ConfigData) TableName() string {
	return "config_data"
}
