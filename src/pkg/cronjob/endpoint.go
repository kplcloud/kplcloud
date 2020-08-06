/**
 * Created by GoLand.
 * Email: xzghua@gmail.com
 * Date: 2019-07-09
 * Time: 14:59
 */
package cronjob

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"gopkg.in/guregu/null.v3"
	"k8s.io/api/batch/v1beta1"
	"time"
)

type Endpoints struct {
	AddCronJobEndPoint       endpoint.Endpoint
	CronJobListEndPoint      endpoint.Endpoint
	CronJobDelEndPoint       endpoint.Endpoint
	CronJobAllDelEndPoint    endpoint.Endpoint
	CronJobUpdateEndPoint    endpoint.Endpoint
	CronJobDetailEndPoint    endpoint.Endpoint
	CronJobLogUpdateEndPoint endpoint.Endpoint
	TriggerEndpoint          endpoint.Endpoint
}

type DetailReturnData struct {
	Name          string                   `json:"name"`
	Namespace     string                   `json:"namespace"`
	Schedule      string                   `json:"schedule"`
	GitType       string                   `json:"git_type"`
	GitPath       string                   `json:"git_path"`
	Image         string                   `json:"image"`
	Suspend       bool                     `json:"suspend"`
	Active        int64                    `json:"active"`
	LastSchedule  null.Time                `json:"last_schedule"`
	ConfMapName   string                   `json:"conf_map_name"`
	Args          []interface{}            `json:"args"`
	Command       string                   `json:"Command"`
	LogPath       string                   `json:"log_path"`
	AddType       string                   `json:"add_type"`
	CronjobInfo   map[string]interface{}   `json:"cronjob_info"`
	CronjobPods   []map[string]interface{} `json:"cronjob_pods"`
	CronjobEvents []map[string]interface{} `json:"cronjob_events"`
	CronjobYaml   *v1beta1.CronJob         `json:"cronjob_yaml"`
}

type JenkinsCronjobData struct {
	Name      string
	Namespace string
	BuildId   int64
	BuildTime time.Time
}

type addCronJob struct {
	Name      string   `json:"name"`
	Args      []string `json:"args"`
	GitType   string   `json:"git_type"`
	GitPath   string   `json:"git_path"`
	Image     string   `json:"image"`
	Namespace string   `json:"namespace"`
	Schedule  string   `json:"schedule"`
	ConfMap   string   `json:"conf_map"`
	LogPath   string   `json:"log_path"`
	AddType   string   `json:"add_type"`
	ParamName string   `json:"param_name"`
}

type cronJobList struct {
	Name      string `json:"name"`
	Group     string `json:"group"`
	Namespace string `json:"namespace"`
	Limit     int    `json:"limit,omitempty"`
	Page      int    `json:"page,omitempty"`
}

type cronJobDel struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type cronJobDetail struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type cronJobAllDel struct {
	Namespace string `json:"namespace"`
}

type cronJobLogUpdate struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	LogPath   string `json:"log_path"`
}

var code int

func NewEndpoint(svc Service, mw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		AddCronJobEndPoint:       makeAddCronJobEndPoint(svc),
		CronJobListEndPoint:      makeCronJobListEndPoint(svc),
		CronJobDelEndPoint:       makeCronJobDelEndPoint(svc),
		CronJobAllDelEndPoint:    makeCronJobAllDelEndPoint(svc),
		CronJobUpdateEndPoint:    makeCronJobUpdateEndPoint(svc),
		CronJobDetailEndPoint:    makeCronJobDetailEndPoint(svc),
		CronJobLogUpdateEndPoint: makeCronJobUpdateLogEndPoint(svc),
		TriggerEndpoint:          makeTriggerEndpoint(svc),
	}

	for _, m := range mw["addCronJob"] {
		eps.AddCronJobEndPoint = m(eps.AddCronJobEndPoint)
	}
	for _, m := range mw["cronJobList"] {
		eps.CronJobListEndPoint = m(eps.CronJobListEndPoint)
	}
	for _, m := range mw["cronJobDetail"] {
		eps.CronJobDetailEndPoint = m(eps.CronJobDetailEndPoint)
	}
	for _, m := range mw["cronJobDel"] {
		eps.CronJobDelEndPoint = m(eps.CronJobDelEndPoint)
	}
	for _, m := range mw["cronJobAllDel"] {
		eps.CronJobAllDelEndPoint = m(eps.CronJobAllDelEndPoint)
	}
	for _, m := range mw["cronJobUpdate"] {
		eps.CronJobUpdateEndPoint = m(eps.CronJobUpdateEndPoint)
	}
	for _, m := range mw["cronJobLogUpdate"] {
		eps.CronJobLogUpdateEndPoint = m(eps.CronJobLogUpdateEndPoint)
	}
	for _, m := range mw["Trigger"] {
		eps.TriggerEndpoint = m(eps.TriggerEndpoint)
	}

	return eps
}

func makeTriggerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(cronJobDetail)
		err = s.Trigger(ctx, req.Name, req.Namespace)
		return encode.Response{Err: err}, err
	}
}

func makeCronJobUpdateLogEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(cronJobLogUpdate)
		err = s.UpdateLog(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeCronJobDetailEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(cronJobDetail)
		res, err := s.Detail(ctx, req.Name, req.Namespace)
		if err != nil {
			code = -1
		}
		return encode.Response{Code: code, Err: err, Data: res}, err
	}
}

func makeCronJobUpdateEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addCronJob)
		err = s.Put(ctx, req.ParamName, req)
		return encode.Response{Err: err}, err
	}
}

func makeCronJobAllDelEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(cronJobAllDel)
		err = s.DeleteJobAll(ctx, req.Namespace)
		return encode.Response{Err: err, Data: nil}, err
	}
}

func makeCronJobDelEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(cronJobDel)
		err = s.Delete(ctx, req.Name, req.Namespace)
		if err != nil {
			code = -1
		}
		return encode.Response{Code: code, Err: err, Data: nil}, err
	}
}

func makeCronJobListEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(cronJobList)
		res, err := s.List(ctx, req.Name, req.Namespace, req.Group, req.Page, req.Limit)
		if err != nil {
			code = -1
		}
		return encode.Response{Code: code, Err: err, Data: res}, err
	}
}

func makeAddCronJobEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addCronJob)
		err = s.AddCronJob(ctx, req)
		if err != nil {
			code = -1
		}
		return encode.Response{Code: code, Err: err}, err
	}
}
