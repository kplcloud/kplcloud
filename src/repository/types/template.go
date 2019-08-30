package types

import "gopkg.in/guregu/null.v3"

type Template struct {
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	Detail    string    `gorm:"column:detail;type:text;size(10000)" json:"detail"`
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	Kind      string    `gorm:"column:kind;size(255)" json:"kind"`
	Name      string    `gorm:"column:name;size(255)" json:"name"`
	Rules     string    `gorm:"column:rules;size(5000)" json:"rules"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (t *Template) TableName() string {
	return "templates"
}
