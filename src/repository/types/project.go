// Package types
package types

import "time"

type Project struct {
	Id          int64  `gorm:"column:id;primary_key" json:"id"`
	Alias       string `gorm:"column:alias;notnull;comment:'别名'" json:"alias"`
	Name        string `gorm:"column:name;index;24;notnull;comment:'名称'" json:"name"`
	Namespace   string `gorm:"column:namespace;index;24;notnull;;comment:'空间'" json:"namespace"`
	Replicas    int    `json:"replicas"`
	Cpu         int64  `gorm:"column:cpu;null;comment:'基础CPU'" json:"cpu"`
	MaxCpu      int64  `gorm:"column:max_cpu;null;comment:'最大CPU'" json:"max_cpu"`
	Memory      int64  `gorm:"column:memory;null;comment:'基础内存'" json:"memory"`
	MaxMemory   int64  `gorm:"column:max_memory;null;comment:'最大内存'" json:"max_memory"`
	GitRepo     string `gorm:"column:git_repo;null;comment:'git仓库'" json:"git_repo"`
	Version     string `gorm:"column:alias;notnull;comment:'别名'" json:"version"`
	Status      string `gorm:"column:alias;notnull;comment:'别名'" json:"status"`
	Remark      string `gorm:"column:alias;notnull;comment:'别名'" json:"remark"`
	AuditStatus int    `gorm:"column:alias;notnull;comment:'别名'" json:"audit_status"`

	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	Deployment Deployment `json:"deployment"`
}

type ProjectContainer struct {
	Id            int64  `json:"id"`
	ProjectId     int64  `json:"project_id"`
	ContainerType int    `json:"container_type"` // 容器类型, initContainer or Container
	Image         string `json:"image"`
	Command       string `json:"command"`
	Args          string `json:"args"`
	Ports         string `json:"ports"`
	Env           string `json:"env"`
	Cpu           int64  `gorm:"column:cpu;null;comment:'基础CPU'" json:"cpu"`
	MaxCpu        int64  `gorm:"column:max_cpu;null;comment:'最大CPU'" json:"max_cpu"`
	Memory        int64  `gorm:"column:memory;null;comment:'基础内存'" json:"memory"`
	MaxMemory     int64  `gorm:"column:max_memory;null;comment:'最大内存'" json:"max_memory"`
}

//import "gopkg.in/guregu/null.v3"
//
//type Project struct {
//	AuditState   int64     `gorm:"column:audit_state" json:"audit_state"`
//	CreatedAt    null.Time `gorm:"column:created_at" json:"created_at"`
//	Desc         string    `gorm:"column:desc;varchar(500)" json:"desc"`
//	ID           int64     `gorm:"column:id;primary_key" json:"id"`
//	Language     string    `gorm:"column:language" json:"language"`
//	MemberID     int64     `gorm:"column:member_id" json:"member_id"`
//	DisplayName  string    `gorm:"column:display_name" json:"display_name"`
//	Name         string    `gorm:"column:name" json:"name"`
//	Namespace    string    `gorm:"column:namespace" json:"namespace"`
//	PublishState int64     `gorm:"column:publish_state" json:"publish_state"`
//	Step         int64     `gorm:"column:step" json:"step"`
//	UpdatedAt    null.Time `gorm:"column:updated_at" json:"updated_at"`
//
//	Member           Member             `gorm:"ForeignKey:id;AssociationForeignKey:MemberId" json:"member"`
//	Groups           []*Groups          `gorm:"many2many:groups_projectss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:groups_id;jointable_foreignkey:projects_id;" json:"groups"`
//	ProjectTemplates []*ProjectTemplate `gorm:"foreignkey:ProjectID" json:"project_templates"`
//}
//
//// TableName sets the insert table name for this struct type
//func (p *Project) TableName() string {
//	return "projects"
//}
