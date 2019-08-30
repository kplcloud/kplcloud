/**
 * Created by GoLand.
 * Email: xzghua@gmail.com
 * Date: 2019-07-09
 * Time: 15:01
 */
package cronjob

import (
	"context"
	"encoding/json"
	"errors"
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

var (
	errBadRoute   = errors.New("bad route")
	errBadStrconv = errors.New("strconv error")
)

type endpoints struct {
	AddCronJobEndPoint       endpoint.Endpoint
	CronJobListEndPoint      endpoint.Endpoint
	CronJobDelEndPoint       endpoint.Endpoint
	CronJobAllDelEndPoint    endpoint.Endpoint
	CronJobUpdateEndPoint    endpoint.Endpoint
	CronJobDetailEndPoint    endpoint.Endpoint
	CronJobLogUpdateEndPoint endpoint.Endpoint
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

	epsMap := endpoints{
		AddCronJobEndPoint:       makeAddCronJobEndPoint(svc),
		CronJobListEndPoint:      makeCronJobListEndPoint(svc),
		CronJobDelEndPoint:       makeCronJobDelEndPoint(svc),
		CronJobAllDelEndPoint:    makeCronJobAllDelEndPoint(svc),
		CronJobUpdateEndPoint:    makeCronJobUpdateEndPoint(svc),
		CronJobDetailEndPoint:    makeCronJobDetailEndPoint(svc),
		CronJobLogUpdateEndPoint: makeCronJobUpdateLogEndPoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.NamespaceMiddleware(logger), // 3
		middleware.CheckAuthMiddleware(logger),
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory),
	}

	emsc := []endpoint.Middleware{
		middleware.CronJobMiddleware(logger, repository.CronJob(), repository.Groups()),
		middleware.NamespaceMiddleware(logger), // 3
		middleware.CheckAuthMiddleware(logger),
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory),
	}

	mw := map[string][]endpoint.Middleware{
		"addCronJob":  ems,
		"cronJobList": ems,

		"cronJobDetail":    emsc,
		"cronJobDel":       emsc,
		"cronJobAllDel":    emsc,
		"cronJobUpdate":    emsc,
		"cronJobLogUpdate": emsc,
	}

	for _, m := range mw["addCronJob"] {
		epsMap.AddCronJobEndPoint = m(epsMap.AddCronJobEndPoint)
	}
	for _, m := range mw["cronJobList"] {
		epsMap.CronJobListEndPoint = m(epsMap.CronJobListEndPoint)
	}
	for _, m := range mw["cronJobDetail"] {
		epsMap.CronJobDetailEndPoint = m(epsMap.CronJobDetailEndPoint)
	}
	for _, m := range mw["cronJobDel"] {
		epsMap.CronJobDelEndPoint = m(epsMap.CronJobDelEndPoint)
	}
	for _, m := range mw["cronJobAllDel"] {
		epsMap.CronJobAllDelEndPoint = m(epsMap.CronJobAllDelEndPoint)
	}
	for _, m := range mw["cronJobUpdate"] {
		epsMap.CronJobUpdateEndPoint = m(epsMap.CronJobUpdateEndPoint)
	}
	for _, m := range mw["cronJobLogUpdate"] {
		epsMap.CronJobLogUpdateEndPoint = m(epsMap.CronJobLogUpdateEndPoint)
	}

	r := mux.NewRouter()

	r.Handle("/cronjob/{namespace}", kithttp.NewServer(
		epsMap.AddCronJobEndPoint,
		decodeAddCronJob,
		encode.EncodeResponse,
		opts...)).Methods("POST")

	r.Handle("/cronjob/{namespace}", kithttp.NewServer(
		epsMap.CronJobListEndPoint,
		decodeCronJobList,
		encode.EncodeResponse,
		opts...)).Methods("GET")
	r.Handle("/cronjob/{namespace}/{name}", kithttp.NewServer(
		epsMap.CronJobDelEndPoint,
		decodeCronJobDel,
		encode.EncodeResponse,
		opts...)).Methods("DELETE")
	r.Handle("/cronjob/job/{namespace}/del", kithttp.NewServer(
		epsMap.CronJobAllDelEndPoint,
		decodeCronJobAllDel,
		encode.EncodeResponse,
		opts...)).Methods("DELETE")
	r.Handle("/cronjob/{namespace}/{name}", kithttp.NewServer(
		epsMap.CronJobUpdateEndPoint,
		decodeCronJobUpdate,
		encode.EncodeResponse,
		opts...)).Methods("PUT")
	r.Handle("/cronjob/{namespace}/{name}", kithttp.NewServer(
		epsMap.CronJobDetailEndPoint,
		decodeCronJobDetail,
		encode.EncodeResponse,
		opts...)).Methods("GET")

	r.Handle("/cronjob/log/{namespace}/{name}", kithttp.NewServer(
		epsMap.CronJobLogUpdateEndPoint,
		decodeCronJobLogUpdate,
		encode.EncodeResponse,
		opts...)).Methods("POST")
	return r
}

func decodeCronJobLogUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, errBadRoute
	}

	name, ok := vars["name"]
	if !ok {
		return nil, errBadRoute
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req cronJobLogUpdate
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	req.Namespace = ns
	req.Name = name
	return req, nil
}

func decodeCronJobDetail(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, errBadRoute
	}

	name, ok := vars["name"]
	if !ok {
		return nil, errBadRoute
	}

	return cronJobDetail{
		Namespace: ns,
		Name:      name,
	}, nil
}

func decodeCronJobUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, errBadRoute
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req addCronJob

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	req.ParamName = name
	return req, nil
}

func decodeCronJobAllDel(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, errBadRoute
	}
	return cronJobAllDel{
		Namespace: ns,
	}, nil
}

func decodeCronJobDel(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, errBadRoute
	}

	name, ok := vars["name"]
	if !ok {
		return nil, errBadRoute
	}

	return cronJobDel{
		Namespace: ns,
		Name:      name,
	}, nil
}

func decodeCronJobList(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, errBadRoute
	}
	name := r.URL.Query().Get("name")
	group := r.URL.Query().Get("group")
	limit := r.URL.Query().Get("limit")

	var limitInt int
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}

		if limitInt == 0 {
			limitInt = 10
		}
	} else {
		limitInt = 10
	}

	page := r.URL.Query().Get("p")
	var pageInt int
	var err error
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			return nil, err
		}

		if pageInt == 0 {
			pageInt = 1
		}
	} else {
		pageInt = 1
	}
	return cronJobList{
		Name:      name,
		Namespace: ns,
		Group:     group,
		Page:      pageInt,
		Limit:     limitInt,
	}, nil
}

func decodeAddCronJob(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req addCronJob

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}
