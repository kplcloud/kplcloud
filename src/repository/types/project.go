/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-06-25
 * Time: 17:10
 */
package types

import "time"

type Project struct {
	Id          int64  `json:"id"`
	Alias       string `json:"alias"`
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Cpu         int64  `json:"cpu"`
	MaxCpu      int64  `json:"max_cpu"`
	Memory      int64  `json:"memory"`
	MaxMemory   int64  `json:"max_memory"`
	GitRepo     string `json:"git_repo"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	State       int    `json:"state"`
	Remark      string `json:"remark"`
	AuditStatus int    `json:"audit_status"`
	Step        int    `json:"step"`

	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	Deployment Deployment `json:"deployment"`
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
