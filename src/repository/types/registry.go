/**
 * @Time : 2021/9/3 12:03 PM
 * @Author : solacowa@gmail.com
 * @File : registry
 * @Software: GoLand
 */

package types

import "time"

// Registry 镜像仓库
type Registry struct {
	Id        int64      `gorm:"column:id;rimary_key" json:"id"`
	Name      string     `gorm:"column:name;notnull;size:32;unique;comment:'名称'" json:"name"`
	Host      string     `gorm:"column:host;notnull;size:64;comment:'仓库地址'" json:"host"`
	Username  string     `gorm:"column:username;notnull;size:64;comment:'仓库用户名'" json:"username"`
	Password  string     `gorm:"column:password;notnull;size:64;comment:'仓库密码'" json:"password"`
	Remark    string     `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName set table
func (*Registry) TableName() string {
	return "registry"
}
