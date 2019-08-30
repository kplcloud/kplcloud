/**
 * @Time : 2019-07-29 11:35
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package market

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
	"github.com/kplcloud/kplcloud/src/util/encode"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type endpoints struct {
	PostEndpoint   endpoint.Endpoint
	DetailEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	PutEndpoint    endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		PostEndpoint:   makePostEndpoint(svc),
		DetailEndpoint: makeDetailEndpoint(svc),
		ListEndpoint:   makeListEndpoint(svc),
		PutEndpoint:    makePutEndpoint(svc),
		DeleteEndpoint: makeDeleteEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Post":     ems,
		"Detail":   ems,
		"List":     ems,
		"Put":      ems,
		"Delete":   ems,
		"Download": ems,
	}

	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["Detail"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Put"] {
		eps.PutEndpoint = m(eps.PutEndpoint)
	}
	for _, m := range mw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mw["Download"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/market/dockerfile/{id:[0-9]+}", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/market/dockerfile", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/market/dockerfile", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/market/dockerfile/{id:[0-9]+}", kithttp.NewServer(
		eps.PutEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/market/dockerfile/{id:[0-9]+}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/market/dockerfile/{id:[0-9]+}/download", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeDetailRequest,
		encodeDownloadLogResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req postRequest

	vars := mux.Vars(r)
	queryId, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	id, err := strconv.ParseInt(queryId, 10, 64)
	if err != nil {
		return nil, encode.ErrBadRoute
	}

	req.Id = id

	return req, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req postRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	language := strings.Split(r.URL.Query().Get("language"), ",")
	name := r.URL.Query().Get("name")
	status, _ := strconv.Atoi(r.URL.Query().Get("status"))

	return listRequest{Page: page, Limit: limit, language: language, name: name, status: status}, nil
}

type response struct {
	Code int           `json:"code"`
	Data io.ReadCloser `json:"data,omitempty"`
	Err  error         `json:"error,omitempty"`
}

func (r response) error() error { return r.Err }

type errorer interface {
	error() error
}

func encodeDownloadLogResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	if e, ok := resp.(errorer); ok && e.error() != nil {
		encode.EncodeError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("content-disposition", "attachment; filename=Eockerfile")
	_, err := io.Copy(w, resp.(response).Data)
	return err
}
