/**
 * @Time : 2021/11/3 6:25 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package audits

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"List": ems,
	})

	r := mux.NewRouter()

	r.Handle("/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req.page = page
	req.pageSize = pageSize
	req.query = r.URL.Query().Get("query")
	req.status = r.URL.Query().Get("status")

	return req, nil
}
