package types

import "gopkg.in/guregu/null.v3"

type Namespace struct {
	CreatedAt   null.Time `gorm:"column:created_at" json:"created_at"`
	ID          int64     `gorm:"column:id;primary_key" json:"id"`
	DisplayName string    `gorm:"column:display_name" json:"display_name"`
	Name        string    `gorm:"column:name" json:"name"`
	UpdatedAt   null.Time `gorm:"column:updated_at" json:"updated_at"`

	Members []Member `gorm:"many2many:namespaces_memberss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:members_id;jointable_foreignkey:namespaces_id"`
}

// TableName sets the insert table name for this struct type
func (n *Namespace) TableName() string {
	return "namespaces"
}
