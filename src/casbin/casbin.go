/**
 * @Time : 2019-07-12 11:06
 * @Author : solacowa@gmail.com
 * @File : base
 * @Software: GoLand
 */

package casbin

import (
	"context"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
	"github.com/casbin/casbin/util"
	"github.com/casbin/gorm-adapter"
	kitcasbin "github.com/go-kit/kit/auth/casbin"
	"github.com/icowan/config"
	redis "github.com/icowan/redis-client"
	"github.com/jinzhu/gorm"
	wrds "github.com/kplcloud/kplcloud/src/casbin/watcher/redis"
	"regexp"
	"strings"
)

type Casbin interface {
	GetContext() context.Context
	GetEnforcer() *casbin.Enforcer
}

type service struct {
	context  context.Context
	enforcer *casbin.Enforcer
	watcher  persist.Watcher
}

var getCasbin *service

func GetCasbin() *service {
	return getCasbin
}

func NewCasbin(cf *config.Config, db *gorm.DB, rds redis.RedisClient) (Casbin, error) {
	m := casbin.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch3(r.obj, p.obj)) && keyMatch2(r.act, p.act)")

	watcher, err := wrds.NewWatcher(rds)
	if err != nil {
		return nil, err
	}

	adapter := gormadapter.NewAdapterByDB(db)
	e, err := casbin.NewEnforcerSafe(m, adapter)
	if err != nil {
		return nil, err
	}

	e.AddFunction("keyMatch3", util.KeyMatch3Func)

	e.SetWatcher(watcher)
	_ = watcher.SetUpdateCallback(func(string) {
		if err = e.LoadPolicy(); err != nil {
			_ = fmt.Errorf("err: %v", err)
		}
	})

	//ctx := context.WithValue(context.Background(), kitcasbin.CasbinModelContextKey, m)
	//ctx = context.WithValue(ctx, kitcasbin.CasbinPolicyContextKey, adapter)
	//ctx = context.WithValue(ctx, kitcasbin.CasbinEnforcerContextKey, e)
	ctx := context.WithValue(context.Background(), kitcasbin.CasbinEnforcerContextKey, e)

	e.EnableLog(cf.GetBool("server", "debug"))

	getCasbin = &service{context: ctx, enforcer: e}
	return getCasbin, nil

}

func (c *service) GetContext() context.Context {
	return c.context
}

func (c *service) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}

func (c *service) Close() {
	c.watcher.Close()
}

func KeyMatch3(key1 string, key2 string) bool {
	re := regexp.MustCompile(`(.*)\{[^/]+\}(.*)`)
	for {
		if !strings.Contains(key2, "/{") {
			break
		}

		key2 = re.ReplaceAllString(key2, "$1[^/]+$2")
	}

	return RegexMatch(key1, key2)
}

func RegexMatch(key1 string, key2 string) bool {
	res, err := regexp.MatchString(key2, key1)
	if err != nil {
		panic(err)
	}
	return res
}
