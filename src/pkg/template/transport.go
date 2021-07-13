/**
 * @Time : 2019/6/25 4:07 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package template

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"io/ioutil"
	"net/http"
	"strconv"
)

var errBadRoute = errors.New("bad route")

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	epsMap := map[string]endpoint.Endpoint{
		"get":    makeGetEndpoint(svc),
		"post":   makePostEndpoint(svc),
		"update": makeUpdateEndpoint(svc),
		"delete": makeDeleteEndpoint(svc),
		"list":   makeListEndpoint(svc),
	}
	get := kithttp.NewServer(epsMap["get"], decodeGetRequest, encode.EncodeResponse, opts...)
	post := kithttp.NewServer(epsMap["post"], decodePostRequest, encode.EncodeResponse, opts...)
	update := kithttp.NewServer(epsMap["update"], decodeUpdateRequest, encode.EncodeResponse, opts...)
	deleteTemp := kithttp.NewServer(epsMap["delete"], decodeGetRequest, encode.EncodeResponse, opts...)
	list := kithttp.NewServer(epsMap["list"], decodeListRequest, encode.EncodeResponse, opts...)

	r := mux.NewRouter()
	r.Handle("/template/{id:[0-9]+}", get).Methods("GET")
	r.Handle("/template/{id:[0-9]+}", update).Methods("PUT")
	r.Handle("/template/{id:[0-9]+}", deleteTemp).Methods("DELETE")
	r.Handle("/template", post).Methods("POST")
	r.Handle("/template", list).Methods("GET")

	return r
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	id, err := strconv.ParseInt(name, 0, 64)
	if err != nil {
		return nil, errBadRoute
	}
	return templateRequest{
		Id: int(id),
	}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req templateRequest

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	// get request Id
	id, err := strconv.ParseInt(idstr, 0, 64)
	if err != nil {
		return nil, errBadRoute
	}
	// get request body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req templateRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	req.Id = int(id)
	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	pageStr := r.URL.Query().Get("p")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("name")
	p, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 15
	}
	return templateListRequest{Name: name, Page: p, Limit: limit}, nil
}
