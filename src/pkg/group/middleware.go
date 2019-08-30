/**
 * Created by GoLand.
 * Email: xzghua@gmail.com
 * Date: 2019-07-05
 * Time: 10:49
 */
package group

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
)

// owner
// 只有组长才能操作
// 组长必须有当前ns权限(算了,就创建一个方法要判断,就不加了)

// admin
// 必须是超管才能访问,否则不允许访问

var (
	ErrLimitAdminVisit = errors.New("limit admin visit ")
	ErrGetGroupInfById = errors.New("get Group info by id failed")
	ErrDoForbid        = errors.New("Permission denied ")
)

func adminVisitMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			isAdmin := ctx.Value(middleware.IsAdmin).(bool)
			if !isAdmin {
				_ = logger.Log("adminVisitMiddleware", "user is not admin")
				return nil, ErrLimitAdminVisit
			}
			return next(ctx, request)
		}
	}
}

func ownerDoGroupMiddleware(logger log.Logger, groupRepository repository.GroupsRepository) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			groupId := ctx.Value(middleware.GroupIdContext).(int64)
			group, err := groupRepository.Find(groupId)
			if err != nil {
				_ = logger.Log("ownerDoGroupMiddleware", "Find Group by id", "err", err.Error())
				return nil, ErrGetGroupInfById
			}

			memberId := ctx.Value(middleware.UserIdContext).(int64)
			if group.MemberId != memberId {
				_ = logger.Log("ownerDoGroupMiddleware", "compare with group owner", "err", "Permission denied")
				return nil, ErrDoForbid
			}

			return next(ctx, request)
		}
	}
}
