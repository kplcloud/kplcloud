/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-06-25
 * Time: 21:08
 */
package types

import "gopkg.in/guregu/null.v3"

type Cronjob struct {
	Active       int       `gorm:"column:active" json:"active"`
	AddType      string    `gorm:"column:add_type" json:"add_type"`
	Args         string    `gorm:"column:args" json:"args"`
	ConfMapName  string    `gorm:"column:conf_map_name" json:"conf_map_name"`
	CreatedAt    null.Time `gorm:"column:created_at" json:"created_at"`
	GitPath      string    `gorm:"column:git_path" json:"git_path"`
	GitType      string    `gorm:"column:git_type" json:"git_type"`
	ID           int64     `gorm:"column:id;primary_key" json:"id"`
	Image        string    `gorm:"column:image" json:"image"`
	LastSchedule null.Time `gorm:"column:last_schedule" json:"last_schedule"`
	LogPath      string    `gorm:"column:log_path" json:"log_path"`
	MemberID     int64     `gorm:"column:member_id" json:"member_id"`
	Name         string    `gorm:"column:name" json:"name"`
	Namespace    string    `gorm:"column:namespace" json:"namespace"`
	Schedule     string    `gorm:"column:schedule" json:"schedule"`
	Suspend      int       `gorm:"column:suspend" json:"suspend"`
	UpdatedAt    null.Time `gorm:"column:updated_at" json:"updated_at"`

	Groups []*Groups `gorm:"many2many:groups_cronjobss;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:groups_id;jointable_foreignkey:cronjobs_id;" json:"groups"`
}

// TableName sets the insert table name for this struct type
func (c *Cronjob) TableName() string {
	return "cronjobs"
}
