/**
 * @Time : 2019/7/17 2:16 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package consul

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/config"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type endpoints struct {
	SyncEndpoint     endpoint.Endpoint
	DetailEndpoint   endpoint.Endpoint
	ListEndpoint     endpoint.Endpoint
	PostEndpoint     endpoint.Endpoint
	UpdateEndpoint   endpoint.Endpoint
	DeleteEndpoint   endpoint.Endpoint
	KVDetailEndpoint endpoint.Endpoint
	KVListEndpoint   endpoint.Endpoint
	KVPostEndpoint   endpoint.Endpoint
	KVDeleteEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger, cf *config.Config) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}
	eps := endpoints{
		SyncEndpoint:     makeSyncEndpoint(svc),
		DetailEndpoint:   makeDetailEndpoint(svc),
		ListEndpoint:     makeListEndpoint(svc),
		PostEndpoint:     makePostEndpoint(svc),
		UpdateEndpoint:   makeUpdateEndpoint(svc),
		DeleteEndpoint:   makeDeleteEndpoint(svc),
		KVDetailEndpoint: makeKVDetailEndpoint(svc),
		KVListEndpoint:   makeKVListEndpoint(svc),
		KVPostEndpoint:   makeKVPostEndpoint(svc),
		KVDeleteEndpoint: makeKVDeleteEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		consulVisitMiddleware(logger, cf),
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}
	noneNsEms := []endpoint.Middleware{
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory),
	}

	mw := map[string][]endpoint.Middleware{
		"Sync":     noneNsEms,
		"Detail":   ems,
		"List":     ems,
		"Post":     ems,
		"Update":   ems,
		"Delete":   ems,
		"KVDetail": ems,
		"KVList":   ems,
		"KVPost":   ems,
		"KVDelete": ems,
	}

	for _, m := range mw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range mw["Detail"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}
	for _, m := range mw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mw["KVDetail"] {
		eps.KVDetailEndpoint = m(eps.KVDetailEndpoint)
	}
	for _, m := range mw["KVList"] {
		eps.KVListEndpoint = m(eps.KVListEndpoint)
	}
	for _, m := range mw["KVPost"] {
		eps.KVPostEndpoint = m(eps.KVPostEndpoint)
	}
	for _, m := range mw["KVDelete"] {
		eps.KVDeleteEndpoint = m(eps.KVDeleteEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/consul/sync", kithttp.NewServer(
		eps.SyncEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/consul/{namespace}/one/{name}", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/consul/{namespace}/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/consul/{namespace}", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/consul/{namespace}/one/{name}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeUpdateRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/consul/{namespace}/one/{name}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/consul/kv/{namespace}/one/{name}", kithttp.NewServer(
		eps.KVDetailEndpoint,
		decodeKVDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/consul/kv/{namespace}/list/{name}", kithttp.NewServer(
		eps.KVListEndpoint,
		decodeKVDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/consul/kv/{namespace}/one/{name}", kithttp.NewServer(
		eps.KVPostEndpoint,
		decodeKVPostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/consul/kv/{namespace}/one/{name}", kithttp.NewServer(
		eps.KVDeleteEndpoint,
		decodeKVDeleteRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	return r
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	return detailRequest{
		Name:      name,
		Namespace: ns,
	}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	// get request body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req postRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	var ruleData = map[string]interface{}{
		"keyring": "read",
	}
	if ruleParams(req.Rules.Key) != nil {
		ruleData["key"] = ruleParams(req.Rules.Key)
	}
	if ruleParams(req.Rules.Event) != nil {
		ruleData["event"] = ruleParams(req.Rules.Event)
	}
	if ruleParams(req.Rules.Service) != nil {
		ruleData["service"] = ruleParams(req.Rules.Service)
	}
	if ruleParams(req.Rules.Query) != nil {
		ruleData["query"] = ruleParams(req.Rules.Query)
	}

	ruleByte, err := json.MarshalIndent(ruleData, "", "  ")
	if err != nil {
		return nil, encode.ErrBadRoute
	}
	req.Namespace = ns
	req.JsonRule = string(ruleByte)

	return req, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	// get request body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req postRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	var ruleData = map[string]interface{}{
		"keyring": "read",
	}
	if ruleParams(req.Rules.Key) != nil {
		ruleData["key"] = ruleParams(req.Rules.Key)
	}
	if ruleParams(req.Rules.Event) != nil {
		ruleData["event"] = ruleParams(req.Rules.Event)
	}
	if ruleParams(req.Rules.Service) != nil {
		ruleData["service"] = ruleParams(req.Rules.Service)
	}
	if ruleParams(req.Rules.Query) != nil {
		ruleData["query"] = ruleParams(req.Rules.Query)
	}

	ruleByte, err := json.MarshalIndent(ruleData, "", "  ")
	if err != nil {
		return nil, encode.ErrBadRoute
	}
	req.Namespace = ns
	req.Name = name
	req.JsonRule = string(ruleByte)

	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	pageStr := r.URL.Query().Get("p")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("name")
	p, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 15
	}
	return listRequest{
		detailRequest: detailRequest{Namespace: ns, Name: name},
		Page:          p,
		Limit:         limit,
	}, nil
}

func ruleParams(data []*rules) (params []map[string]interface{}) {
	for _, v := range data {
		if v.Name != "" {
			param := map[string]interface{}{
				v.Name: map[string]interface{}{
					"policy": v.Policy,
				},
			}
			params = append(params, param)
		}
	}
	return
}

func decodeKVDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	prefix := r.URL.Query().Get("prefix")
	if prefix == "" {
		prefix = ns + "." + name + "/"
	}
	return kvDetailRequest{
		Prefix: prefix,
		detailRequest: detailRequest{
			Name:      name,
			Namespace: ns,
		}}, nil
}
func decodeKVPostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req kvPostRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	req.Namespace = ns
	req.Name = name
	return req, nil
}

func decodeKVDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	prefix := r.URL.Query().Get("prefix")
	folder := r.URL.Query().Get("folder")
	var folderState bool
	if folder == "1" {
		folderState = true
	}
	if prefix == "" || !strings.Contains(prefix, ns+"."+name) {
		return nil, encode.ErrBadRoute
	}
	return kvDetailRequest{
		detailRequest{ns, name},
		prefix,
		folderState,
	}, nil
}
