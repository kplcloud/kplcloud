/**
 * Created by GoLand.
 * Email: xzghua@gmail.com
 * Date: 2019-07-17
 * Time: 15:07
 */
package types

import "gopkg.in/guregu/null.v3"

type ConfigEnv struct {
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	EnvDesc   string    `gorm:"column:env_desc" json:"env_desc"`
	EnvKey    string    `gorm:"column:env_key" json:"env_key"`
	EnvVar    string    `gorm:"column:env_var" json:"env_var"`
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Namespace string    `gorm:"column:namespace" json:"namespace"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (c *ConfigEnv) TableName() string {
	return "config_env"
}
