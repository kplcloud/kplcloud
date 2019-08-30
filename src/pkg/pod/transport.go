/**
 * @Time : 2019-07-04 16:17
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package pod

import (
	"context"
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
	"io"
	"net/http"
)

type endpoints struct {
	DetailEndpoint      endpoint.Endpoint
	ProjectPodsEndpoint endpoint.Endpoint
	GetLogEndpoint      endpoint.Endpoint
	DownloadLogEndpoint endpoint.Endpoint
	DeleteEndpoint      endpoint.Endpoint
	PodsMetricsEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger, repository repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(middleware.CookieToContext()),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		DetailEndpoint:      makeDetailEndpoint(svc),
		ProjectPodsEndpoint: makeProjectPodsEndpoint(svc),
		GetLogEndpoint:      makeGetLogEndpoint(svc),
		DownloadLogEndpoint: makeDownloadLogEndpoint(svc),
		DeleteEndpoint:      makeDeleteEndpoint(svc),
		PodsMetricsEndpoint: makePodsMetricsEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()),           // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Detail":      ems,
		"ProjectPods": ems,
		"GetLog":      ems,
		"DownloadLog": ems,
		"Delete":      ems,
		"PodsMetrics": ems,
	}

	for _, m := range mw["Detail"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}
	for _, m := range mw["ProjectPods"] {
		eps.ProjectPodsEndpoint = m(eps.ProjectPodsEndpoint)
	}
	for _, m := range mw["GetLog"] {
		eps.GetLogEndpoint = m(eps.GetLogEndpoint)
	}
	for _, m := range mw["DownloadLog"] {
		eps.DownloadLogEndpoint = m(eps.DownloadLogEndpoint)
	}
	for _, m := range mw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mw["PodsMetrics"] {
		eps.PodsMetricsEndpoint = m(eps.PodsMetricsEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/pods/{namespace}/detail/{name}/pod/{podName}", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/pods/{namespace}/metrics/{name}", kithttp.NewServer(
		eps.PodsMetricsEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/pods/{namespace}/delete/{name}/pod/{podName}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/pods/{namespace}/detail/{name}/pod", kithttp.NewServer(
		eps.ProjectPodsEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/pods/{namespace}/detail/{name}/logs/{podName}/container/{container}", kithttp.NewServer(
		eps.GetLogEndpoint,
		decodeGetLogsRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/pods/{namespace}/detail/{name}/logs/{podName}/container/{container}/download", kithttp.NewServer(
		eps.DownloadLogEndpoint,
		decodeGetLogsRequest,
		encodeDownloadLogResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeGetLogsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	podName, ok := vars["podName"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	container, ok := vars["container"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	return getLogsRequest{
		PodName:   podName,
		Container: container,
	}, nil
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	podName, ok := vars["podName"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	return detailRequest{
		PodName: podName,
	}, nil
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

	logName := ctx.Value(middleware.NameContext).(string)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+logName+".log\"")
	_, err := io.Copy(w, resp.(response).Data)
	return err
}
