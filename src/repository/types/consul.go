/**
 * @Time : 2019/7/17 2:18 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : consul
 * @Software: GoLand
 */

package types

type Consul struct {
	CreateIndex  int64  `gorm:"column:create_index" json:"create_index"`
	ID           int64  `gorm:"column:id;primary_key" json:"id"`
	ModifyIndex  int64  `gorm:"column:modify_index" json:"modify_index"`
	Name         string `gorm:"column:name" json:"name"`
	Namespace    string `gorm:"column:namespace" json:"namespace"`
	Rules        string `gorm:"column:rules;type:text" json:"rules"`
	Token        string `gorm:"column:token" json:"token"`
	Type         string `gorm:"column:type" json:"type"`
	EncryptToken string `gorm:"-" json:"encrypt_token"`
}

// TableName sets the insert table name for this struct type
func (c *Consul) TableName() string {
	return "consul"
}
