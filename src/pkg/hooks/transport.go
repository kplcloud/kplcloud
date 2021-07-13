/**
 * @Time : 2019/6/27 10:10 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package hooks

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
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

func MakeHandler(svc Service, logger kitlog.Logger, repository repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}
	epsMap := map[string]endpoint.Endpoint{
		"get":         makeGetEndpoint(svc),
		"list":        makeListEndpoint(svc),
		"post":        makePostEndpoint(svc),
		"update":      makeUpdateEndpoint(svc),
		"delete":      makeDeleteEndpoint(svc),
		"testSend":    makeTestSendEndpoint(svc),
		"listNoApp":   makeListNoAppEndpoint(svc),
		"postNoApp":   makePostEndpoint(svc),
		"getNoApp":    makeGetEndpoint(svc),
		"updateNoApp": makeUpdateEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()),           // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	emsNoApp := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory),
	}

	mw := map[string][]endpoint.Middleware{
		"get":         ems,
		"list":        ems,
		"post":        ems,
		"update":      ems,
		"delete":      ems,
		"testSend":    ems,
		"listNoApp":   emsNoApp,
		"postNoApp":   emsNoApp,
		"getNoApp":    emsNoApp,
		"updateNoApp": emsNoApp,
	}

	for _, m := range mw["get"] {
		epsMap["get"] = m(epsMap["get"])
	}
	for _, m := range mw["list"] {
		epsMap["list"] = m(epsMap["list"])
	}
	for _, m := range mw["get"] {
		epsMap["post"] = m(epsMap["post"])
	}
	for _, m := range mw["get"] {
		epsMap["update"] = m(epsMap["update"])
	}
	for _, m := range mw["get"] {
		epsMap["delete"] = m(epsMap["delete"])
	}
	for _, m := range mw["get"] {
		epsMap["testSend"] = m(epsMap["testSend"])
	}
	for _, m := range mw["listNoApp"] {
		epsMap["listNoApp"] = m(epsMap["listNoApp"])
	}
	for _, m := range mw["listNoApp"] {
		epsMap["postNoApp"] = m(epsMap["postNoApp"])
	}
	for _, m := range mw["getNoApp"] {
		epsMap["getNoApp"] = m(epsMap["getNoApp"])
	}
	for _, m := range mw["updateNoApp"] {
		epsMap["updateNoApp"] = m(epsMap["updateNoApp"])
	}

	get := kithttp.NewServer(epsMap["get"], decodeGetRequest, encode.EncodeResponse, opts...)
	list := kithttp.NewServer(epsMap["list"], decodeListRequest, encode.EncodeResponse, opts...)
	post := kithttp.NewServer(epsMap["post"], decodePostRequest, encode.EncodeResponse, opts...)
	update := kithttp.NewServer(epsMap["update"], decodeUpdateRequest, encode.EncodeResponse, opts...)
	deleteHook := kithttp.NewServer(epsMap["delete"], decodeGetRequest, encode.EncodeResponse, opts...)
	testSendHook := kithttp.NewServer(epsMap["testSend"], decodeGetRequest, encode.EncodeResponse, opts...)
	listNoApp := kithttp.NewServer(epsMap["listNoApp"], decodeListNoAppRequest, encode.EncodeResponse, opts...)
	postNoApp := kithttp.NewServer(epsMap["postNoApp"], decodePostNoAppRequest, encode.EncodeResponse, opts...)
	getNoApp := kithttp.NewServer(epsMap["getNoApp"], decodeGetNoAppRequest, encode.EncodeResponse, opts...)
	updateNoApp := kithttp.NewServer(epsMap["updateNoApp"], decodeUpdateRequest, encode.EncodeResponse, opts...)

	r := mux.NewRouter()
	r.Handle("/hooks/webhooks/{namespace}/project/{name}/{id:[0-9]+}", get).Methods("GET")
	r.Handle("/hooks/webhooks/{namespace}/project/{name}", list).Methods("GET")
	r.Handle("/hooks/webhooks/{namespace}/project/{name}/{id:[0-9]+}", deleteHook).Methods("DELETE")
	r.Handle("/hooks/webhooks/{namespace}/project/{name}/{id:[0-9]+}", update).Methods("PUT")
	r.Handle("/hooks/webhooks/{namespace}/project/{name}", post).Methods("POST")
	r.Handle("/hooks/webhooks/{namespace}/project/{name}/test-send/{id:[0-9]+}", testSendHook).Methods("POST")
	r.Handle("/hooks/webhooks", listNoApp).Methods("GET")
	r.Handle("/hooks/webhooks", postNoApp).Methods("POST")
	r.Handle("/hooks/webhooks/{id:[0-9]+}", getNoApp).Methods("GET")
	r.Handle("/hooks/webhooks/{id:[0-9]+}", updateNoApp).Methods("PUT")

	return r
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	hookId, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	id, err := strconv.ParseInt(hookId, 0, 64)
	if err != nil {
		return nil, encode.ErrBadRoute
	}

	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	namespace, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	return hookRequest{
		Id:        int(id),
		AppName:   name,
		Namespace: namespace,
	}, nil
}

func decodeGetNoAppRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	hookId, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	id, err := strconv.ParseInt(hookId, 0, 64)
	if err != nil {
		return nil, encode.ErrBadRoute
	}
	return hookRequest{
		Id: int(id),
	}, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	pageStr := r.URL.Query().Get("p")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("name")
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	appName, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	p, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 15
	}
	return hookListRequest{Name: name, AppName: appName, Namespace: ns, Page: p, Limit: limit}, nil
}

func decodeListNoAppRequest(_ context.Context, r *http.Request) (interface{}, error) {
	pageStr := r.URL.Query().Get("p")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("name")
	p, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 15
	}
	return hookListRequest{Name: name, Page: p, Limit: limit}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req hookRequest

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	if req.Name == "" {
		return nil, encode.ErrParamsRefused
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	appName, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	req.AppName = appName
	req.Namespace = ns

	return req, nil
}

func decodePostNoAppRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req hookRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	if req.Name == "" {
		return nil, encode.ErrParamsRefused
	}

	return req, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req hookRequest

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	name, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	id, err := strconv.ParseInt(name, 0, 64)
	if err != nil {
		return nil, encode.ErrBadRoute
	}
	req.Id = int(id)

	return req, nil
}
