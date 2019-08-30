/**
 * @Time : 2019-07-09 16:06
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package build

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"io/ioutil"
	"net/http"
	"strconv"
)

type endpoints struct {
	BuildEndpoint            endpoint.Endpoint
	BuildConsoleEndpoint     endpoint.Endpoint
	ReceiverBuildEndpoint    endpoint.Endpoint
	AbortBuildEndpoint       endpoint.Endpoint
	HistoryEndpoint          endpoint.Endpoint
	RollBackEndpoint         endpoint.Endpoint
	BuildConfEndpoint        endpoint.Endpoint
	CronHistoryEndpoint      endpoint.Endpoint
	CronBuildConsoleEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger, repository repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		BuildEndpoint:            makeBuildEndpoint(svc),
		BuildConsoleEndpoint:     makeBuildConsoleEndpoint(svc),
		AbortBuildEndpoint:       makeAbortBuildEndpoint(svc),
		HistoryEndpoint:          makeHistoryEndpoint(svc),
		CronHistoryEndpoint:      makeCronHistoryEndPoint(svc),
		RollBackEndpoint:         makeRollBackEndpoint(svc),
		BuildConfEndpoint:        makeBuildConfEndpoint(svc),
		CronBuildConsoleEndpoint: makeCronBuildConsoleEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()),           // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	emsc := []endpoint.Middleware{
		middleware.CronJobMiddleware(logger, repository.CronJob(), repository.Groups()),
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Build":            ems,
		"BuildConsole":     ems,
		"AbortBuild":       ems,
		"History":          ems,
		"RollBack":         ems,
		"BuildConf":        ems,
		"CronHistory":      emsc,
		"CronBuildConsole": emsc,
	}

	for _, m := range mw["Build"] {
		eps.BuildEndpoint = m(eps.BuildEndpoint)
	}
	for _, m := range mw["BuildConsole"] {
		eps.BuildConsoleEndpoint = m(eps.BuildConsoleEndpoint)
	}
	for _, m := range mw["AbortBuild"] {
		eps.AbortBuildEndpoint = m(eps.AbortBuildEndpoint)
	}
	for _, m := range mw["History"] {
		eps.HistoryEndpoint = m(eps.HistoryEndpoint)
	}
	for _, m := range mw["RollBack"] {
		eps.RollBackEndpoint = m(eps.RollBackEndpoint)
	}
	for _, m := range mw["BuildConf"] {
		eps.BuildConfEndpoint = m(eps.BuildConfEndpoint)
	}
	for _, m := range mw["CronHistory"] {
		eps.CronHistoryEndpoint = m(eps.CronHistoryEndpoint)
	}
	for _, m := range mw["CronBuildConsole"] {
		eps.CronBuildConsoleEndpoint = m(eps.CronBuildConsoleEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/build/jenkins/{namespace}/project/{name}/building", kithttp.NewServer(
		eps.BuildEndpoint,
		decodeBuildRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/build/jenkins/{namespace}/project/{name}/rollback/{id:[0-9]+}", kithttp.NewServer(
		eps.RollBackEndpoint,
		decodeRollbackRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/build/jenkins/{namespace}/project/{name}/history", kithttp.NewServer(
		eps.HistoryEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/build/jenkins/{namespace}/project/{name}/console/{number}", kithttp.NewServer(
		eps.BuildConsoleEndpoint,
		decodeBuildConsoleRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/build/jenkins/{namespace}/project/{name}/abort/{number}", kithttp.NewServer(
		eps.AbortBuildEndpoint,
		decodeAbortBuildRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/build/jenkins/{namespace}/project/{name}/conf", kithttp.NewServer(
		eps.BuildConfEndpoint,
		decodeBuildConfRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/build/jenkins/{namespace}/cronjob/{name}/history", kithttp.NewServer(
		eps.CronHistoryEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/build/jenkins/{namespace}/cronjob/{name}/console/{number}", kithttp.NewServer(
		eps.CronBuildConsoleEndpoint,
		decodeBuildConsoleRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeAbortBuildRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	number, ok := vars["number"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	num, _ := strconv.Atoi(number)

	return abortRequest{
		Number: num,
	}, nil
}

func decodeRollbackRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	number, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	num, _ := strconv.Atoi(number)

	return abortRequest{
		Number: num,
	}, nil
}

func decodeBuildRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := buildRequest{
		GitType: "tag",
		Version: "master",
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeBuildConsoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	number, ok := vars["number"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	start, _ := strconv.Atoi(r.URL.Query().Get("start"))
	num, _ := strconv.Atoi(number)

	return buildConsoleRequest{
		Number: num,
		Start:  start,
	}, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	p := r.URL.Query().Get("page")
	l := r.URL.Query().Get("limit")
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	t := r.URL.Query().Get("types")
	if t != "cronjob" {
		t = "project"
	}

	return listRequest{
		Page:  page,
		Limit: limit,
		Types: t,
	}, nil
}

func decodeBuildConfRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	return confRequest{Name: name, Namespace: ns}, nil
}
