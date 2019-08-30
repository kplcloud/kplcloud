package auth

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"time"
)

const rateBucketNum = 3

type endpoints struct {
	LoginEndpoint         endpoint.Endpoint
	GithubEndpoint        endpoint.Endpoint
	GoogleEndpoint        endpoint.Endpoint
	AuthLoginTypeEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	//ctx := context.Background()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	eps := endpoints{
		LoginEndpoint:         makeLoginEndpoint(svc),
		AuthLoginTypeEndpoint: makeAuthLoginTypeEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		newTokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)),
	}

	mw := map[string][]endpoint.Middleware{
		"Login":         ems,
		"AuthLoginType": ems,
	}

	for _, m := range mw["Login"] {
		eps.LoginEndpoint = m(eps.LoginEndpoint)
	}
	for _, m := range mw["AuthLoginType"] {
		eps.AuthLoginTypeEndpoint = m(eps.AuthLoginTypeEndpoint)
	}

	login := kithttp.NewServer(
		eps.LoginEndpoint,
		decodeLoginRequest,
		encodeLoginResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/auth/login", login).Methods("POST")
	r.Handle("/auth/login/type", kithttp.NewServer(
		eps.AuthLoginTypeEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		}, encode.EncodeResponse, opts...))
	r.HandleFunc("/auth/github/callback", svc.AuthLoginGithubCallback).Methods("GET")
	r.HandleFunc("/auth/github/login", svc.AuthLoginGithub).Methods("GET")

	return r
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req authRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	res := response.(authResponse)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Authorization", res.Token)
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    res.Token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7200})
	return json.NewEncoder(w).Encode(response)
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
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  -1,
		"error": err.Error(),
	})
}
