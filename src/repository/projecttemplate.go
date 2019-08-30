/**
 * @Time : 2019-06-27 15:32
 * @Author : solacowa@gmail.com
 * @File : projecttemplate
 * @Software: GoLand
 */

package repository

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type ProjectTemplateRepository interface {
	Create(projectTemplate *types.ProjectTemplate) error
	FirstOrCreate(projectId int64, kind Kind, fields string, finalYaml string, state int64) (res *types.ProjectTemplate, err error)
	FindByProjectId(projectId int64, kind Kind) (res *types.ProjectTemplate, err error)
	UpdateTemplate(projectTpl *types.ProjectTemplate) error
	CreateOrUpdate(projectTpl *types.ProjectTemplate) error
	FindProjectTemplateByProjectId(projectId int64) (res []*types.ProjectTemplate, err error)
	Count(namespace string, kind Kind) (count int, err error)
	FindOffsetLimit(namespace string, kind Kind, offset, limit int) (res []*types.ProjectTemplate, err error)
	Delete(projectId int64, kind Kind) error
	DeleteByProjectId(projectId int64) error
	UpdateFieldsByNsProjectId(projectId int64, kind Kind, fields, finalTpl string) error
	UpdateProjectTemplate(tpl *types.ProjectTemplate) error
	UpdateState(projectTpl *types.ProjectTemplate) error
}
type Kind string

const (
	Deployment     Kind = "Deployment"
	Service        Kind = "Service"
	VirtualService Kind = "VirtualService"
	Ingress        Kind = "Ingress"
	ConfigMap      Kind = "ConfigMap"

	FieldMesh   = "mesh"
	FieldNormal = "normal"
)

func (c Kind) String() string {
	return string(c)
}

var (
	ErrProjectIdIsNull = errors.New("项目id不能为空")
	ErrKindNull        = errors.New("项目Kind不能为空")
)

type projectTemplate struct {
	db *gorm.DB
}

func NewProjectTemplateRepository(db *gorm.DB) ProjectTemplateRepository {
	return &projectTemplate{db: db}
}

func (c *projectTemplate) UpdateFieldsByNsProjectId(projectId int64, kind Kind, fields, finalTpl string) error {
	return c.db.Model(&types.ProjectTemplate{}).Where("project_id = ? AND kind = ?", projectId, kind.String()).
		Update("fields", fields).Update("final_template", finalTpl).Error
}

func (c *projectTemplate) Create(projectTemplate *types.ProjectTemplate) error {
	return c.db.Save(projectTemplate).Error
}

func (c *projectTemplate) FirstOrCreate(projectId int64, kind Kind, fields string, finalYaml string, state int64) (res *types.ProjectTemplate, err error) {
	class := types.ProjectTemplate{
		Kind:          kind.String(),
		Fields:        fields,
		FinalTemplate: finalYaml,
		State:         state,
	}
	err = c.db.FirstOrCreate(&class, types.ProjectTemplate{
		ProjectID: projectId,
		Kind:      kind.String(),
	}).Error
	return &class, err
}

func (c *projectTemplate) FindByProjectId(projectId int64, kind Kind) (res *types.ProjectTemplate, err error) {
	var rs types.ProjectTemplate
	err = c.db.First(&rs, "project_id = ? AND kind = ?", projectId, kind.String()).Error
	if err == nil {
		switch rs.Kind {
		case Deployment.String():
			var fields types.TemplateField
			_ = json.Unmarshal([]byte(rs.Fields), &fields)
			rs.FieldStruct = fields
		case Ingress.String():
			var fields types.IngressField
			_ = json.Unmarshal([]byte(rs.Fields), &fields)
			rs.IngressFieldStruct = fields
		}

	}
	return &rs, err
}

func (c *projectTemplate) UpdateProjectTemplate(tpl *types.ProjectTemplate) error {
	if tpl.ProjectID == 0 {
		return ErrProjectIdIsNull
	}

	if tpl.Kind == "" {
		return ErrKindNull
	}

	query := c.db.Model(tpl).Where("project_id = ? AND kind = ?", tpl.ProjectID, tpl.Kind)

	if tpl.Fields != "" {
		query.Update("fields", tpl.Fields)
	}
	if tpl.FinalTemplate != "" {
		query.Update("final_template", tpl.FinalTemplate)
	}
	return query.Error
}

