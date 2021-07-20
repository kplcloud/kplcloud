/**
 * @Time : 2020/12/28 11:34 AM
 * @Author : solacowa@gmail.com
 * @File : redis
 * @Software: GoLand
 */

package redis

import (
	redisclient "github.com/icowan/redis-client"
	"github.com/opentracing/opentracing-go"
)

func New(hosts, password, prefix string, db int, tracer opentracing.Tracer) (rds redisclient.RedisClient, err error) {
	rds, err = redisclient.NewRedisClient(hosts, password, prefix, db)
	if err == nil && tracer != nil {
		rds = NewRedisMiddleware(rds, tracer)(rds)
	}
	return rds, err
}
