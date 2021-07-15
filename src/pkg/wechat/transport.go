/**
 * @Time : 2019-07-09 18:52
 * @Author : soupzhb@gmail.com
 * @File : transport.go
 * @Software: GoLand
 */

package wechat

import (
	"context"
	"encoding/json"
	"fmt"
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
)

type endpoints struct {
	ReceiveEndpoint  endpoint.Endpoint
	GetQrEndpoint    endpoint.Endpoint
	TestSendEndpoint endpoint.Endpoint
	MenuEndpoint     endpoint.Endpoint
}

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(httpToContest()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}
	opts2 := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(httpToContest()),
	}

	eps := endpoints{
		ReceiveEndpoint:  makeReceiveEndpoint(svc),
		GetQrEndpoint:    makeGetQrEndpoint(svc),
		TestSendEndpoint: makeTestSendEndpoint(svc),
		MenuEndpoint:     makeMenuEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		//middleware.ProjectMiddleware(logger, projectRepository),                                   // 4
		//middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"GetQr":    ems,
		"TestSend": ems,
		"Menu":     ems,
	}

	for _, m := range mw["GetQr"] {
		eps.GetQrEndpoint = m(eps.GetQrEndpoint)
	}
	for _, m := range mw["TestSend"] {
		eps.TestSendEndpoint = m(eps.TestSendEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/wechat/server", kithttp.NewServer(
		eps.ReceiveEndpoint,
		decodeReceiveRequest,
		encodeResponse,
		opts2...,
	)).Methods("GET", "POST")

	r.Handle("/wechat/qr", kithttp.NewServer(
		eps.GetQrEndpoint,
		decodeGetQrRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/wechat/testSend", kithttp.NewServer(
		eps.TestSendEndpoint,
		decodeTestSendRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/wechat/menu", kithttp.NewServer(
		eps.MenuEndpoint,
		decodeMenuRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r

}

func decodeReceiveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetQrRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req qrRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeTestSendRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeMenuRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	data := response.(receiveResponse).Data
	contentType := response.(receiveResponse).ContentType

	fmt.Println("wechat sdk response:", data)

	header := w.Header()
	header["Content-Type"] = []string{contentType}
	w.WriteHeader(200)
	w.Write([]byte(data))

	return nil
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  -1,
		"error": err.Error(),
	})
}
