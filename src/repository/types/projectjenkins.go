/**
 * @Time : 2019/7/15 3:38 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : projectjenkins
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type ProjectJenkins struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	Namespace  string    `gorm:"column:namespace" json:"namespace"`
	Command    string    `gorm:"column:command;type:text" json:"command"`
	GitAddr    string    `gorm:"column:git_addr" json:"git_addr"`
	GitType    string    `gorm:"column:git_type" json:"git_type"`
	GitVersion string    `gorm:"column:git_version" json:"git_version"`
	CreatedAt  null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (p *ProjectJenkins) TableName() string {
	return "project_jenkins"
}
