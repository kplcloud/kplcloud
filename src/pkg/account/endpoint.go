/**
 * @Time : 2021/9/17 3:08 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package account

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/middleware"
)

type (
	userInfoResult struct {
		Username    string   `json:"username"`
		Avatar      string   `json:"avatar"`
		Roles       []string `json:"roles"`
		Permissions []string `json:"permissions"`
		Clusters    []string `json:"clusters"`
	}

	userMenuResult struct {
		Id       int64            `json:"-"`
		ParentId int64            `json:"-"`
		Icon     string           `json:"icon"`
		Key      string           `json:"-"`
		Text     string           `json:"text"`
		Link     string           `json:"link"`
		Alias    string           `json:"-"`
		Items    []userMenuResult `json:"items"`
	}

	nsResult struct {
		Name  string `json:"name"`
		Alias string `json:"alias"`
	}
)

type Endpoints struct {
	UserInfoEndpoint   endpoint.Endpoint
	MenusEndpoint      endpoint.Endpoint
	LogoutEndpoint     endpoint.Endpoint
	NamespacesEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		UserInfoEndpoint:   makeUserInfoEndpoint(s),
		MenusEndpoint:      makeMenusEndpoint(s),
		LogoutEndpoint:     makeLogoutEndpoint(s),
		NamespacesEndpoint: makeNamespacesEndpoint(s),
	}

	for _, m := range dmw["UserInfo"] {
		eps.UserInfoEndpoint = m(eps.UserInfoEndpoint)
	}
	for _, m := range dmw["Menus"] {
		eps.MenusEndpoint = m(eps.MenusEndpoint)
	}
	for _, m := range dmw["Logout"] {
		eps.LogoutEndpoint = m(eps.LogoutEndpoint)
	}
	for _, m := range dmw["Namespaces"] {
		eps.NamespacesEndpoint = m(eps.NamespacesEndpoint)
	}
	return eps
}

func makeLogoutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.ContextUserId).(int64)
		if !ok {
			return nil, encode.ErrAccountNotLogin.Error()
		}
		err = s.Logout(ctx, userId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeMenusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.ContextUserId).(int64)
		if !ok {
			return nil, encode.ErrAccountNotLogin.Error()
		}
		info, err := s.Menus(ctx, userId)
		return encode.Response{
			Data:  info,
			Error: err,
		}, err
	}
}

func makeUserInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.ContextUserId).(int64)
		if !ok {
			return nil, encode.ErrAccountNotLogin.Error()
		}
		info, err := s.UserInfo(ctx, userId)
		return encode.Response{
			Data:  info,
			Error: err,
		}, err
	}
}

func makeNamespacesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.ContextUserId).(int64)
		if !ok {
			return nil, encode.ErrAccountNotLogin.Error()
		}
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		info, err := s.Namespaces(ctx, userId, clusterId)
		return encode.Response{
			Data:  info,
			Error: err,
		}, err
	}
}
