/**
 * @Time: 2021/12/5 21:55
 * @Author: solacowa@gmail.com
 * @File: http
 * @Software: GoLand
 */

package pod

import (
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"net/http"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opt []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware

	ems = append(ems, dmw...)
	//
	//eps := NewEndpoint(s, map[string][]endpoint.Middleware{
	//	"Sync":   ems,
	//	"SyncPv": ems,
	//	"Create": ems,
	//	"Info":   ems,
	//	"List":   ems,
	//	"Update": ems,
	//})
	opts := []kithttp.ServerOption{}

	opts = append(opts, opt...)

	r := mux.NewRouter()

	r.Handle("/ws/pods/console/exec/", sockjs.NewHandler("/ws/pods/console/exec", sockjs.DefaultOptions, func(session sockjs.Session) {
		//termilanSvc.HandleTerminalSession(session)
	})).Methods(http.MethodGet)

	//r.Handle("/{cluster}/list", kithttp.NewServer(
	//	eps.ListEndpoint,
	//	decodeListRequest,
	//	encode.JsonResponse,
	//	opts...,
	//)).Methods(http.MethodGet)

	return r
}
