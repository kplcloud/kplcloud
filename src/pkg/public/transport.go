/**
 * @Time : 2019-06-27 13:39
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package public

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"time"
)

type endpoints struct {
	GitHooksEndpoint        endpoint.Endpoint
	PrometheusAlertEndpoint endpoint.Endpoint
	ConfigEndpoint          endpoint.Endpoint
}

const rateBucketNum = 6

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
	}

	eps := endpoints{
		GitHooksEndpoint:        makeGitPostEndpoint(svc),
		PrometheusAlertEndpoint: makePrometheusEndpoint(svc),
		ConfigEndpoint:          makeConfigEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		newTokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)),
	}

	mw := map[string][]endpoint.Middleware{
		//"GitHooks":     ems,
		"PrometheusAlert": ems,
	}

	for _, m := range mw["PrometheusAlert"] {
		eps.PrometheusAlertEndpoint = m(eps.PrometheusAlertEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/public/build/{namespace}/app/{name}", kithttp.NewServer(
		eps.GitHooksEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/public/prometheus/alerts", kithttp.NewServer(
		eps.PrometheusAlertEndpoint,
		decodePromAlertRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/public/config", kithttp.NewServer(
		eps.ConfigEndpoint,
		decodeConfigRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	var req gitlabHook
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	word := r.URL.Query().Get("word")
	if word == "" {
		word = "build"
	}

	branch := r.URL.Query().Get("branch")

	req.Name = name
	req.Namespace = ns
	req.KeyWord = word
	req.Branch = branch
	req.Token = r.Header.Get("X-Git-Token")

	return req, nil
}

func decodePromAlertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)

	fmt.Println("promAlert:", string(b))

	if err != nil {
		return nil, err
	}
	req, err := NewPrometheusAlerts(b)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func decodeConfigRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
