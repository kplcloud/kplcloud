/**
 * @Time : 2019-07-09 11:33
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package git

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
	"net/http"
)

type endpoints struct {
	TagsEndpoint              endpoint.Endpoint
	BranchesEndpoint          endpoint.Endpoint
	TagsByGitPathEndpoint     endpoint.Endpoint
	BranchesByGitPathEndpoint endpoint.Endpoint
	GetDockerfileEndpoint     endpoint.Endpoint
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
		TagsEndpoint:              makeTagsEndpoint(svc),
		BranchesEndpoint:          makeBranchesEndpoint(svc),
		TagsByGitPathEndpoint:     makeTagsByGitPathEndpoint(svc),
		BranchesByGitPathEndpoint: makeBranchesByGitPathEndpoint(svc),
		GetDockerfileEndpoint:     makeGetDockerfileEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()),           // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	emsNoNs := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Tags":              ems,
		"Branches":          ems,
		"TagsByGitPath":     emsNoNs,
		"BranchesByGitPath": emsNoNs,
		"GetDockerfile":     ems,
	}

	for _, m := range mw["Tags"] {
		eps.TagsEndpoint = m(eps.TagsEndpoint)
	}
	for _, m := range mw["Branches"] {
		eps.BranchesEndpoint = m(eps.BranchesEndpoint)
	}
	for _, m := range mw["TagsByGitPath"] {
		eps.TagsByGitPathEndpoint = m(eps.TagsByGitPathEndpoint)
	}
	for _, m := range mw["BranchesByGitPath"] {
		eps.BranchesByGitPathEndpoint = m(eps.BranchesByGitPathEndpoint)
	}
	for _, m := range mw["GetDockerfile"] {
		eps.GetDockerfileEndpoint = m(eps.GetDockerfileEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/git/tags/{namespace}/project/{name}", kithttp.NewServer(
		eps.TagsEndpoint,
		func(ctx context.Context, req *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/git/branches/{namespace}/project/{name}", kithttp.NewServer(
		eps.BranchesEndpoint,
		func(ctx context.Context, req *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")
	r.Handle("/git/tags", kithttp.NewServer(
		eps.TagsByGitPathEndpoint,
		decodeGitRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")
	r.Handle("/git/branches", kithttp.NewServer(
		eps.BranchesByGitPathEndpoint,
		decodeGitRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/git/dockerfile/{namespace}/project/{name}", kithttp.NewServer(
		eps.GetDockerfileEndpoint,
		decodeGetDockerfileRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeGetDockerfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getDockerfile{
		FileName: "Dockerfile",
	}
	return req, nil
}

func decodeGitRequest(_ context.Context, r *http.Request) (interface{}, error) {
	git := r.URL.Query().Get("git")
	if git == "" {
		return nil, encode.ErrBadRoute
	}
	return gitPath{git}, nil
}
