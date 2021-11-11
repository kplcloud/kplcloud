package middleware

import (
	"context"
	"fmt"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	kitcache "github.com/icowan/kit-cache"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/kplcloud/kplcloud/src/encode"
	asdjwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/repository"
)

type ASDContext string

const (
	UserIdContext    ASDContext = "userId"
	EmailContext     ASDContext = "email"
	NamespaceContext ASDContext = "namespace"
	NameContext      ASDContext = "name"
	GroupIdsContext  ASDContext = "groupIds"
	ProjectContext   ASDContext = "project"
	GroupIdContext   ASDContext = "groupId"
	IsAdmin          ASDContext = "isAdmin"
	CronJobContext   ASDContext = "cronJob"

	ContextKeyClusterName   ASDContext = "ctx-cluster-name"   // 集群名称
	ContextKeyClusterId     ASDContext = "ctx-cluster-id"     // 集群ID
	ContextKeyNamespaceName ASDContext = "ctx-namespace-name" // 空间标识
	ContextKeyNamespaceList ASDContext = "ctx-namespace-list" // 空间标识列表
	ContextKeyNamespaceId   ASDContext = "ctx-namespace-id"   // 空间ID
	ContextKeyName          ASDContext = "ctx-name"           // 名称
	ContextUserId           ASDContext = "ctx-user-id"        // 用户ID
	ContextPermissionId     ASDContext = "ctx-permission-id"  // 权限ID
	ContextKeyClusters      ASDContext = "ctx-clusters"       // 集群列表
)

var (
	ErrProjectNotExists      = errors.New("项目可能不存在")
	ErrCronJobNotExists      = errors.New("定时任务可能不存在")
	ErrNotPermission         = errors.New("没有权限")
	ErrCheckPermissionFailed = errors.New("校验权限失败")
)

func CheckAuthMiddleware(logger log.Logger, cache kitcache.Service, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if tracer != nil {
				var span opentracing.Span
				span, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "CheckAuthMiddleware", opentracing.Tag{
					Key:   string(ext.Component),
					Value: "Middleware",
				})
				defer func() {
					span.LogKV("err", err)
					span.Finish()
				}()
			}

			token := ctx.Value(kithttp.ContextKeyRequestAuthorization).(string)
			if token == "" {
				_ = level.Warn(logger).Log("ctx", "Value", "err", encode.ErrAuthNotLogin.Error())
				return nil, encode.ErrAuthNotLogin.Error()
			}

			var clustom asdjwt.ArithmeticCustomClaims
			tk, err := jwt.ParseWithClaims(token, &clustom, asdjwt.JwtKeyFunc)
			if err != nil || tk == nil {
				_ = level.Error(logger).Log("jwt", "ParseWithClaims", "err", err)
				err = encode.ErrAuthNotLogin.Wrap(err)
				return
			}

			claim, ok := tk.Claims.(*asdjwt.ArithmeticCustomClaims)
			if !ok {
				_ = level.Error(logger).Log("tk", "Claims", "err", ok)
				err = encode.ErrAccountASD.Error()
				return
			}

			// 查询用户是否退出
			if _, err = cache.Get(ctx, fmt.Sprintf("login:%d:token", claim.UserId), nil); err != nil {
				_ = level.Error(logger).Log("cache", "Get", "err", err)
				err = encode.ErrAuthNotLogin.Wrap(err)
				return
			}

			var clusters, namespaces []string
			if _, err := cache.Get(ctx, fmt.Sprintf("user:%d:namespaces", claim.UserId), &namespaces); err == nil {
				ctx = context.WithValue(ctx, ContextKeyNamespaceList, namespaces)
			}
			if _, err := cache.Get(ctx, fmt.Sprintf("user:%d:clusters", claim.UserId), &clusters); err == nil {
				ctx = context.WithValue(ctx, ContextKeyClusters, clusters)
			}
			ctx = context.WithValue(ctx, ContextUserId, claim.UserId)
			ctx = context.WithValue(ctx, "Authorization", token)
			return next(ctx, request)
		}
	}
}

func CheckPermissionMiddleware(logger log.Logger, cacheSvc kitcache.Service, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if tracer != nil {
				var span opentracing.Span
				span, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "CheckPermissionMiddleware", opentracing.Tag{
					Key:   string(ext.Component),
					Value: "Middleware",
				}, opentracing.Tag{
					Key:   "userId",
					Value: ctx.Value(ContextUserId),
				})
				defer func() {
					span.LogKV("err", err)
					span.Finish()
				}()
			}

			userId, ok := ctx.Value(ContextUserId).(int64)
			if !ok {
				err = encode.ErrAccountNotLogin.Error()
				_ = level.Error(logger).Log("userIdContext", "is null")
				return
			}
			var sysUser types.SysUser
			_, err = cacheSvc.Get(ctx, fmt.Sprintf("user:%d:info", userId), &sysUser)
			if err != nil {
				err = encode.ErrAccountASD.Wrap(err)
				_ = level.Error(logger).Log("cacheSvc", "Get", "err", err)
				return
			}

			if sysUser.Locked {
				err = encode.ErrAccountLocked.Error()
				_ = level.Error(logger).Log("sysUser", "Locked", "err", err)
				return
			}

			var roles []types.SysRole
			// TODO: 超管标识,直接通过
			for _, v := range sysUser.SysRoles {
				roles = append(roles, v)
			}

			var requestPath, requestMethod = ctx.Value(kithttp.ContextKeyRequestPath).(string),
				ctx.Value(kithttp.ContextKeyRequestMethod).(string)

			u, _ := url.Parse(ctx.Value(kithttp.ContextKeyRequestURI).(string))
			requestPath = u.Path

			var perms []types.SysPermission
			for _, v := range sysUser.SysRoles {
				for _, p := range v.SysPermissions {
					perms = append(perms, p)
				}
			}

			var pass bool
			var perm types.SysPermission
			for _, v := range perms {
				//fmt.Println("strings.EqualFold(v.Path, requestPath)", strings.EqualFold(v.Path, requestPath), "keyMatch3(requestPath, v.Path)", keyMatch3(requestPath, v.Path))
				//fmt.Println(v.Path, requestPath)
				if !strings.EqualFold(v.Path, requestPath) && !keyMatch3(requestPath, v.Path) {
					continue
				}
				if !strings.EqualFold(v.Method, requestMethod) {
					continue
				}
				perm = v
				pass = true
				break
			}
			if !pass {
				err = encode.ErrAccountASD.Error()
				_ = level.Warn(logger).Log("userId", userId, "requestPath", requestPath, "method", requestMethod, "msg", "权限校验失败")
				return
			}

			ctx = context.WithValue(ctx, ContextPermissionId, perm.Id)

			return next(ctx, request)
		}
	}
}

