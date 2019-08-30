/**
 * @Time : 2019/7/15 3:43 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : projectjenkins
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type ProjectJenkinsRepository interface {
	CreateOrUpdate(projectJenkins *types.ProjectJenkins) error
	Find(ns, name string) (res *types.ProjectJenkins, err error)
	Delete(ns, name string) error
}

type projectJenkins struct {
	db *gorm.DB
}

func NewProjectJenkins(db *gorm.DB) ProjectJenkinsRepository {
	return &projectJenkins{db: db}
}

func (c *projectJenkins) CreateOrUpdate(projectJenkins *types.ProjectJenkins) error {
	var rs types.ProjectJenkins
	if c.db.First(&rs, "name = ? AND namespace = ?", projectJenkins.Name, projectJenkins.Namespace).RecordNotFound() == true {
		return c.db.Save(projectJenkins).Error
	}
	return c.db.Model(&rs).Where("name = ? AND namespace = ?", projectJenkins.Name, projectJenkins.Namespace).Update(projectJenkins).Error
}

func (c *projectJenkins) Find(ns, name string) (res *types.ProjectJenkins, err error) {
	var rs types.ProjectJenkins
	err = c.db.First(&rs, "name = ? AND namespace = ?", name, ns).Error
	return &rs, err
}

func (c *projectJenkins) Delete(ns, name string) error {
	return c.db.Delete(&types.ProjectJenkins{
		Name:      name,
		Namespace: ns,
	}, "name = ? AND namespace = ?", name, ns).Error
}
