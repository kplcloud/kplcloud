/**
 * @Time : 2019-07-16 10:33
 * @Author : solacowa@gmail.com
 * @File : redis
 * @Software: GoLand
 */

package redis

import (
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	kplredis "github.com/icowan/redis-client"
)

type adapter struct {
	rds kplredis.RedisClient
}

func NewAdapter(rds kplredis.RedisClient) persist.Adapter {
	return &adapter{rds: rds}
}

func (c *adapter) LoadPolicy(model model.Model) error {

	return nil
}

func (c *adapter) SavePolicy(model model.Model) error {

	return nil
}

func (c *adapter) AddPolicy(sec string, ptype string, rule []string) error {

	return nil
}

func (c *adapter) RemovePolicy(sec string, ptype string, rule []string) error {

	return nil
}

func (c *adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {

	return nil
}
