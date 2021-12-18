/**
 * @Time : 3/4/21 6:05 PM
 * @Author : solacowa@gmail.com
 * @File : sys_group
 * @Software: GoLand
 */

package types

import "time"

// SysGroup 系统组，与namespace类似，预留还没想好怎么用
// 大概想好了怎么用
type SysGroup struct {
	Id        int64      `gorm:"column:id;primary_key" json:"id"`
	ClusterId int64      `gorm:"column:cluster_id;index;notnull;comment:'集群ID'" json:"cluster_id"`
	Namespace string     `gorm:"column:namespace;index;size:32;notnull;comment:'空间'" json:"namespace"`
	Name      string     `gorm:"column:name;unique;size:64;notnull;comment:'名称'" json:"name"`
	Alias     string     `gorm:"column:alias;size:64;notnull;comment:'别名'" json:"alias"`
	Remark    string     `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	UserId    int64      `gorm:"column:user_id;index;notnull;comment:'用户ID'" json:"user_id"`
	OnlyRead  bool       `gorm:"column:only_read;null;default:false;comment:'是否只读'" json:"only_read"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"` // 删除时间

	User  SysUser       `gorm:"ForeignKey:id;AssociationForeignKey:user_id;" json:"user"`
	Users []SysUser     `gorm:"many2many:sys_group_users;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:sys_user_id;jointable_foreignkey:group_id;" json:"users"`
	Apps  []Application `gorm:"many2many:sys_group_apps;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:app_id;jointable_foreignkey:group_id;" json:"apps"`
	// CronJobs
	// StatefulSets
}

// TableName set table
func (*SysGroup) TableName() string {
	return "sys_group"
}
