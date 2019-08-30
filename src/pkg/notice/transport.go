/**
 * @Time : 2019-07-02 10:40
 * @Author : soupzhb@gmail.com
 * @File : transport.go
 * @Software: GoLand
 */

package notice

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
	"github.com/kplcloud/kplcloud/src/util/encode"
	"io/ioutil"
	"net/http"
	"strconv"
)

type endpoints struct {
	ListEndpoint      endpoint.Endpoint
	TipsEndpoint      endpoint.Endpoint
	CountReadEndpoint endpoint.Endpoint
	ClearAllEndpoint  endpoint.Endpoint
	DetailEndpoint    endpoint.Endpoint
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
		ListEndpoint:      makeListEndpoint(svc),
		TipsEndpoint:      makeTipsEndpoint(svc),
		CountReadEndpoint: makeCountReadEndpoint(svc),
		ClearAllEndpoint:  makeClearAllEndpoint(svc),
		DetailEndpoint:    makeDetailEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		//middleware.ProjectMiddleware(logger, projectRepository),                                   // 4
		//middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"List":      ems,
		"Tips":      ems,
		"CountRead": ems,
		"ClearAll":  ems,
		"Detail":    ems,
	}

	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Tips"] {
		eps.TipsEndpoint = m(eps.TipsEndpoint)
	}
	for _, m := range mw["CountRead"] {
		eps.CountReadEndpoint = m(eps.CountReadEndpoint)
	}
	for _, m := range mw["ClearAll"] {
		eps.ClearAllEndpoint = m(eps.ClearAllEndpoint)
	}
	for _, m := range mw["Detail"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/notice/message", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/notice/tips", kithttp.NewServer(
		eps.TipsEndpoint,
		decodeTipsRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/notice/message/readcount", kithttp.NewServer(
		eps.CountReadEndpoint,
		decodeCountReadRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/notice/clear/all", kithttp.NewServer(
		eps.ClearAllEndpoint,
		decodeClearAllRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/notice/detail/{id:[0-9]+}", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeCountReadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	title := r.URL.Query().Get("title")
	t := r.URL.Query().Get("type")
	return noticeListRequest{Title: title, Type: t, Page: 0, Limit: 0}, nil
}

func decodeTipsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	pageStr := r.URL.Query().Get("p")
	limitStr := r.URL.Query().Get("limit")
	title := r.URL.Query().Get("title")
	t := r.URL.Query().Get("type")

	p, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 15
	}
	return noticeListRequest{Title: title, Type: t, Page: p, Limit: limit}, nil
}

func decodeClearAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req clearRequest
	var param map[string]string

	if err = json.Unmarshal(b, &param); err != nil {
		return nil, err
	}

	if param["0"] == "公" && param["1"] == "告" {
		req.ClearType = 1
	} else if param["0"] == "通" && param["1"] == "知" {
		req.ClearType = 2
	} else if param["0"] == "告" && param["1"] == "警" {
		req.ClearType = 3
	}

	return req, nil
}

func decodeDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	noticeId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return detailRequest{
		NoticeId: noticeId,
	}, nil
}
