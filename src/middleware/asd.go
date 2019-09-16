package middleware

import (
	"context"
	"errors"
	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	kitcasbin "github.com/go-kit/kit/auth/casbin"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/repository"
	"net/http"
	"strconv"
	"strings"
)

var ErrorASD = errors.New("权限验证失败！")

type checkRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type ASDContext string

const (
	UserIdContext     ASDContext = "userId"
	EmailContext      ASDContext = "email"
	NamespaceContext  ASDContext = "namespace"
	NameContext       ASDContext = "name"
	RoleIdsContext    ASDContext = "roleIds"
	GroupIdsContext   ASDContext = "groupIds"
	NamespacesContext ASDContext = "namespaces"
	ProjectContext    ASDContext = "project"
	GroupIdContext    ASDContext = "groupId"
	IsAdmin           ASDContext = "isAdmin"
	CronJobContext    ASDContext = "cronJob"
	StartTime         ASDContext = "start-time"
)

var (
	ErrProjectNotExists      = errors.New("项目可能不存在")
	ErrCronJobNotExists      = errors.New("定时任务可能不存在")
	ErrNotPermission         = errors.New("没有权限")
	ErrCheckPermissionFailed = errors.New("校验权限失败")
)

func CheckAuthMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			token := ctx.Value(kithttp.ContextKeyRequestAuthorization).(string)

			if token == "" {
				return nil, ErrorASD
			}
			token = strings.Split(token, "Bearer ")[1]

			var clustom kpljwt.ArithmeticCustomClaims
			tk, err := jwt.ParseWithClaims(token, &clustom, kpljwt.JwtKeyFunc)
			if err != nil || tk == nil {
				_ = level.Error(logger).Log("jwt", "ParseWithClaims", "err", err)
				return
			}

			claim, ok := tk.Claims.(*kpljwt.ArithmeticCustomClaims)
			if !ok {
				_ = level.Error(logger).Log("tk", "Claims", "err", ok)
				err = ErrorASD
				return
			}

			ctx = context.WithValue(ctx, UserIdContext, claim.UserId)
			ctx = context.WithValue(ctx, EmailContext, claim.Name)
			ctx = context.WithValue(ctx, NamespacesContext, claim.Namespaces)
			ctx = context.WithValue(ctx, GroupIdsContext, claim.Groups)
			ctx = context.WithValue(ctx, IsAdmin, claim.IsAdmin)
			ctx = context.WithValue(ctx, RoleIdsContext, claim.RoleIds)

			if !claim.IsAdmin {
				var path, method = ctx.Value(kithttp.ContextKeyRequestPath), ctx.Value(kithttp.ContextKeyRequestMethod).(string)
				var casbinErr error
				var ctx2 interface{}
				for _, groupId := range claim.RoleIds {
					// casbin
					mw := NewEnforcer(strconv.Itoa(int(groupId)), path, method)(func(ctx context.Context, request interface{}) (response interface{}, err error) {
						return ctx, nil
					})
					ctx2, casbinErr = mw(ctx, request)
					if casbinErr == nil {
						break
					}
				}
				if casbinErr != nil || ctx2 == nil {
					return nil, kitcasbin.ErrUnauthorized
				}
				return next(ctx2.(context.Context), request)
			}
			return next(ctx, request)
		}
	}
}

func NamespaceMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var namespace, name string
			namespace, _ = ctx.Value(NamespaceContext).(string)
			name, _ = ctx.Value(NameContext).(string)

			var permission bool

			var namespaces []string
			if ctx.Value(NamespacesContext) != nil {
				namespaces = ctx.Value(NamespacesContext).([]string)
			}

			for _, v := range namespaces {
				if v == namespace {
					permission = true
					break
				}
			}

			if !permission {
				_ = level.Error(logger).Log("name", name, "namespace", namespace, "permission", permission)
				return nil, ErrorASD
			}

			return next(ctx, request)
		}
	}
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

				memberId := ctx.Value(UserIdContext).(int64)
				groupIds, _ := ctx.Value(GroupIdsContext).([]int64)
				isAdmin := ctx.Value(IsAdmin).(bool)

				// 如果project创建者就是登录用户,直接过
				// 否则看是否项目是否在该登录用户所属组
				// 如果是超管,直接过
				if !isAdmin {
					if method != http.MethodGet && project.MemberID != memberId {
						notFound, err := groupsRepository.CheckPermissionForMidProject(project.ID, groupIds)
						if err != nil {
							_ = level.Error(logger).Log("groupsRepository", "CheckPermissionForMidProject", "err", err.Error())
							return nil, ErrCheckPermissionFailed
						}
						if notFound {
							_ = level.Error(logger).Log("groupsRepository", "CheckPermissionForMidProject", "err", "project not in this group")
							return nil, ErrNotPermission
						}

					}
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

func NewEnforcer(
	subject string, object interface{}, action string,
) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (
			response interface{}, err error,
		) {
			enforcer := ctx.Value(kitcasbin.CasbinEnforcerContextKey).(*casbin.Enforcer)
			if !enforcer.Enforce(subject, object, action) {
				return nil, kitcasbin.ErrUnauthorized
			}
			return next(ctx, request)
		}
	}
}
