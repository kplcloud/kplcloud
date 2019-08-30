/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-06-24
 * Time: 18:06
 */
package group

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type gResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Err  error       `json:"error,omitempty"`
}

type gRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	NameSpace   string `json:"namespace"`
	MemberId    int64  `json:"memberId"`
	GroupId     int64  `json:"groupId"`
}

type getAllRequest struct {
	Name  string `json:"name,omitempty"`
	Limit int    `json:"limit,omitempty"`
	Page  int    `json:"page,omitempty"`
	Ns    string `json:"ns,omitempty"`
}

type memberLikeRequest struct {
	Email string `json:"email"`
	Ns    string `json:"ns,omitempty"`
}

type newGroup struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Namespace   newNs        `json:"namespace"`
	Members     []newMember  `json:"members"`
	Projects    []newProject `json:"projects"`
	Owner       newOwner     `json:"owner"`
}

type newMember struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type newProject struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	NameEn    string `json:"nameEn"`
	Language  string `json:"language"`
	Namespace string `json:"namespace"`
}

type newOwner struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type newNs struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type adminDestroyRequest struct {
	GroupId int64 `json:"groupId"`
}

type nsListRequest struct {
	Name string `json:"name"`
	Ns   string `json:"ns"`
}

type adminDoProjectRequest struct {
	ProjectId int64 `json:"projectId"`
	GroupId   int64 `json:"groupId"`
}

type adminDoCronjobRequest struct {
	CronjobId int64 `json:"cronjobId"`
	GroupId   int64 `json:"groupId"`
}

type adminDoMemberRequest struct {
	MemberId int64 `json:"memberId"`
	GroupId  int64 `json:"groupId"`
}

type ownerDoGroup struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Namespace   string `json:"namespace"`
}

type ownerDoMember struct {
	MemberId int64 `json:"memberId"`
	GroupId  int64 `json:"groupId"`
}

type ownerDoProject struct {
	ProjectId int64 `json:"projectId"`
	GroupId   int64 `json:"groupId"`
}

type ownerDoCronjob struct {
	CronjobId int64 `json:"cronjobId"`
	GroupId   int64 `json:"groupId"`
}

type userGroupListRequest struct {
	Ns   string `json:"ns"`
	Name string `json:"name"`
}

type relDetailRequest struct {
	GroupId int64 `json:"groupId"`
}

type isExists struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

func makeGroupNameIsExistsEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(isExists)
		res, err := s.GroupNameExists(ctx, req.Name)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeGroupDisplayNameIsExistsEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(isExists)
		res, err := s.GroupDisplayNameExists(ctx, req.DisplayName)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeRelDetailEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(relDetailRequest)
		if res, err := s.RelDetail(ctx, req.GroupId); err == nil {
			return encode.Response{Data: res}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeUserNsListEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if res, err := s.NsMyList(ctx); err == nil {
			return encode.Response{Data: res}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeUserGroupListEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(userGroupListRequest)
		if res, err := s.UserMyList(ctx, req); err == nil {
			return encode.Response{Data: res}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeOwnerDelGroupEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ownerDoGroup)
		if err := s.OwnerDelGroup(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeOwnerDelCronjobEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ownerDoCronjob)
		if err := s.OwnerDelCronjob(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeOwnerAddCronjobEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ownerDoCronjob)
		if err := s.OwnerAddCronjob(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeOwnerDelProjectEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ownerDoProject)
		if err := s.OwnerDelProject(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeOwnerAddProjectEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ownerDoProject)
		if err := s.OwnerAddProject(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeOwnerDelMemberEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ownerDoMember)
		if err := s.OwnerDelMember(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeOwnerAddMemberEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ownerDoMember)
		if err := s.OwnerAddMember(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeOwnerUpdateGroupEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ownerDoGroup)
		if err := s.OwnerUpdateGroup(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeOwnerAddGroupEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ownerDoGroup)
		if err := s.OwnerAddGroup(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeAdminDelCronjobEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(adminDoCronjobRequest)
		if err := s.AdminDelCronjob(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeAdminDelProjectEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(adminDoProjectRequest)
		if err := s.AdminDelProject(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeAdminDelMemberEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(adminDoMemberRequest)
		if err := s.AdminDelMember(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeAdminAddMemberEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(adminDoMemberRequest)
		if err := s.AdminAddMember(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeAdminAddCronjobEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(adminDoCronjobRequest)
		if err := s.AdminAddCronjob(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeAdminAddProjectEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(adminDoProjectRequest)
		if err := s.AdminAddProject(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeNsCronjobListEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(nsListRequest)
		if res, err := s.NamespaceCronjobList(ctx, req); err == nil {
			return encode.Response{Data: res}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeNsProjectListEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(nsListRequest)
		if res, err := s.NamespaceProjectList(ctx, req); err == nil {
			return encode.Response{Data: res}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeAdminDestroyEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(adminDestroyRequest)
		if err := s.AdminDestroy(ctx, req.GroupId); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makePostEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(gRequest)
		if err := s.Post(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeGetAllEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getAllRequest)
		if rs, err := s.GetAll(ctx, req); err == nil {
			return encode.Response{Data: rs}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeAdminAddEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(gRequest)
		if err := s.AdminAddGroup(ctx, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeAdminUpdateGroup(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(gRequest)
		if err := s.AdminUpdateGroup(ctx, req.GroupId, req); err == nil {
			return encode.Response{Data: nil}, err
		}
		return encode.Response{Err: err}, err
	}
}

func makeGetMemberByEmailLike(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(memberLikeRequest)
		if res, err := s.GetMemberByEmailLike(ctx, req.Email, req.Ns); err == nil {
			return encode.Response{Data: res}, err
		}
		return encode.Response{Err: err}, err
	}
}
