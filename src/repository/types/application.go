/**
 * @Time: 2021/10/24 11:22
 * @Author: solacowa@gmail.com
 * @File: application
 * @Software: GoLand
 */

package types

import "time"

// Application 应用主表
type Application struct {
	Id           int64        `gorm:"column:id;primary_key" json:"id"`
	Alias        string       `gorm:"column:alias;notnull;comment:'别名'" json:"alias"`
	Name         string       `gorm:"column:name;index;24;notnull;comment:'名称'" json:"name"`
	Namespace    string       `gorm:"column:namespace;index;24;notnull;;comment:'空间'" json:"namespace"`
	Replicas     int          `gorm:"column:replicas;null;default:1;comment:'pod数量'" json:"replicas"`
	Cpu          int64        `gorm:"column:cpu;null;comment:'基础CPU'" json:"cpu"`
	MaxCpu       int64        `gorm:"column:max_cpu;null;comment:'最大CPU'" json:"max_cpu"`
	Memory       int64        `gorm:"column:memory;null;comment:'基础内存'" json:"memory"`
	MaxMemory    int64        `gorm:"column:max_memory;null;comment:'最大内存'" json:"max_memory"`
	GitRepo      string       `gorm:"column:git_repo;null;comment:'git仓库'" json:"git_repo"`
	Version      string       `gorm:"column:version;null;comment:'版本'" json:"version"`
	Status       string       `gorm:"column:status;null;comment:'状态'" json:"status"`
	Remark       string       `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	AuditStatus  AuditStatus  `gorm:"column:audit_status;null;comment:'审核状态'" json:"audit_status"`
	DeployMethod DeployMethod `gorm:"column:deploy_method;null;default:'git';comment:'部署方式'" json:"deploy_method"`
	GitHash      string       // 构建的时候获取有多少个项目使用的是这个创建，让勾选要更新的项目

	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	Deployment Deployment `json:"deployment"`
}
