/**
 * @Time : 2019-07-29 11:31
 * @Author : soupzhb@gmail.com
 * @File : transport.go
 * @Software: GoLand
 */

package statistics

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"net/http"
	"strconv"
)

var errBadRoute = errors.New("bad route")

type endpoints struct {
	BuildEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		BuildEndpoint: makeBuildEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		//middleware.ProjectMiddleware(logger, projectRepository),                                   // 4
		//middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Build": ems,
	}

	for _, m := range mw["Build"] {
		eps.BuildEndpoint = m(eps.BuildEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/statistics/{namespace}/build", kithttp.NewServer(
		eps.BuildEndpoint,
		decodeBuildRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeBuildRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, errBadRoute
	}
	pageStr := r.URL.Query().Get("p")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("nameEn")
	etime := r.URL.Query().Get("eTime")
	stime := r.URL.Query().Get("sTime")
	groupId := r.URL.Query().Get("groupId")
	gid, _ := strconv.Atoi(groupId)
	p, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 10
	}
	return buildRequest{
		Namespace: ns,
		Name:      name,
		GroupID:   gid,
		ETime:     etime,
		STime:     stime,
		Page:      p,
		Limit:     limit,
	}, nil
}
