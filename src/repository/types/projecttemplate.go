/**
 * @Time : 2019-06-27 15:17
 * @Author : solacowa@gmail.com
 * @File : projecttemplate
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type Port struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type TemplateField struct {
	Args         []string `json:"args"`
	Command      []string `json:"command"`
	GitAddr      string   `json:"git_addr"`
	BuildPath    string   `json:"build_path"`
	GitType      string   `json:"git_type"`
	GitVersion   string   `json:"git_version"`
	Branch       string   `json:"branch"`
	Image        string   `json:"image"`
	Language     string   `json:"language"`
	Name         string   `json:"name"`
	Namespace    string   `json:"namespace"`
	Ports        []Port   `json:"ports"`
	Replicas     int      `json:"replicas"`
	Mesh         string   `json:"mesh"`
	ResourceType string   `json:"resource_type"`
	Resources    string   `json:"resources"`
	Step         int      `json:"step"`
	PomFile      string   `json:"pom_file"`
}

type ServiceField struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Ports     []Port `json:"ports"`
}

type IngressField struct {
	Namespace string        `json:"namespace"`
	Name      string        `json:"name"`
	Rules     []*RuleStruct `json:"rules"`
}
type RuleStruct struct {
	Domain string   `json:"domain"`
	Paths  []*Paths `json:"paths"`
}
type Paths struct {
	Path        string `json:"path"`
	ServiceName string `json:"serviceName"`
	PortName    int    `json:"port"`
}

type ProjectTemplate struct {
	CreatedAt          null.Time     `gorm:"column:created_at" json:"created_at"`
	Fields             string        `gorm:"column:fields;size(10000)" json:"fields"`
	FinalTemplate      string        `gorm:"column:final_template;type:text;(10000)" json:"final_template"`
	ID                 int64         `gorm:"column:id;primary_key" json:"id"`
	Kind               string        `gorm:"column:kind;size(255)" json:"kind"`
	ProjectID          int64         `gorm:"column:project_id" json:"project_id"`
	State              int64         `gorm:"column:state" json:"state"`
	UpdatedAt          null.Time     `gorm:"column:updated_at" json:"updated_at"`
	Project            Project       `gorm:"ForeignKey:id;AssociationForeignKey:ProjectID"`
	FieldStruct        TemplateField `gorm:"_" json:"field_struct"`
	IngressFieldStruct IngressField  `gorm:"_" json:"field_struct"`
}

// TableName sets the insert table name for this struct type
func (p *ProjectTemplate) TableName() string {
	return "project_template"
}
