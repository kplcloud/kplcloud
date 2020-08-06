/**
 * @Time : 2020/6/23 6:13 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package virtualservice

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/kplcloud/kplcloud/src/util/encode"
)

func MakeHTTPHandler(s Service, logger log.Logger, opts []kithttp.ServerOption, dmw []endpoint.Middleware) http.Handler {
	ems := []endpoint.Middleware{}

	//s = NewLoggingServer(logger, s)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Get": append(ems, dmw...),
	})

	r := mux.NewRouter()

	r.Handle("/virtualservice/{namespace}/{name}", kithttp.NewServer(
		eps.GetEndpoint,
		kithttp.NopRequestDecoder,
		encode.EncodeResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}
