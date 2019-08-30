/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-06-24
 * Time: 17:40
 */
package types

import "time"

type Groups struct {
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	ID          int64     `gorm:"column:id;primary_key" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	DisplayName string    `gorm:"column:display_name" json:"display_name"`
	Namespace   string    `gorm:"column:namespace" json:"namespace"`
	MemberId    int64     `gorm:"column:member_id" json:"member_id"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`

	Member Member    `gorm:"ForeignKey:id;AssociationForeignKey:MemberId;" json:"member"`
	Ns     Namespace `gorm:"ForeignKey:Name;AssociationForeignKey:Namespace" json:"ns"`

	Members  []Member  `gorm:"many2many:groups_memberss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:members_id;jointable_foreignkey:groups_id;" json:"members"`
	Cronjobs []Cronjob `gorm:"many2many:groups_cronjobss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:cronjobs_id;jointable_foreignkey:groups_id;" json:"cronjobs"`
	Projects []Project `gorm:"many2many:groups_projectss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:projects_id;jointable_foreignkey:groups_id;" json:"projects"`
}

// TableName sets the insert table name for this struct type
func (g *Groups) TableName() string {
	return "groups"
}
