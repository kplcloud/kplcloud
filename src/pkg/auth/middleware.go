/**
 * @Time: 2019-06-29 09:59
 * @Author: solacowa@gmail.com
 * @File: middleware
 * @Software: GoLand
 */

package auth

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("Rate limit exceed!")

func newTokenBucketLimitter(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
