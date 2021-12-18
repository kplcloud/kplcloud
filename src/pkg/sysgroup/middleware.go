/**
 * @Time : 2021/12/17 2:33 PM
 * @Author : solacowa@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package sysgroup

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"strings"
)

type groupMiddleware string

const (
	ctxGroupName groupMiddleware = "ctx-group-name"
	ctxGroupIds  groupMiddleware = "ctx-group-ids"
)

// 加此中间件将会对group进行校验
func checkGroupMiddleware(store repository.Repository) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			userId := ctx.Value(middleware.ContextUserId).(int64)
			clusterId := ctx.Value(middleware.ContextKeyClusterId).(int64)
			namespace := ctx.Value(middleware.ContextKeyNamespaceName).(string)
			groupName, ok := ctx.Value(ctxGroupName).(string)
			sysUser, err := store.SysUser().Find(ctx, userId, "SysGroups", "SysRoles")
			if err != nil {
				return nil, err
			}
			var groupIds []int64
			if sysUser.IsAdmin() {
				groupIds, _ = store.SysGroup(ctx).FindIds(ctx, clusterId, namespace)
			} else {
				groupIds = sysUser.GroupIds()
				if ok && !strings.EqualFold(groupName, "") {
					group, err := store.SysGroup(ctx).FindByName(ctx, clusterId, namespace, groupName)
					if err != nil {
						err = encode.ErrSysGroupNotfound.Wrap(err)
						return nil, err
					}
					var paas bool
					for _, v := range groupIds {
						if group.Id == v {
							paas = true
							break
						}
					}
					if !paas {
						return nil, encode.ErrSysGroupNotPermission.Error()
					}
					if group.OnlyRead {
						// TODO 如果为只读，这个是在上层中间件处理
					}
				}
			}
			ctx = context.WithValue(ctx, ctxGroupIds, groupIds)
			return next(ctx, request)
		}
	}
}
