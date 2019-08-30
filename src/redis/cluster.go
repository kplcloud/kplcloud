/**
 * @Time : 2019-06-25 19:38
 * @Author : solacowa@gmail.com
 * @File : cluster
 * @Software: GoLand
 */

package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type cluster struct {
	client *redis.ClusterClient
	//prefix func(s string) string
	prefix string
}

func NewRedisCluster(hosts []string, password string) RedisInterface {
	return &cluster{client: redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    hosts,
		Password: password,
	}), prefix: "kplcloud:"}
}

func (c *cluster) Set(k string, v interface{}, expir ...time.Duration) (err error) {
	var val string
	switch v.(type) {
	case string:
		val = v.(string)
	default:
		b, _ := json.Marshal(v)
		val = string(b)
	}

	exp := expiration
	if len(expir) == 1 {
		exp = expir[0]
	}

	return c.client.Set(c.setPrefix(k), val, exp).Err()
}

func (c *cluster) Get(k string) (v string, err error) {
	return c.client.Get(c.setPrefix(k)).Result()
}

func (c *cluster) Del(k string) (err error) {
	return c.client.Del(c.setPrefix(k)).Err()
}

func (c *cluster) HSet(k string, field string, v interface{}) (err error) {
	var val string
	switch v.(type) {
	case string:
		val = v.(string)
	default:
		b, _ := json.Marshal(v)
		val = string(b)
	}
	return c.client.HSet(c.setPrefix(k), field, val).Err()
}

func (c *cluster) HGet(k string, field string) (res string, err error) {
	return c.client.HGet(c.setPrefix(k), field).Result()
}

func (c *cluster) HDelAll(k string) (err error) {
	res, err := c.client.HKeys(c.setPrefix(k)).Result()
	if err != nil {
		return
	}
	return c.client.HDel(c.setPrefix(k), res...).Err()
}

func (c *cluster) HDel(k string, field string) (err error) {
	return c.client.HDel(c.setPrefix(k), field).Err()
}

func (c *cluster) setPrefix(s string) string {
	return c.prefix + s
}

func (c *cluster) Close() error {
	return c.client.Close()
}

func (c *cluster) Subscribe(channels ...string) *redis.PubSub {
	return c.client.Subscribe(channels...)
}

func (c *cluster) Publish(channel string, message interface{}) error {
	return c.client.Publish(channel, message).Err()
}
