/**
 * @Time : 2019-07-29 10:39
 * @Author : solacowa@gmail.com
 * @File : dockerfile
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type DockerfileRepository interface {
	Create(data *types.Dockerfile) error
	FindById(id int64) (res *types.Dockerfile, err error)
	Update(data *types.Dockerfile) (err error)
	Delete(id int64) (err error)
	FindBy(language []string, status int, name string, offset, limit int) (res []*types.Dockerfile, count int64, err error)
}

type dockerfile struct {
	db *gorm.DB
}

func NewDockerfileRepository(db *gorm.DB) DockerfileRepository {
	return &dockerfile{db}
}

func (c *dockerfile) Create(data *types.Dockerfile) error {
	return c.db.Save(data).Error
}

func (c *dockerfile) FindById(id int64) (res *types.Dockerfile, err error) {
	var rs types.Dockerfile
	err = c.db.Preload("Member", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,username,email")
	}).First(&rs, "id = ?", id).Error
	return &rs, err
}

func (c *dockerfile) Update(data *types.Dockerfile) (err error) {
	data.Member = types.Member{}
	return c.db.Save(data).Error
}

func (c *dockerfile) Delete(id int64) (err error) {
	return c.db.Delete(&types.Dockerfile{ID: id}).Error
}

func (c *dockerfile) FindBy(language []string, status int, name string, offset, limit int) (res []*types.Dockerfile, count int64, err error) {
	query := c.db.Model(&res).Order(gorm.Expr("id DESC"))

	if len(language) > 0 && language[0] != "" {
		query = query.Where("language in (?)", language)
	}

	if status != 0 {
		query = query.Where("status = ?", status)
	}

	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}

	err = query.Count(&count).Find(&res).Offset(offset).Limit(limit).Error
	return
}
