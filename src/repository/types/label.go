/**
 * @Time : 8/12/21 5:49 PM
 * @Author : solacowa@gmail.com
 * @File : label
 * @Software: GoLand
 */

package types

import "time"

type Label struct {
	Id        int64      `gorm:"column:id;rimary_key" json:"id"`
	Name      string     `gorm:"column:name;notnull;unique;comment:'标签名'" json:"name"`
	Alias     string     `gorm:"column:alias;notnull;comment:'标签别名'" json:"alias"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName set table
func (*Label) TableName() string {
	return "label"
}
