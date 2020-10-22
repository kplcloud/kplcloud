/**
 * @Time : 2019-06-28 10:33
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package deployment

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
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"io/ioutil"
	"net/http"
	"strings"
)

type endpoints struct {
	GetYamlEndpoint         endpoint.Endpoint
	CommandArgsEndpoint     endpoint.Endpoint
	ExpansionEndpoint       endpoint.Endpoint
	StretchEndpointEndpoint endpoint.Endpoint
	GetPvcEndpoint          endpoint.Endpoint
	BindPvcEndpoint         endpoint.Endpoint
	UnBindPvcEndpoint       endpoint.Endpoint
	AddPortEndpoint         endpoint.Endpoint
	DelPortEndpoint         endpoint.Endpoint
	LoggingEndpoint         endpoint.Endpoint
	ProbeEndpoint           endpoint.Endpoint
	MeshEndpoint            endpoint.Endpoint
	HostsEndpoint           endpoint.Endpoint
	VolumeConfigEndpoint    endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger, cf *config.Config, repository repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
		kithttp.ServerBefore(middleware.RequestIdToContext()),
		kithttp.ServerAfter(middleware.RequestIdToResponse()),
	}

	eps := endpoints{
		GetYamlEndpoint:         makeGetYamlEndpoint(svc),
		CommandArgsEndpoint:     makeCommandArgsEndpoint(svc),
		ExpansionEndpoint:       makeExpansionEndpoint(svc),
		StretchEndpointEndpoint: makeStretchEndpoint(svc),
		GetPvcEndpoint:          makeGetPvcEndpoint(svc),
		BindPvcEndpoint:         makeBindPvcEndpoint(svc),
		UnBindPvcEndpoint:       makeUnBindPvcEndpoint(svc),
		AddPortEndpoint:         makeAddPortEndpoint(svc),
		DelPortEndpoint:         makeDelPortEndpoint(svc),
		LoggingEndpoint:         makeLoggingEndpoint(svc),
		ProbeEndpoint:           makeProbeEndpoint(svc),
		MeshEndpoint:            makeMeshEndpoint(svc),
		HostsEndpoint:           makeHostsEndpoint(svc),
		VolumeConfigEndpoint:    makeVolumeConfigEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()),           // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	ems2 := ems
	ems2 = append(ems2, checkServiceMesh(cf))

	mw := map[string][]endpoint.Middleware{
		"GetYaml":      ems,
		"CommandArgs":  ems,
		"Expansion":    ems,
		"Stretch":      ems,
		"GetPvc":       ems,
		"BindPvc":      ems,
		"UnBindPvc":    ems,
		"AddPort":      ems,
		"DelPort":      ems,
		"Logging":      ems,
		"Probe":        ems,
		"Mesh":         ems2,
		"Hosts":        ems,
		"VolumeConfig": ems,
	}

	for _, m := range mw["GetYaml"] {
		eps.GetYamlEndpoint = m(eps.GetYamlEndpoint)
	}
	for _, m := range mw["CommandArgs"] {
		eps.CommandArgsEndpoint = m(eps.CommandArgsEndpoint)
	}
	for _, m := range mw["Expansion"] {
		eps.ExpansionEndpoint = m(eps.ExpansionEndpoint)
	}
	for _, m := range mw["Stretch"] {
		eps.StretchEndpointEndpoint = m(eps.StretchEndpointEndpoint)
	}
	for _, m := range mw["GetPvc"] {
		eps.GetPvcEndpoint = m(eps.GetPvcEndpoint)
	}
	for _, m := range mw["BindPvc"] {
		eps.BindPvcEndpoint = m(eps.BindPvcEndpoint)
	}
	for _, m := range mw["UnBindPvc"] {
		eps.UnBindPvcEndpoint = m(eps.UnBindPvcEndpoint)
	}
	for _, m := range mw["AddPort"] {
		eps.AddPortEndpoint = m(eps.AddPortEndpoint)
	}
	for _, m := range mw["DelPort"] {
		eps.DelPortEndpoint = m(eps.DelPortEndpoint)
	}
	for _, m := range mw["Logging"] {
		eps.LoggingEndpoint = m(eps.LoggingEndpoint)
	}
	for _, m := range mw["Probe"] {
		eps.ProbeEndpoint = m(eps.ProbeEndpoint)
	}
	for _, m := range mw["Mesh"] {
		eps.MeshEndpoint = m(eps.MeshEndpoint)
	}
	for _, m := range mw["Hosts"] {
		eps.HostsEndpoint = m(eps.HostsEndpoint)
	}
	for _, m := range mw["VolumeConfig"] {
		eps.VolumeConfigEndpoint = m(eps.VolumeConfigEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/deployment/{namespace}/yaml/{name}", kithttp.NewServer(
		eps.GetYamlEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/deployment/{namespace}/pvc/{name}", kithttp.NewServer(
		eps.GetPvcEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/deployment/{namespace}/command/{name}", kithttp.NewServer(
		eps.CommandArgsEndpoint,
		decodeCommandRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/deployment/{namespace}/expansion/{name}", kithttp.NewServer(
		eps.ExpansionEndpoint,
		decodeExpansionRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/deployment/{namespace}/scale/{name}", kithttp.NewServer(
		eps.StretchEndpointEndpoint,
		decodeStretchRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/deployment/{namespace}/pvc/{name}/bind", kithttp.NewServer(
		eps.BindPvcEndpoint,
		decodeBindPvcRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/deployment/{namespace}/pvc/{name}/unbind", kithttp.NewServer(
		eps.UnBindPvcEndpoint,
		decodeBindPvcRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/deployment/{namespace}/service/{name}/port", kithttp.NewServer(
		eps.AddPortEndpoint,
		decodeSvcPortRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/deployment/{namespace}/service/{name}/port", kithttp.NewServer(
		eps.DelPortEndpoint,
		decodeDelPortRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/deployment/{namespace}/logging/{name}", kithttp.NewServer(
		eps.LoggingEndpoint,
		decodeLoggingRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/deployment/{namespace}/probe/{name}", kithttp.NewServer(
		eps.ProbeEndpoint,
		decodeProbeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/deployment/{namespace}/mesh/{name}", kithttp.NewServer(
		eps.MeshEndpoint,
		decodeMeshRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/deployment/{namespace}/hosts/{name}", kithttp.NewServer(
		eps.HostsEndpoint,
		decodeHostsRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/deployment/{namespace}/volume-config/{name}", kithttp.NewServer(
		eps.VolumeConfigEndpoint,
		decodeVolumeConfigRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	return r
}

func decodeVolumeConfigRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req volumeConfigRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeHostsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req hostRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	req.Hosts = strings.Split(strings.Trim(strings.TrimSpace(req.Body), "\n"), "\n")

	return req, nil
}

func decodeMeshRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := meshRequest{
		Model: "normal",
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	req.Namespace = ns
	req.Name = name

	return req, nil
}

func decodeProbeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := probeRequest{
		InitialDelaySeconds: 10,
		TimeoutSeconds:      1,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	req.Namespace = ns
	req.Name = name

	return req, nil
}

func decodeLoggingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := loggingRequest{
		Paths:   []string{"/var/log/"},
		Pattern: "^20[0-9]{2}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}",
		Suffix:  "*.log",
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	req.Namespace = ns
	req.Name = name

	return req, nil
}

func decodeSvcPortRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req portRequest

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	req.Namespace = ns
	req.Name = name

	return req, nil
}

func decodeDelPortRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req delPortRequest

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	req.Namespace = ns
	req.Name = name

	return req, nil
}

func decodeBindPvcRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := bindPvcRequest{
		Path: "/opt/data/",
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	req.Namespace = ns
	req.Name = name

	return req, nil
}

func decodeStretchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := stretchRequest{
		Replicas: 1,
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	req.Namespace = ns
	req.Name = name

	return req, nil
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

func decodeExpansionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := expansionRequest{
		Cpu:       "0",
		MaxCpu:    "0",
		Memory:    "0",
		MaxMemory: "0",
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	req.Name = name
	req.Namespace = ns

	return req, nil
}

func decodeCommandRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	var req struct {
		Name      string      `json:"name"`
		Namespace string      `json:"namespace"`
		Args      interface{} `json:"args"`
		Command   interface{} `json:"command"`
	}
	var commands, args []string
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	switch req.Args.(type) {
	case string:
		args = strings.Split(strings.TrimSpace(req.Args.(string)), ",")
	case []interface{}:
		for _, v := range req.Args.([]interface{}) {
			val, ok := v.(string)
			if !ok || val == "" {
				continue
			}
			args = append(args, strings.TrimSpace(val))
		}
	}

	switch req.Command.(type) {
	case string:
		commands = strings.Split(strings.TrimSpace(req.Command.(string)), ",")
	case []interface{}:
		for _, v := range req.Command.([]interface{}) {
			val, ok := v.(string)
			if !ok || val == "" {
				continue
			}
			commands = append(commands, strings.TrimSpace(val))
		}
	}

	return commandArgsRequest{
		getRequest{
			Name:      name,
			Namespace: ns,
		},
		commands,
		args,
	}, nil
}
