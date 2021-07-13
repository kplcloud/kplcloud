/**
 * @Time : 6/15/21 5:09 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package nodes

import (
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/kplcloud/kplcloud/src/util/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"List": ems,
		//"Sync": ems,
	})

	r := mux.NewRouter()

	r.Handle("/sync", kithttp.NewServer(
		eps.SyncEndpoint,
		kithttp.NopRequestDecoder,
		encode.EncodeResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}
