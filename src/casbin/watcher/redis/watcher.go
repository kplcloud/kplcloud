/**
 * @Time : 2019-07-16 10:36
 * @Author : solacowa@gmail.com
 * @File : redis
 * @Software: GoLand
 */

package redis

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/persist"
	"github.com/go-redis/redis"
	kplrds "github.com/icowan/redis-client"
	"sync"
)

type Watcher struct {
	rds      kplrds.RedisClient
	callback func(string)
	closed   chan struct{}
	once     sync.Once
	channel  string
}

func NewWatcher(rds kplrds.RedisClient) (persist.Watcher, error) {
	watcher := &Watcher{
		closed:  make(chan struct{}),
		channel: "/casbin",
		rds:     rds,
	}

	go func() {
		for {
			select {
			case <-watcher.closed:
				return
			default:
				err := watcher.subscribe()
				if err != nil {
					watcher.Close()
					fmt.Printf("Failure from Redis subscription: %v", err)
				}
			}
		}
	}()

	return watcher, nil
}

func (c *Watcher) SetUpdateCallback(callback func(string)) error {
	c.callback = callback
	return nil
}

func (c *Watcher) subscribe() error {
	pubSub := c.rds.Subscribe(c.channel)
	defer func() {
		_ = pubSub.Close()
	}()
	for {
		if pubSub == nil {
			return errors.New("pubSub is nil")
		}
		n, err := pubSub.Receive()
		if err != nil {
			return err
		}
		switch n.(type) {
		case *redis.Message:
			if c.callback != nil {
				c.callback(n.(*redis.Message).String())
			}
		case *redis.Subscription:
			if n.(*redis.Subscription).Count == 0 {
				return nil
			}
		}
	}
}

func (c *Watcher) Update() error {
	if err := c.rds.Publish(c.channel, "update rules"); err != nil {
		return err
	}

	return nil
}

func (c *Watcher) Close() {
	c.once.Do(func() {
		close(c.closed)
	})
	//_ = c.rds.Close()
}
