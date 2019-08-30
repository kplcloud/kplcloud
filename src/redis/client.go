/**
 * @Time : 2019-06-25 19:26
 * @Author : solacowa@gmail.com
 * @File : client
 * @Software: GoLand
 */

package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/kplcloud/kplcloud/src/config"
	"time"
)

type RedisInterface interface {
	Set(k string, v interface{}, expir ...time.Duration) (err error)
	Get(k string) (v string, err error)
	Del(k string) (err error)
	HSet(k string, field string, v interface{}) (err error)
	HGet(k string, field string) (res string, err error)
	HDelAll(k string) (err error)
	HDel(k string, field string) (err error)
	Close() error
	Subscribe(channels ...string) *redis.PubSub
	Publish(channel string, message interface{}) error
}

const (
	RedisCluster = "cluster"
	RedisSingle  = "single"
	expiration   = 600 * time.Second
)

func NewRedisClient(cf *config.Config) (RedisInterface, error) {
	if cf.GetString("redis", "redis_drive") == RedisCluster {
		return NewRedisCluster(cf.GetStrings("redis", "redis_hosts"), cf.GetString("redis", "redis_password")), nil
	} else if cf.GetString("redis", "redis_drive") == RedisSingle {
		return NewRedisSingle(cf.GetString("redis", "redis_hosts"), cf.GetString("redis", "redis_password"), cf.GetInt("redis", "redis_db")), nil
	}

	return nil, errors.New("redis drive is nil!")
}
