/**
 * @Time : 8/13/21 2:27 PM
 * @Author : solacowa@gmail.com
 * @File : cache
 * @Software: GoLand
 */

package cluster

import (
	"fmt"
	kitcache "github.com/icowan/kit-cache"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"time"
)

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type cache struct {
	logger   log.Logger
	next     Service
	traceId  string
	pkgName  string
	kitcache kitcache.Service
}

func (s *cache) SaveRole(ctx context.Context, clusterRole *types.ClusterRole, roles []types.PolicyRule) (err error) {
	return s.next.SaveRole(ctx, clusterRole, roles)
}

func (s *cache) FindAll(ctx context.Context, status int) (res []types.Cluster, err error) {
	return s.next.FindAll(ctx, status)
}

func (s *cache) FindByName(ctx context.Context, name string) (res types.Cluster, err error) {
	if err = s.kitcache.GetCall(ctx, fmt.Sprintf("%s:%s", s.pkgName, name), func(key string) (res interface{}, err error) {
		return s.next.FindByName(ctx, name)
	}, time.Minute*5, &res); err != nil {
		return res, err
	}
	return res, nil
	//return s.next.FindByName(ctx, name)
}

func (s *cache) Save(ctx context.Context, data *types.Cluster, calls ...Call) (err error) {
	return s.next.Save(ctx, data, calls...)
}

func (s *cache) Delete(ctx context.Context, id int64, unscoped bool) (err error) {
	defer func() {
		// TODO: name 可能删不掉哦
		if err = s.kitcache.Del(ctx, fmt.Sprintf("%s:%d", s.pkgName, id)); err != nil {
			_ = level.Error(s.logger).Log("kitcache", "Del", "err", err.Error())
		}
	}()
	return s.next.Delete(ctx, id, unscoped)
}

func NewCache(logger log.Logger, traceId string, kitcache kitcache.Service) Middleware {
	return func(next Service) Service {
		return &cache{
			logger:   level.Info(logger),
			next:     next,
			traceId:  traceId,
			kitcache: kitcache,
			pkgName:  "repository:cache",
		}
	}
}
