/**
 * @Time : 2021/9/3 2:05 PM
 * @Author : solacowa@gmail.com
 * @File : cache
 * @Software: GoLand
 */

package registry

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

func (s *cache) SaveCall(ctx context.Context, reg *types.Registry, call Call) (err error) {
	panic("implement me")
}

func (s *cache) List(ctx context.Context, query string, page, pageSize int) (res []types.Registry, total int, err error) {
	return s.next.List(ctx, query, page, pageSize)
}

func (s *cache) FindByNames(ctx context.Context, names []string) (res []types.Registry, err error) {
	return s.next.FindByNames(ctx, names)
}

func (s *cache) FindByName(ctx context.Context, name string) (res types.Registry, err error) {
	if err = s.kitcache.GetCall(ctx, fmt.Sprintf("%s:%s", s.pkgName, name), func(key string) (res interface{}, err error) {
		return s.next.FindByName(ctx, name)
	}, time.Minute*5, &res); err != nil {
		return res, err
	}
	return res, nil
	//return s.next.FindByName(ctx, name)
}

func (s *cache) Save(ctx context.Context, data *types.Registry) (err error) {
	return s.next.Save(ctx, data)
}

func NewCache(logger log.Logger, traceId string, kitcache kitcache.Service) Middleware {
	return func(next Service) Service {
		return &cache{
			logger:   level.Info(logger),
			next:     next,
			traceId:  traceId,
			kitcache: kitcache,
			pkgName:  "repository:registry",
		}
	}
}
