/**
 * @Time : 2019-06-28 10:57
 * @Author : solacowa@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package deployment

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/repository"
)

func checkServiceMesh(cf *config.Config) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(meshRequest)
			if req.Model == repository.FieldMesh && !cf.GetBool("server", "service_mesh") {
				return nil, ErrDeploymentServiceMesh
			}
			// ErrDeploymentServiceMesh
			return next(ctx, request)
		}
	}
}
