/**
 * @Time : 6/15/21 5:07 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package nodes

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type Endpoints struct {
	SyncEndpoint endpoint.Endpoint
	ListEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint: nil,
		SyncEndpoint: makeSyncEndpoint(s),
	}

	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}

	return eps
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = s.Sync(ctx)
		return encode.Response{
			Err: err,
		}, err
	}
}
