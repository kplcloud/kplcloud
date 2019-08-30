/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-06-25
 * Time: 17:10
 */
package types

import "gopkg.in/guregu/null.v3"

type Project struct {
	AuditState   int64     `gorm:"column:audit_state" json:"audit_state"`
	CreatedAt    null.Time `gorm:"column:created_at" json:"created_at"`
	Desc         string    `gorm:"column:desc;varchar(500)" json:"desc"`
	ID           int64     `gorm:"column:id;primary_key" json:"id"`
	Language     string    `gorm:"column:language" json:"language"`
	MemberID     int64     `gorm:"column:member_id" json:"member_id"`
	DisplayName  string    `gorm:"column:display_name" json:"display_name"`
	Name         string    `gorm:"column:name" json:"name"`
	Namespace    string    `gorm:"column:namespace" json:"namespace"`
	PublishState int64     `gorm:"column:publish_state" json:"publish_state"`
	Step         int64     `gorm:"column:step" json:"step"`
	UpdatedAt    null.Time `gorm:"column:updated_at" json:"updated_at"`

	Member           Member             `gorm:"ForeignKey:id;AssociationForeignKey:MemberId" json:"member"`
	Groups           []*Groups          `gorm:"many2many:groups_projectss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:groups_id;jointable_foreignkey:projects_id;" json:"groups"`
	ProjectTemplates []*ProjectTemplate `gorm:"foreignkey:ProjectID" json:"project_templates"`
}

// TableName sets the insert table name for this struct type
func (p *Project) TableName() string {
	return "projects"
}
