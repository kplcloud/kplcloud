/**
 * @Time : 2019/7/5 11:03 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package configmap

import (
	"context"
	"encoding/json"
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
	"io/ioutil"
	"net/http"
	"strconv"
)

type endpoints struct {
	GetOneEndpoint              endpoint.Endpoint
	GetOnePullEndpoint          endpoint.Endpoint
	ListEndpoint                endpoint.Endpoint
	PostEndpoint                endpoint.Endpoint
	UpdateEndpoint              endpoint.Endpoint
	DeleteEndpoint              endpoint.Endpoint
	SyncEndpoint                endpoint.Endpoint
	CreateConfigMapEndpoint     endpoint.Endpoint //以下按单条处理
	GetConfigMapEndpoint        endpoint.Endpoint
	GetConfigMapDataEndpoint    endpoint.Endpoint
	CreateConfigMapDataEndpoint endpoint.Endpoint
	UpdateConfigMapDataEndpoint endpoint.Endpoint
	DeleteConfigMapDataEndpoint endpoint.Endpoint
	SyncConfigMapYamlEndpoint   endpoint.Endpoint
	UpdateConfigMapYamlEndpoint endpoint.Endpoint
	GetConfigEnvEndPoint        endpoint.Endpoint
	CreateConfigEnvEndPoint     endpoint.Endpoint
	UpdateConfigEnvEndPoint     endpoint.Endpoint
	DelConfigEnvEndPoint        endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger, store repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		GetOneEndpoint:              makeGetOneEndpoint(svc),
		GetOnePullEndpoint:          makeGetOnePullEndpoint(svc),
		ListEndpoint:                makeListEndpoint(svc),
		PostEndpoint:                makePostEndpoint(svc),
		UpdateEndpoint:              makeUpdateEndpoint(svc),
		DeleteEndpoint:              makeDeleteEndpoint(svc),
		SyncEndpoint:                makeSyncEndpoint(svc),
		CreateConfigMapEndpoint:     makeCreateConfigMapEndpoint(svc),
		GetConfigMapEndpoint:        makeGetConfigMapEndpoint(svc),
		GetConfigMapDataEndpoint:    makeGetConfigMapDataEndpoint(svc),
		CreateConfigMapDataEndpoint: makeCreateConfigMapDataEndpoint(svc),
		UpdateConfigMapDataEndpoint: makeUpdateConfigMapDataEndpoint(svc),
		DeleteConfigMapDataEndpoint: makeDeleteConfigMapDataEndpoint(svc),
		GetConfigEnvEndPoint:        makeGetConfigEnvEndPoint(svc),
		CreateConfigEnvEndPoint:     makeCreateConfigEnvEndPoint(svc),
		UpdateConfigEnvEndPoint:     makeUpdateConfigEnvEndPoint(svc),
		DelConfigEnvEndPoint:        makeDeleteConfigEnvEndPoint(svc),
	}
	ems := []endpoint.Middleware{
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}
	emsd := []endpoint.Middleware{
		//middleware.ProjectMiddleware(logger, store.Project(), store.Groups()),
		//middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}
	emse := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, store.Project(), store.Groups()),
		//middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}
	mw := map[string][]endpoint.Middleware{
		"GetOne":              ems,
		"GetOnePull":          ems,
		"List":                ems,
		"Post":                ems,
		"Update":              ems,
		"Delete":              ems,
		"Sync":                ems,
		"GetConfigEnv":        ems,
		"CreateConfigEnv":     ems,
		"UpdateConfigEnv":     ems,
		"DelConfigEnv":        ems,
		"CreateConfigMap":     emsd,
		"GetConfigMap":        emsd,
		"GetConfigMapData":    emsd,
		"CreateConfigMapData": emse,
		"UpdateConfigMapData": emse,
		"DeleteConfigMapData": emse,
	}

	for _, m := range mw["GetOne"] {
		eps.GetOneEndpoint = m(eps.GetOneEndpoint)
	}
	for _, m := range mw["GetOnePull"] {
		eps.GetOnePullEndpoint = m(eps.GetOnePullEndpoint)
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
	for _, m := range mw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}

	for _, m := range mw["CreateConfigMap"] {
		eps.CreateConfigMapEndpoint = m(eps.CreateConfigMapEndpoint)
	}
	for _, m := range mw["GetConfigMap"] {
		eps.GetConfigMapEndpoint = m(eps.GetConfigMapEndpoint)
	}
	for _, m := range mw["GetConfigMapData"] {
		eps.GetConfigMapDataEndpoint = m(eps.GetConfigMapDataEndpoint)
	}
	for _, m := range mw["CreateConfigMapData"] {
		eps.CreateConfigMapDataEndpoint = m(eps.CreateConfigMapDataEndpoint)
	}
	for _, m := range mw["UpdateConfigMapData"] {
		eps.UpdateConfigMapDataEndpoint = m(eps.UpdateConfigMapDataEndpoint)
	}
	for _, m := range mw["DeleteConfigMapData"] {
		eps.DeleteConfigMapDataEndpoint = m(eps.DeleteConfigMapDataEndpoint)
	}
	for _, m := range mw["DeleteConfigMapData"] {
		eps.DeleteConfigMapDataEndpoint = m(eps.DeleteConfigMapDataEndpoint)
	}
	for _, m := range mw["GetConfigEnv"] {
		eps.GetConfigEnvEndPoint = m(eps.GetConfigEnvEndPoint)
	}
	for _, m := range mw["CreateConfigEnv"] {
		eps.CreateConfigEnvEndPoint = m(eps.CreateConfigEnvEndPoint)
	}
	for _, m := range mw["UpdateConfigEnv"] {
		eps.UpdateConfigEnvEndPoint = m(eps.UpdateConfigEnvEndPoint)
	}
	for _, m := range mw["DelConfigEnv"] {
		eps.DelConfigEnvEndPoint = m(eps.DelConfigEnvEndPoint)
	}

	r := mux.NewRouter()
	r.Handle("/config/{namespace}/map/{name}", kithttp.NewServer(
		eps.GetOneEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/config/{namespace}/map/{name}/onePull", kithttp.NewServer(
		eps.GetOnePullEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/config/{namespace}", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/config/{namespace}", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/config/{namespace}/map/{name}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeUpdateRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/config/{namespace}/delete/{name}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/config/{namespace}/delete/{name}/data", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/config/{namespace}/sync", kithttp.NewServer(
		eps.SyncEndpoint,
		decodeSyncRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	//以下接口配置字典和配置数据单独处理
	r.Handle("/configmap/{namespace}/project/{name}", kithttp.NewServer(
		eps.CreateConfigMapEndpoint,
		decodeCreateConfigMapRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/configmap/{namespace}/project/{name}", kithttp.NewServer(
		eps.GetConfigMapEndpoint,
		decodeGetConfigMapRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/configmap/{namespace}/project/{name}/data", kithttp.NewServer(
		eps.CreateConfigMapDataEndpoint,
		decodeCreateConfigMapDataRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/configmap/{namespace}/project/{name}/data", kithttp.NewServer(
		eps.GetConfigMapDataEndpoint,
		decodeGetConfigMapDataRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/configmap/{namespace}/project/{name}/data/{id:[0-9]+}", kithttp.NewServer(
		eps.UpdateConfigMapDataEndpoint,
		decodeUpdateConfigMapDataRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/configmap/{namespace}/project/{name}/data/{id:[0-9]+}", kithttp.NewServer(
		eps.DeleteConfigMapDataEndpoint,
		decodeDeleteConfigMapDataRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/config/{namespace}/env/{name}", kithttp.NewServer(
		eps.GetConfigEnvEndPoint,
		decodeConfigEnvRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/config/{namespace}/env/{name}", kithttp.NewServer(
		eps.CreateConfigEnvEndPoint,
		decodeCreateConfEnvRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/config/{namespace}/env/{name}/{id:[0-9]+}", kithttp.NewServer(
		eps.UpdateConfigEnvEndPoint,
		decodeUpdateConfEnvRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/config/{namespace}/env/{name}/{id:[0-9]+}", kithttp.NewServer(
		eps.DelConfigEnvEndPoint,
		decodeDelConfEnvRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	return r
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	return getRequest{Name: name, Namespace: ns}, nil
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
		getRequest: getRequest{Namespace: ns, Name: name},
		Page:       p,
		Limit:      limit,
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
	req.Namespace = ns

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
	req.Name = name
	req.Namespace = ns

	return req, nil
}

func decodeSyncRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	return getRequest{Namespace: ns}, nil
}

func decodeCreateConfigMapRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req createConfigMapRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetConfigMapRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	return getRequest{Namespace: ns, Name: name}, nil
}

func decodeGetConfigMapDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	pageStr := r.URL.Query().Get("p")
	limitStr := r.URL.Query().Get("limit")
	p, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 15
	}
	return listRequest{
		getRequest: getRequest{Namespace: ns, Name: name},
		Page:       p,
		Limit:      limit,
	}, nil
}

func decodeCreateConfigMapDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req createConfigMapDataRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateConfigMapDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req configMapDataRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeDeleteConfigMapDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return configMapDataRequest{ConfigMapDataId: id}, nil
}

func decodeConfigEnvRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	pageStr := r.URL.Query().Get("p")
	limitStr := r.URL.Query().Get("limit")
	p, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 15
	}
	return listRequest{
		getRequest: getRequest{Namespace: ns, Name: name},
		Page:       p,
		Limit:      limit,
	}, nil
}

func decodeCreateConfEnvRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
	var req configEnvRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	req.Namespace = ns
	req.Name = name
	return req, nil
}

func decodeUpdateConfEnvRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	confEnvId, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	confEnvIdInt, err := strconv.Atoi(confEnvId)
	if err != nil {
		return nil, encode.ErrBadRoute
	}
	// get request body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req configEnvRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	req.Id = int64(confEnvIdInt)
	req.Namespace = ns
	req.Name = name
	return req, nil
}

func decodeDelConfEnvRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	confEnvId, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	confEnvIdInt, err := strconv.Atoi(confEnvId)
	if err != nil {
		return nil, encode.ErrBadRoute
	}

	return configEnvRequest{
		Id:        int64(confEnvIdInt),
		Name:      name,
		Namespace: ns,
	}, nil
}
