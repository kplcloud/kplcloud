/**
 * @Time : 2021/9/17 3:29 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package account

import (
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kitcache "github.com/icowan/kit-cache"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/opentracing/opentracing-go"
	"net/http"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption, store repository.Repository, cacheSvc kitcache.Service, tracer opentracing.Tracer) http.Handler {
	var ems []endpoint.Middleware

	ems = append(ems, dmw...)
	clusterEms := []endpoint.Middleware{
		middleware.ClusterMiddleware(store, cacheSvc, tracer),
	}
	clusterEms = append(clusterEms, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"UserInfo":   ems,
		"Menus":      ems,
		"Logout":     ems,
		"Namespaces": clusterEms,
	})

	r := mux.NewRouter()

	r.Handle("/user-info", kithttp.NewServer(
		eps.UserInfoEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/menus", kithttp.NewServer(
		eps.MenusEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/logout", kithttp.NewServer(
		eps.LogoutEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/namespaces/{cluster}", kithttp.NewServer(
		eps.NamespacesEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}
