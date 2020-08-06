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
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
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

func MakeHandler(svc Service, logger log.Logger, repository repository.Repository, opts []kithttp.ServerOption, dmw []endpoint.Middleware) http.Handler {

	ems := []endpoint.Middleware{
		middleware.CronJobMiddleware(logger, repository.CronJob(), repository.Groups()),
	}

	eps := NewEndpoint(svc, map[string][]endpoint.Middleware{
		"addCronJob":  ems,
		"cronJobList": ems,

		"cronJobDetail":    append(ems, dmw...),
		"cronJobDel":       append(ems, dmw...),
		"cronJobAllDel":    append(ems, dmw...),
		"cronJobUpdate":    append(ems, dmw...),
		"cronJobLogUpdate": append(ems, dmw...),
		"Trigger":          append(ems, dmw...),
	})

	r := mux.NewRouter()

	r.Handle("/cronjob/{namespace}", kithttp.NewServer(
		eps.AddCronJobEndPoint,
		decodeAddCronJob,
		encode.EncodeResponse,
		opts...)).Methods(http.MethodPost)

	r.Handle("/cronjob/{namespace}", kithttp.NewServer(
		eps.CronJobListEndPoint,
		decodeCronJobList,
		encode.EncodeResponse,
		opts...)).Methods(http.MethodGet)
	r.Handle("/cronjob/{namespace}/{name}", kithttp.NewServer(
		eps.CronJobDelEndPoint,
		decodeCronJobDel,
		encode.EncodeResponse,
		opts...)).Methods(http.MethodDelete)
	r.Handle("/cronjob/job/{namespace}/del", kithttp.NewServer(
		eps.CronJobAllDelEndPoint,
		decodeCronJobAllDel,
		encode.EncodeResponse,
		opts...)).Methods(http.MethodDelete)
	r.Handle("/cronjob/{namespace}/{name}", kithttp.NewServer(
		eps.CronJobUpdateEndPoint,
		decodeCronJobUpdate,
		encode.EncodeResponse,
		opts...)).Methods(http.MethodPut)
	r.Handle("/cronjob/{namespace}/{name}", kithttp.NewServer(
		eps.CronJobDetailEndPoint,
		decodeCronJobDetail,
		encode.EncodeResponse,
		opts...)).Methods(http.MethodGet)

	r.Handle("/cronjob/log/{namespace}/{name}", kithttp.NewServer(
		eps.CronJobLogUpdateEndPoint,
		decodeCronJobLogUpdate,
		encode.EncodeResponse,
		opts...)).Methods(http.MethodPost)

	r.Handle("/cronjob/log/{namespace}/{name}/trigger", kithttp.NewServer(
		eps.TriggerEndpoint,
		decodeCronJobDetail,
		encode.EncodeResponse,
		opts...)).Methods(http.MethodPut)
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
