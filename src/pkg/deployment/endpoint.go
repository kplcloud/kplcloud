/**
 * @Time : 2019-06-28 10:30
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package deployment

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type Probe string

const (
	ProbeReadiness Probe = "Readiness"
	ProbeLiveness  Probe = "Liveness"
)

func (c Probe) String() string {
	return string(c)
}

type getRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type commandArgsRequest struct {
	getRequest
	Command []string `json:"command"`
	Args    []string `json:"args"`
}

type expansionRequest struct {
	getRequest
	Cpu       string `json:"cpu"`
	MaxCpu    string `json:"maxCpu"`
	Memory    string `json:"memory"`
	MaxMemory string `json:"maxMemory"`
}

type stretchRequest struct {
	getRequest
	Replicas int `json:"replicas"`
}

type bindPvcRequest struct {
	getRequest
	Path      string `json:"path"`
	ClaimName string `json:"claim_name"`
}

type ports struct {
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
	Name     string `json:"name"`
}
type portRequest struct {
	getRequest
	Ports []ports
}

type delPortRequest struct {
	getRequest
	Port     int32  `json:"port"`
	PortName string `json:"port_name"`
}

type loggingRequest struct {
	getRequest
	Paths   []string `json:"paths" yaml:"paths"`
	Pattern string   `json:"pattern" yaml:"pattern"`
	Suffix  string   `json:"suffix"`
}

type probeRequest struct {
	getRequest
	Probe               []string `json:"probe"` // Readiness 就绪探针,Liveness 存活探针
	Port                int32    `json:"port"`
	InitialDelaySeconds int32    `json:"initial_delay_seconds"` // 容器启动后第一次执行探测是需要等待多少秒。
	TimeoutSeconds      int32    `json:"timeout_seconds"`       // 探测超时时间。默认1秒，最小1秒。
	PeriodSeconds       int32    `json:"period_seconds"`        // 执行探测的频率。默认是10秒，最小1秒。
	SuccessThreshold    int32    `json:"success_threshold"`     // 探测失败后，最少连续探测成功多少次才被认定为成功。默认是 1。对于 liveness 必须是 1。最小值是 1。
	FailureThreshold    int32    `json:"failure_threshold"`     // 探测成功后，最少连续探测失败多少次才被认定为失败。默认是 3。最小值是 1。
	Path                string   `json:"path"`                  // 如果是http 需要填path
}

type meshRequest struct {
	getRequest
	Model string `json:"model"`
}

type hostRequest struct {
	Body  string   `json:"body"`
	Hosts []string `json:"hosts"`
}

type volumeConfigRequest struct {
	MountPath string `json:"mount_path"`
	SubPath   string `json:"sub_path"`
}

func makeGetYamlEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetYaml(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeCommandArgsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(commandArgsRequest)
		err = s.CommandArgs(ctx, req.Command, req.Args)
		return encode.Response{Err: err}, err
	}
}

func makeExpansionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(expansionRequest)
		err = s.Expansion(ctx, req.Cpu, req.MaxCpu, req.Memory, req.MaxMemory)
		return encode.Response{Err: err}, err
	}
}

func makeStretchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(stretchRequest)
		err = s.Stretch(ctx, req.Replicas)
		return encode.Response{Err: err}, err
	}
}

func makeGetPvcEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		rs, err := s.GetPvc(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err, Data: rs}, err
	}
}

func makeBindPvcEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bindPvcRequest)
		err = s.BindPvc(ctx, req.Namespace, req.Name, req.Path, req.ClaimName)
		return encode.Response{Err: err}, err
	}
}

func makeUnBindPvcEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bindPvcRequest)
		err = s.UnBindPvc(ctx, req.Namespace, req.Name, req.ClaimName)
		return encode.Response{Err: err}, err
	}
}

func makeAddPortEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(portRequest)
		err = s.AddPort(ctx, req.Namespace, req.Name, req)
		return encode.Response{Err: err}, err
	}
}

func makeDelPortEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(delPortRequest)
		err = s.DelPort(ctx, req.Namespace, req.Name, req.PortName, req.Port)
		return encode.Response{Err: err}, err
	}
}

func makeLoggingEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loggingRequest)
		err = s.Logging(ctx, req.Namespace, req.Name, req.Pattern, req.Suffix, req.Paths)
		return encode.Response{Err: err}, err
	}
}

func makeProbeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(probeRequest)
		err = s.Probe(ctx, req.Namespace, req.Name, req)

		return encode.Response{Err: err}, err
	}
}

func makeMeshEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(meshRequest)
		err = s.Mesh(ctx, req.Namespace, req.Name, req.Model)

		return encode.Response{Err: err}, err
	}
}

func makeHostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(hostRequest)
		err = s.Hosts(ctx, req.Hosts)
		return encode.Response{Err: err}, err
	}
}

func makeVolumeConfigEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(volumeConfigRequest)
		err = s.VolumeConfig(ctx, req.MountPath, req.SubPath)
		return encode.Response{Err: err}, err
	}
}
