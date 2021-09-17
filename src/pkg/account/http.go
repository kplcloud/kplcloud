/**
 * @Time : 2021/9/17 3:29 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package account

import (
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"UserInfo": ems,
	})

	r := mux.NewRouter()

	r.Handle("/user-info", kithttp.NewServer(
		eps.UserInfoEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}
