package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type NamespaceRepository interface {
	Find(name string) (res *types.Namespace, err error)
	Create(ns *types.Namespace) error
	UserMyNsList(memberId int64) ([]types.Namespace, error)
	FindByNames(names []string) (res []*types.Namespace, err error)
	FindAll() (res []*types.Namespace, err error)
}

type ns struct {
	db *gorm.DB
}

func NewNamespaceRepository(db *gorm.DB) NamespaceRepository {
	return &ns{db: db}
}

func (c *ns) Find(name string) (res *types.Namespace, err error) {
	var ns types.Namespace
	if err = c.db.First(&ns, "name = ?", name).Error; err != nil {
		return
	}
	return &ns, nil
}

func (c *ns) Create(ns *types.Namespace) error {
	return c.db.Save(ns).Error
}

func (c *ns) UserMyNsList(memberId int64) ([]types.Namespace, error) {
	var m types.Member
	err := c.db.Preload("Namespaces").First(&m, memberId).Error
	return m.Namespaces, err
}

func (c *ns) FindByNames(names []string) (res []*types.Namespace, err error) {
	err = c.db.Find(&res, "name in (?)", names).Error
	return
}
func (c *ns) FindAll() (res []*types.Namespace, err error) {
	err = c.db.Find(&res).Error
	return
}
