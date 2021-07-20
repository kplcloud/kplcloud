/**
 * @Time : 2020/11/9 2:15 PM
 * @Author : solacowa@gmail.com
 * @File : service_sore
 * @Software: GoLand
 */

package server

import (
	"context"
	"time"

	"github.com/dchest/captcha"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitcache "github.com/icowan/kit-cache"
)

type captchaStore struct {
	cache      kitcache.Service
	expiration time.Duration
	logger     log.Logger
	prefix     string
}

func (s *captchaStore) Set(id string, digits []byte) {
	err := s.cache.Set(context.Background(), s.pre(id), string(digits), s.expiration)
	if err != nil {
		_ = level.Error(s.logger).Log("rds", "set", "id", id, "err", err.Error())
	}
}

func (s *captchaStore) Get(id string, clear bool) (digits []byte) {
	v, err := s.cache.Get(context.Background(), s.pre(id), nil)
	if err != nil {
		_ = level.Error(s.logger).Log("rds", "get", "id", id, "clear", clear, "err", err.Error())
	}
	if clear {
		//_ = s.rds.Del(s.pre(id))
	}

	return []byte(v)
}

func (s *captchaStore) pre(id string) string {
	return s.prefix + id
}

func NewCaptchaStore(cache kitcache.Service, logger log.Logger, expiration time.Duration) captcha.Store {
	return &captchaStore{
		cache:      cache,
		logger:     logger,
		expiration: expiration,
		prefix:     "captcha:",
	}
}
