/**
 * @Time : 2019-07-09 15:54
 * @Author : solacowa@gmail.com
 * @File : build
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type Build struct {
	Address   null.String `gorm:"column:address" json:"address"`
	BuildEnv  null.String `gorm:"column:build_env" json:"build_env"`
	BuildID   null.Int    `gorm:"column:build_id" json:"build_id"`
	BuildTime null.Time   `gorm:"column:build_time" json:"build_time"`
	BuilderID int64       `gorm:"column:builder_id" json:"builder_id"`
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
	Desc      null.String `gorm:"column:desc" json:"desc"`
	GitType   null.String `gorm:"column:git_type" json:"git_type"`
	ID        int64       `gorm:"column:id;primary_key" json:"id"`
	Name      string      `gorm:"column:name" json:"name"`
	Namespace string      `gorm:"column:namespace" json:"namespace"`
	Output    null.String `gorm:"column:output;type:text" json:"output"`
	Path      null.String `gorm:"column:path" json:"path"`
	Status    null.String `gorm:"column:status" json:"status"`
	UpdatedAt null.Time   `gorm:"column:updated_at" json:"updated_at"`
	Version   string      `gorm:"column:version" json:"version"`
	Member    Member      `gorm:"ForeignKey:id;AssociationForeignKey:BuilderID" json:"member"`
}

// TableName sets the insert table name for this struct type
func (b *Build) TableName() string {
	return "builds"
}