func keyMatch3(key1 string, key2 string) bool {
	re := regexp.MustCompile(`(.*)\{[^/]+\}(.*)`)
	for {
		if !strings.Contains(key2, "/{") {
			break
		}

		key2 = re.ReplaceAllString(key2, "$1[^/]+$2")
	}
	return regexMatch(key1, key2)
}

func regexMatch(key1 string, key2 string) bool {
	if !strings.Contains(key2, "[^/]") && !strings.EqualFold(key2, key1) {
		return false
	}
	res, err := regexp.MatchString(key2, key1)
	if err != nil {
		panic(err)
	}
	return res
}

func CronJobMiddleware(logger log.Logger, cronjobRepository repository.CronjobRepository, groupsRepository repository.GroupsRepository) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			name := ctx.Value(NameContext).(string)
			namespace := ctx.Value(NamespaceContext).(string)
			if name != "" {
				cronjob, notFound := cronjobRepository.GetCronJobByNameAndNs(name, namespace)
				if notFound {
					_ = level.Error(logger).Log("cronJobRepository", "GetCronJobByNameAndNs", "err", "data not found")
					return nil, ErrCronJobNotExists
				}

				memberId := ctx.Value(UserIdContext).(int64)
				groupIds, _ := ctx.Value(GroupIdsContext).([]int64)
				isAdmin := ctx.Value(IsAdmin).(bool)

				// 如果是超管,直接过
				if !isAdmin {
					// 如果cronjob创建者就是登录用户,直接过
					// 否则看是否项目是否在该登录用户所属组
					if cronjob.MemberID != memberId {
						notFound, err := groupsRepository.CheckPermissionForMidCronJob(cronjob.ID, groupIds)
						if err != nil {
							_ = level.Error(logger).Log("cronJobRepository", "CheckPermissionForMidCronJob", "err", err.Error())
							return nil, ErrCheckPermissionFailed
						}
						if notFound {
							_ = level.Error(logger).Log("cronJobRepository", "CheckPermissionForMidCronJob", "err", "cronjob not in this group")
							return nil, ErrNotPermission
						}
					}
				}

				ctx = context.WithValue(ctx, CronJobContext, cronjob)
				_ = level.Info(logger).Log("CronJobMiddleware", "ctx", "name", name)
				// 如果为 post 或 put 是否需要考虑存历史? 使用defer 但必须有返回结果才行？
			}
			return next(ctx, request)
		}
	}
}

func ProjectMiddleware(logger log.Logger, projectRepository repository.ProjectRepository, groupsRepository repository.GroupsRepository) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			name := ctx.Value(NameContext).(string)
			namespace := ctx.Value(NamespaceContext).(string)
			method := ctx.Value(kithttp.ContextKeyRequestMethod).(string)
			if name != "" {
				project, err := projectRepository.FindByNsNameOnly(namespace, name)
				if err != nil {
					_ = level.Error(logger).Log("projectRepository", "FindByNsName", "err", err.Error())
					return nil, ErrProjectNotExists
				}

				//memberId := ctx.Value(UserIdContext).(int64)
				//groupIds, _ := ctx.Value(GroupIdsContext).([]int64)
				isAdmin := ctx.Value(IsAdmin).(bool)

				// 如果project创建者就是登录用户,直接过
				// 否则看是否项目是否在该登录用户所属组
				// 如果是超管,直接过
				if !isAdmin {
					//if method != http.MethodGet && project.MemberID != memberId {
					//	notFound, err := groupsRepository.CheckPermissionForMidProject(project.ID, groupIds)
					//	if err != nil {
					//		_ = level.Error(logger).Log("groupsRepository", "CheckPermissionForMidProject", "err", err.Error())
					//		return nil, ErrCheckPermissionFailed
					//	}
					//	if notFound {
					//		_ = level.Error(logger).Log("groupsRepository", "CheckPermissionForMidProject", "err", "project not in this group")
					//		return nil, ErrNotPermission
					//	}
					//
					//}
				}

				ctx = context.WithValue(ctx, ProjectContext, project)
				_ = level.Info(logger).Log("ProjectMiddleware", "ctx", "name", name)
			}

			defer func() {
				if err == nil && method != http.MethodGet {
					// todo 后续处理
					// _ = logger.Log("requestId", ctx.Value(uid.RequestId), "userId", ctx.Value(UserIdContext), "err", err.Error())
				}
			}()

			return next(ctx, request)
		}
	}
}
