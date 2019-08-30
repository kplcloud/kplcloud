/**
 * @Time : 2019-07-12 14:27
 * @Author : solacowa@gmail.com
 * @File : Casbin
 * @Software: GoLand
 */

package types

type Casbin struct {
	ID       int64  `gorm:"-"`
	Ptype    string `json:"ptype" gorm:"ptype"`
	RoleName string `json:"rolename" gorm:"v0"`
	Path     string `json:"path" gorm:"v1"`
	Method   string `json:"method" gorm:"v2"`
}

func (m *Casbin) TableName() string {
	return "casbin_rule"
}
