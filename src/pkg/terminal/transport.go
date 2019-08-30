/**
 * @Time : 2019-06-27 18:14
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package terminal

import (
	"bytes"
	"context"
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
	"html/template"
	"net/http"
	"path"
	"strings"
)

type endpoints struct {
	//AttachEndpoint endpoint.Endpoint
	IndexEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger kitlog.Logger, repository repository.Repository) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CookieToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
		//kithttp.ServerBefore(kithttp.SetRequestHeader("Authorization", "")),
	}

	eps := endpoints{
		//AttachEndpoint: makeAttachEndpoint(svc),
		IndexEndpoint: makeIndexEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()), // 4
		middleware.NamespaceMiddleware(logger),                                          // 3
		middleware.CheckAuthMiddleware(logger),                                          // 2
		//kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
		NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory),
	}

	mw := map[string][]endpoint.Middleware{
		"Index": ems,
	}

	for _, m := range mw["Index"] {
		eps.IndexEndpoint = m(eps.IndexEndpoint)
	}

	r.Handle("/terminal/{namespace}/index/{name}/pod/{podName}/container/{container}",
		kithttp.NewServer(
			eps.IndexEndpoint,
			decodeIndexRequest,
			encodeIndexResponse, opts...)).Methods("GET")

	return r
}

func decodeIndexRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	podName, ok := vars["podName"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	container, ok := vars["container"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	var token string
	if c, err := r.Cookie("Authorization"); err == nil {
		token = c.Value
	}

	return indexRequest{
		PodName:   podName,
		Container: container,
		Token:     strings.TrimSpace(token),
	}, nil
}

func encodeIndexResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encode.EncodeError(ctx, e.error(), w)
		return nil
	}

	resp := response.(indexResponse)

	name := path.Base(resp.Data.TemplateFile)
	t, err := template.New(name).ParseFiles(resp.Data.TemplateFile)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, resp.Data); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())

	return err
}
