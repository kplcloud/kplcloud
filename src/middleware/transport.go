/**
 * @Time : 2019-07-02 18:36
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package middleware

import (
	"context"
	"github.com/go-kit/kit/auth/casbin"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kplcasbin "github.com/kplcloud/kplcloud/src/casbin"
	"github.com/kplcloud/kplcloud/src/util/uid"
	stdhttp "net/http"
	"strconv"
	"strings"
)

func NamespaceToContext() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		vars := mux.Vars(r)
		ns, ok := vars["namespace"]
		if ok {
			ctx = context.WithValue(ctx, NamespaceContext, ns)
		}
		name, ok := vars["name"]
		if ok {
			ctx = context.WithValue(ctx, NameContext, name)
		}
		return ctx
	}
}

func GroupIdToContext() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		vars := mux.Vars(r)
		groupId, ok := vars["groupId"]
		if ok {
			groupIdInt, _ := strconv.Atoi(groupId)
			ctx = context.WithValue(ctx, GroupIdContext, int64(groupIdInt))
		}
		return ctx
	}
}

func CasbinToContext() http.RequestFunc {
	return func(ctx context.Context, request *stdhttp.Request) context.Context {
		casbinCtx := kplcasbin.GetCasbin().GetContext()
		ctx = context.WithValue(ctx, casbin.CasbinModelContextKey, casbinCtx.Value(casbin.CasbinModelContextKey))
		ctx = context.WithValue(ctx, casbin.CasbinPolicyContextKey, casbinCtx.Value(casbin.CasbinPolicyContextKey))
		ctx = context.WithValue(ctx, casbin.CasbinEnforcerContextKey, casbinCtx.Value(casbin.CasbinEnforcerContextKey))
		return ctx
	}
}

func CookieToContext() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		if c, err := r.Cookie("Authorization"); err == nil {
			ctx = context.WithValue(ctx, http.ContextKeyRequestAuthorization, c.Value)
			ctx = context.WithValue(ctx, jwt.JWTTokenContextKey, strings.Split(c.Value, "Bearer ")[1])
			r.Header.Set(string(http.ContextKeyRequestAuthorization), c.Value)
		}
		return ctx
	}
}

func RequestIdToContext() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		if r.Header.Get("X-Request-Id") != "" {
			ctx = context.WithValue(ctx, uid.RequestId, r.Header.Get("X-Request-Id"))
		} else {
			ctx = context.WithValue(ctx, uid.RequestId, uid.GenerateUid())
		}
		return ctx
	}
}

func RequestIdToResponse() http.ServerResponseFunc {
	return func(ctx context.Context, w stdhttp.ResponseWriter) context.Context {
		if requestId, ok := ctx.Value(uid.RequestId).(string); ok {
			w.Header().Set("X-Request-Id", requestId)
		}
		return ctx
	}
}

func AllowCors() http.ServerResponseFunc {
	return func(ctx context.Context, w stdhttp.ResponseWriter) context.Context {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type,Authorization,x-requested-with,Access-Control-Allow-Origin,Access-Control-Allow-Credentials")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")
		return ctx
	}
}
