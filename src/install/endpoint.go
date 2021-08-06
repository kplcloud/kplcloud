/**
 * @Time : 7/21/21 2:26 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package install

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"mime/multipart"
)

type (
	initDbRequest struct {
		Drive    string `json:"drive" valid:"required"`
		Host     string `json:"host" valid:"required"`
		Port     int    `json:"port" valid:"required"`
		User     string `json:"user" valid:"required"`
		Password string `json:"password" valid:"required"`
		Database string `json:"database" valid:"required"`
	}

	initPlatformRequest struct {
		AppName       string `json:"appName" valid:"required"`
		AdminName     string `json:"adminName" valid:"required"`
		AdminPassword string `json:"adminPassword" valid:"required"`
		AppKey        string `json:"appKey"`
		Domain        string `json:"domain" valid:"required"`
		DomainSuffix  string `json:"domainSuffix" valid:"required"`
		LogPath       string `json:"logPath"`
		LogLevel      string `json:"logLevel"`
		UploadPath    string `json:"uploadPath" valid:"required"`
		Debug         bool   `json:"debug"`
	}

	initLogoRequest struct {
		Files []*multipart.FileHeader
	}

	initCorsRequest struct {
		Allow   bool     `json:"allow"`
		Methods []string `json:"methods" valid:"required"`
		Headers string   `json:"headers" valid:"required"`
		Origin  string   `json:"origin"  valid:"required"`
		Method  string   `json:"-"`
	}

	initRedisRequest struct {
		Hosts    string `json:"hosts" valid:"required"`
		Password string `json:"password"`
		Database int    `json:"database"`
		Prefix   string `json:"prefix"`
	}
)

type Endpoints struct {
	InitDbEndpoint       endpoint.Endpoint
	InitPlatformEndpoint endpoint.Endpoint
	InitLogoEndpoint     endpoint.Endpoint
	InitCorsEndpoint     endpoint.Endpoint
	InitRedisEndpoint    endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		InitDbEndpoint:       makeInitDbEndpoint(s),
		InitPlatformEndpoint: makeInitPlatformEndpoint(s),
		InitLogoEndpoint:     makeInitLogoEndpoint(s),
		InitCorsEndpoint:     makeInitCorsEndpoint(s),
		InitRedisEndpoint:    makeInitRedisEndpoint(s),
	}

	for _, m := range dmw["InitDb"] {
		eps.InitDbEndpoint = m(eps.InitDbEndpoint)
	}
	for _, m := range dmw["InitPlatform"] {
		eps.InitPlatformEndpoint = m(eps.InitPlatformEndpoint)
	}
	for _, m := range dmw["InitLogo"] {
		eps.InitLogoEndpoint = m(eps.InitLogoEndpoint)
	}
	for _, m := range dmw["InitCors"] {
		eps.InitCorsEndpoint = m(eps.InitCorsEndpoint)
	}
	for _, m := range dmw["InitRedis"] {
		eps.InitRedisEndpoint = m(eps.InitRedisEndpoint)
	}

	return eps
}

func makeInitRedisEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(initRedisRequest)
		err = s.InitRedis(ctx, req.Hosts, req.Password, req.Database, req.Prefix)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeInitCorsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(initCorsRequest)
		err = s.InitCors(ctx, req.Allow, req.Origin, req.Method, req.Headers)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeInitLogoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(initLogoRequest)
		err = s.InitLogo(ctx, req.Files[0])
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeInitPlatformEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(initPlatformRequest)
		err = s.InitPlatform(ctx, req.AppName, req.AdminName, req.AdminPassword, req.AppKey, req.Domain, req.DomainSuffix, req.LogPath, req.LogLevel, req.UploadPath, req.Debug)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeInitDbEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(initDbRequest)
		err = s.InitDb(ctx, req.Drive, req.Host, req.Port, req.User, req.Password, req.Database)
		return encode.Response{
			Error: err,
		}, err
	}
}