func (c *projectTemplate) UpdateTemplate(projectTpl *types.ProjectTemplate) error {
	var fields []byte
	switch projectTpl.Kind {
	case Deployment.String():
		fields, _ = json.Marshal(projectTpl.FieldStruct)
	default:
		fields = []byte(projectTpl.Fields)
	}
	return c.db.Model(projectTpl).Where("id = ? AND kind = ?", projectTpl.ID, projectTpl.Kind).
		Update("fields", string(fields), "final_template", projectTpl.FinalTemplate, "state", projectTpl.State).Error
}

func (c *projectTemplate) CreateOrUpdate(projectTpl *types.ProjectTemplate) error {
	if isNotExist := c.db.Where("project_id = ? AND kind = ?", projectTpl.ProjectID, projectTpl.Kind).
		First(&types.ProjectTemplate{}).RecordNotFound(); isNotExist == true {
		return c.Create(projectTpl)
	} else {
		return c.UpdateTemplate(projectTpl)
	}

}

func (c *projectTemplate) FindProjectTemplateByProjectId(projectId int64) (res []*types.ProjectTemplate, err error) {
	err = c.db.Find(&res, "project_id = ?", projectId).Error
	if err == nil {
		for k, v := range res {
			switch v.Kind {
			case Deployment.String():
				var fields types.TemplateField
				_ = json.Unmarshal([]byte(v.Fields), &fields)
				res[k].FieldStruct = fields
				//res[k].Fields = ""
			case Ingress.String():
				var fields types.IngressField
				_ = json.Unmarshal([]byte(v.Fields), &fields)
				res[k].IngressFieldStruct = fields
				//res[k].Fields = ""
			}

		}
	}
	return
}

func (c *projectTemplate) Count(namespace string, kind Kind) (count int, err error) {
	var projects types.Project
	var projectTemp types.ProjectTemplate

	query := c.db.Model(&projectTemp).Joins("inner join " + projects.TableName() + " t1 on t1.id = " + projectTemp.TableName() + ".project_id")
	if namespace != "" {
		query = query.Where("t1.namespace = ?", namespace)
	}
	if kind.String() != "" {
		query = query.Where("kind = ?", kind.String())
	}
	err = query.Count(&count).Error
	return
}

func (c *projectTemplate) FindOffsetLimit(namespace string, kind Kind, offset, limit int) (list []*types.ProjectTemplate, err error) {
	var projects types.Project
	var projectTemp types.ProjectTemplate
	query := c.db.Model(&projectTemp).Joins("inner join " + projects.TableName() + " t1 on t1.id = " + projectTemp.TableName() + ".project_id")
	if namespace != "" {
		query = query.Where("t1.namespace = ?", namespace)
	}
	if kind.String() != "" {
		query = query.Where("kind = ?", kind.String())
	}
	err = query.Preload("Project").Preload("Project.Member").Order("id desc").Offset(offset).Limit(limit).Find(&list).Error
	if err == nil {
		for k, v := range list {
			switch v.Kind {
			case Deployment.String():
				var fields types.TemplateField
				_ = json.Unmarshal([]byte(v.Fields), &fields)
				list[k].FieldStruct = fields
			case Ingress.String():
				var fields types.IngressField
				_ = json.Unmarshal([]byte(v.Fields), &fields)
				list[k].IngressFieldStruct = fields
			}
		}
	}
	return
}

func (c *projectTemplate) UpdateState(projectTpl *types.ProjectTemplate) error {
	return c.db.Model(projectTpl).Where("id = ?", projectTpl.ID).Update(projectTpl).Error
}

func (c *projectTemplate) Delete(projectId int64, kind Kind) error {
	return c.db.Delete(&types.ProjectTemplate{
		ProjectID: projectId,
		Kind:      kind.String(),
	}, "project_id = ? AND kind = ?", projectId, kind.String()).Error
}

func (c *projectTemplate) DeleteByProjectId(projectId int64) error {
	return c.db.Delete(&types.ProjectTemplate{
		ProjectID: projectId,
	}, "project_id = ?", projectId).Error
}
