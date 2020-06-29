/**
 * @Time : 2020/6/23 6:13 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package virtualservice

import (
	"context"
	"github.com/go-kit/kit/endpoint"

	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type Endpoints struct {
	GetEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GetEndpoint: makeGetEndpoint(s),
	}
	//
	//for _, m := range dmw["Get"] {
	//	eps.GetEndpoint = m(eps.GetEndpoint)
	//}

	return eps
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ns, ok := ctx.Value(middleware.NamespaceContext).(string)
		if !ok {
			return nil, encode.ErrBadRoute
		}
		name, ok := ctx.Value(middleware.NameContext).(string)
		if !ok {
			return nil, encode.ErrBadRoute
		}
		s.Get(ctx, ns, name)
		return encode.Response{Err: err}, err
	}
}
