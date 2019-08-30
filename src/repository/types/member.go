package types

import "gopkg.in/guregu/null.v3"

type Member struct {
	CreatedAt  null.Time   `gorm:"column:created_at" json:"created_at"`
	Email      string      `gorm:"column:email" json:"email"`
	ID         int64       `gorm:"column:id;primary_key" json:"id"`
	Openid     string      `gorm:"column:openid" json:"openid"`
	Phone      string      `gorm:"column:phone" json:"phone"`
	State      int64       `gorm:"column:state" json:"state"`
	City       string      `gorm:"column:city" json:"city"`
	Department string      `gorm:"column:department" json:"department"`
	UpdatedAt  null.Time   `gorm:"column:updated_at" json:"updated_at"`
	Username   string      `gorm:"column:username" json:"username"`
	Password   null.String `gorm:"column:password" json:"password"`

	Groups     []Groups    `gorm:"many2many:groups_memberss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:groups_id;jointable_foreignkey:members_id;" json:"groups"`
	Namespaces []Namespace `gorm:"many2many:namespaces_memberss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:namespaces_id;jointable_foreignkey:members_id;" json:"namespaces"`
	Roles      []Role      `gorm:"many2many:members_roless;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:roles_id;jointable_foreignkey:members_id;" json:"roles"`
}

// TableName sets the insert table name for this struct type
func (m *Member) TableName() string {
	return "members"
}
