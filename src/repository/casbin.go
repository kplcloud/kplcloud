/**
 * @Time : 2019-07-12 14:28
 * @Author : solacowa@gmail.com
 * @File : casbin
 * @Software: GoLand
 */

package repository

import (
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type CasbinRepository interface {
	AddCasbin(r *types.Casbin) bool
}

type casbinStruct struct {
	e  *casbin.Enforcer
	db *gorm.DB
}

func NewCasbin(db *gorm.DB) CasbinRepository {
	m := casbin.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)")

	return &casbinStruct{
		db: db,
		e:  casbin.NewEnforcer(m, gormadapter.NewAdapterByDB(db)),
	}
}

func (c *casbinStruct) AddCasbin(r *types.Casbin) bool {
	return c.e.AddPolicy(r.RoleName, r.Path, r.Method)
}
