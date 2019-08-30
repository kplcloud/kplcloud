/**
 * @Time : 2019-06-25 19:38
 * @Author : solacowa@gmail.com
 * @File : single
 * @Software: GoLand
 */

package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type single struct {
	client *redis.Client
	prefix string
}

func NewRedisSingle(host, password string, db int) RedisInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	return &single{client: client, prefix: "kplcloud:"}
}

func (c *single) Set(k string, v interface{}, expir ...time.Duration) (err error) {
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

func (c *single) Get(k string) (v string, err error) {
	return c.client.Get(c.setPrefix(k)).Result()
}

func (c *single) Del(k string) (err error) {
	return c.client.Del(c.setPrefix(k)).Err()
}

func (c *single) HSet(k string, field string, v interface{}) (err error) {
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

func (c *single) HGet(k string, field string) (res string, err error) {
	return c.client.HGet(c.setPrefix(k), field).Result()
}

func (c *single) HDelAll(k string) (err error) {
	res, err := c.client.HKeys(c.setPrefix(k)).Result()
	if err != nil {
		return
	}
	return c.client.HDel(c.setPrefix(k), res...).Err()
}

func (c *single) HDel(k string, field string) (err error) {
	return c.client.HDel(c.setPrefix(k), field).Err()
}

func (c *single) setPrefix(s string) string {
	return c.prefix + s
}

func (c *single) Close() error {
	return c.client.Close()
}

func (c *single) Subscribe(channels ...string) *redis.PubSub {
	return c.client.Subscribe(channels...)
}

func (c *single) Publish(channel string, message interface{}) error {
	return c.client.Publish(channel, message).Err()
}
