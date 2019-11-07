/**
 * @Time : 2019-07-02 17:19
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package project

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
	"strings"
)

type endpoints struct {
	PostEndpoint      endpoint.Endpoint
	BasicPostEndpoint endpoint.Endpoint
	ListEndpoint      endpoint.Endpoint
	ListByNsEndpoint  endpoint.Endpoint
	PomFileEndpoint   endpoint.Endpoint
	GitAddrEndpoint   endpoint.Endpoint
	DetailEndpoint    endpoint.Endpoint
	UpdateEndpoint    endpoint.Endpoint
	WorkspaceEndpoint endpoint.Endpoint
	SyncEndpoint      endpoint.Endpoint
	DeleteEndpoint    endpoint.Endpoint
	ConfigEndpoint    endpoint.Endpoint
	MonitorEndpoint   endpoint.Endpoint
	AlertsEndpoint    endpoint.Endpoint
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
		PomFileEndpoint:   makePomFileEndpoint(svc),
		GitAddrEndpoint:   makeGitAddrEndpoint(svc),
		DetailEndpoint:    makeDetailEndpoint(svc),
		ListEndpoint:      makeListEndpoint(svc),
		ListByNsEndpoint:  makeListByNsEndpoint(svc),
		PostEndpoint:      makePostEndpoint(svc),
		BasicPostEndpoint: makeBasicPostEndpoint(svc),
		UpdateEndpoint:    makeUpdateEndpoint(svc),
		WorkspaceEndpoint: makeWorkspaceEndpoint(svc),
		SyncEndpoint:      makeSyncEndpoint(svc),
		DeleteEndpoint:    makeDeleteEndpoint(svc),
		ConfigEndpoint:    makeConfigEndpoint(svc),
		MonitorEndpoint:   makeMonitorEndpoint(svc),
		AlertsEndpoint:    makeAlertsEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()),           // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"PomFile":       ems,
		"GitAddr":       ems,
		"Detail":        ems,
		"List":          ems[1:],
		"ListByNs":      ems[1:],
		"Post":          ems[1:],
		"BasicPost":     ems[1:],
		"Update":        ems,
		"Workspace":     ems[1:],
		"Sync":          ems[1:],
		"DeleteProject": ems,
		"Config":        ems[2:],
		"Monitor":       ems,
		"Alerts":        ems,
	}

	for _, m := range mw["PomFile"] {
		eps.PomFileEndpoint = m(eps.PomFileEndpoint)
	}
	for _, m := range mw["GitAddr"] {
		eps.GitAddrEndpoint = m(eps.GitAddrEndpoint)
	}
	for _, m := range mw["Detail"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["ListByNs"] {
		eps.ListByNsEndpoint = m(eps.ListByNsEndpoint)
	}
	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["BasicPost"] {
		eps.BasicPostEndpoint = m(eps.BasicPostEndpoint)
	}
	for _, m := range mw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}
	for _, m := range mw["Workspace"] {
		eps.WorkspaceEndpoint = m(eps.WorkspaceEndpoint)
	}
	for _, m := range mw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range mw["DeleteProject"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mw["Config"] {
		eps.ConfigEndpoint = m(eps.ConfigEndpoint)
	}
	for _, m := range mw["Monitor"] {
		eps.MonitorEndpoint = m(eps.MonitorEndpoint)
	}
	for _, m := range mw["Alerts"] {
		eps.AlertsEndpoint = m(eps.AlertsEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/project/{namespace}", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/project/{namespace}/sync", kithttp.NewServer(
		eps.SyncEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/project/{namespace}/basic/{name}", kithttp.NewServer(
		eps.BasicPostEndpoint,
		decodeBasicPostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/project/{namespace}", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/project/{namespace}/list", kithttp.NewServer(
		eps.ListByNsEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/project/{namespace}/pomfile/{name}", kithttp.NewServer(
		eps.PomFileEndpoint,
		decodePomFileRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/project/{namespace}/gitaddr/{name}", kithttp.NewServer(
		eps.GitAddrEndpoint,
		decodeGitAddrRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/project/{namespace}/detail/{name}", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/project/{namespace}/detail/{name}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeUpdateRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/project/{namespace}/workspace", kithttp.NewServer(
		eps.WorkspaceEndpoint,
		decodeWorkspaceRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")
	r.Handle("/project/{namespace}/app/{name}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteProjectRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/project/config/detail", kithttp.NewServer(
		eps.ConfigEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/project/{namespace}/monitor/{name}", kithttp.NewServer(
		eps.MonitorEndpoint,
		decodeMonitorRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/project/{namespace}/alerts/{name}", kithttp.NewServer(
		eps.AlertsEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeMonitorRequest(_ context.Context, r *http.Request) (interface{}, error) {
	metrics := r.URL.Query().Get("metrics")
	podName := r.URL.Query().Get("podName")
	container := r.URL.Query().Get("container")
	//if metrics == "" {
	//	return nil, ErrParamsMetrics
	//}

	return monitorRequest{
		metrics,
		podName,
		container,
	}, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := postRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := postRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	if req.Name == "" || req.Namespace == "" {
		return nil, encode.ErrBadRoute
	}
	return req, nil
}

func decodeBasicPostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := basicRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	req.Name = name

	//校验项目地址
	gitAddr := strings.Trim(req.GitAddr, ":")
	gitAddr = strings.Trim(gitAddr, " ")
	if !strings.Contains(gitAddr, ".git") {
		gitAddr += ".git"
	}
	req.GitAddr = gitAddr
	return req, nil
}

func decodeDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	p := r.URL.Query().Get("p")
	l := r.URL.Query().Get("limit")
	g := r.URL.Query().Get("group")

	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	var group int
	if g == "undefined" || g == "" {
		group = 0
	} else {
		group, _ = strconv.Atoi(g)
	}
	if limit == 0 {
		limit = 5
	}

	return listRequest{
		Name:  name,
		Page:  page,
		Limit: limit,
		Group: int64(group),
	}, nil
}

func decodePomFileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := pomFileRequest{
		PomFilePath: "./pom.xml",
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

func decodeGitAddrRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req gitAddrRequest

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeWorkspaceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeDeleteProjectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req deleteRequest
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	code := r.URL.Query().Get("code")

	req.Namespace = ns
	req.Name = name
	req.Code = code

	return req, nil
}
