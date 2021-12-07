/**
 * @Time : 3/4/21 6:05 PM
 * @Author : solacowa@gmail.com
 * @File : sys_permission
 * @Software: GoLand
 */

package types

import "time"

type SysPermission struct {
	Id          int64  `gorm:"column:id;primary_key" json:"id"`
	ParentId    int64  `gorm:"column:parent_id;null;default:0;comment:'上级ID'" json:"parent_id"`
	Icon        string `gorm:"column:icon;null;comment:'Icon'" json:"icon"`
	Menu        bool   `gorm:"column:menu;null;default:false;comment:'是否是菜单'" json:"menu"`
	Method      string `gorm:"column:method;null;default:'GET';comment:'Method'" json:"method"`
	Alias       string `gorm:"column:alias;notnull;comment:'名称'" json:"alias"`
	Name        string `gorm:"column:name;unique;size:32;notnull;comment:'标识'" json:"name"`
	Path        string `gorm:"column:path;size:255;notnull;comment:'路径'" json:"path"`
	Description string `gorm:"column:description;null;comment:'描述'" json:"description"`
	Sort        int    `gorm:"column:sort;null;default:0;comment:'排序'" json:"sort"`
	//Methods     []SysPermissionMethod `json:"methods"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

//type SysPermissionMethod struct {
//	Id        int64      `gorm:"column:id;primary_key" json:"id"`
//	Name      string     `gorm:"column:name;unique;comment:'方法: GET POST PUT DELETE'" json:"name"`
//	Alias     string     `gorm:"column:alias;comment:'别名: 查看，删除，增加，修改'" json:"alias"`
//	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
//	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
//	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
//}

// TableName sets the insert table name for this struct type
func (p *SysPermission) TableName() string {
	return "sys_permission"
}
