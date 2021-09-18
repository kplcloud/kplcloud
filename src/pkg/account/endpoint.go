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
)

type Endpoints struct {
	UserInfoEndpoint endpoint.Endpoint
	MenusEndpoint    endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		UserInfoEndpoint: makeUserInfoEndpoint(s),
		MenusEndpoint:    makeMenusEndpoint(s),
	}

	for _, m := range dmw["UserInfo"] {
		eps.UserInfoEndpoint = m(eps.UserInfoEndpoint)
	}
	for _, m := range dmw["Menus"] {
		eps.MenusEndpoint = m(eps.MenusEndpoint)
	}
	return eps
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
