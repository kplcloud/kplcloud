/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-06-24
 * Time: 21:38
 */
package group

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
	"github.com/kplcloud/kplcloud/src/util/helper"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	errBadRoute   = errors.New("bad route")
	errBadStrconv = errors.New("strconv error")
)

func MakeHandler(svc Service, logger log.Logger, groupRepository repository.GroupsRepository) http.Handler {

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	ownerOpts := opts
	ownerOpts = append(ownerOpts, kithttp.ServerBefore(middleware.GroupIdToContext()))

	commonEpsMap := map[string]endpoint.Endpoint{
		"memberByEmailLike": makeGetMemberByEmailLike(svc),
		"getAll":            makeGetAllEndPoint(svc),
		"nsProjectList":     makeNsProjectListEndPoint(svc),
		"nsCronjobList":     makeNsCronjobListEndPoint(svc),
		"ownerAddGroup":     makeOwnerAddGroupEndPoint(svc), // 应该放在这,创建是谁都可以创建的
		"userMyList":        makeUserGroupListEndPoint(svc),
		"userNsList":        makeUserNsListEndPoint(svc),
		"relDetail":         makeRelDetailEndPoint(svc),
		"nameExists":        makeGroupNameIsExistsEndPoint(svc),
		"displayNameExists": makeGroupDisplayNameIsExistsEndPoint(svc),
	}

	adminEpsMap := map[string]endpoint.Endpoint{
		"post":            makePostEndPoint(svc),
		"adminAdd":        makeAdminAddEndPoint(svc),
		"adminUpdate":     makeAdminUpdateGroup(svc),
		"adminDestroy":    makeAdminDestroyEndPoint(svc),
		"adminAddProject": makeAdminAddProjectEndPoint(svc),
		"adminAddCronjob": makeAdminAddCronjobEndPoint(svc),
		"adminAddMember":  makeAdminAddMemberEndPoint(svc),
		"adminDelMember":  makeAdminDelMemberEndPoint(svc),
		"adminDelProject": makeAdminDelProjectEndPoint(svc),
		"adminDelCronjob": makeAdminDelCronjobEndPoint(svc),
	}

	ownerEpsMap := map[string]endpoint.Endpoint{
		"ownerUpdateGroup": makeOwnerUpdateGroupEndPoint(svc),
		"ownerDelGroup":    makeOwnerDelGroupEndPoint(svc),
		"ownerAddMember":   makeOwnerAddMemberEndPoint(svc),
		"ownerDelMember":   makeOwnerDelMemberEndPoint(svc),
		"ownerAddProject":  makeOwnerAddProjectEndPoint(svc),
		"ownerDelProject":  makeOwnerDelProjectEndPoint(svc),
		"ownerAddCronjob":  makeOwnerAddCronjobEndPoint(svc),
		"ownerDelCronjob":  makeOwnerDelCronjobEndPoint(svc),
	}

	commonEms := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory),
	}

	ownerEms := []endpoint.Middleware{
		ownerDoGroupMiddleware(logger, groupRepository),
	}

	adminEms := []endpoint.Middleware{
		adminVisitMiddleware(logger),
	}

	ownerEms = helper.EmsArrMerge(ownerEms, commonEms)
	adminEms = helper.EmsArrMerge(adminEms, commonEms)

	ems := map[string][]endpoint.Middleware{
		"common": commonEms,
		"owner":  ownerEms,
		"admin":  adminEms,
	}

	eps := map[string]map[string]endpoint.Endpoint{
		"common": commonEpsMap,
		"owner":  ownerEpsMap,
		"admin":  adminEpsMap,
	}

	for n, m := range ems {
		for _, k := range m {
			for s, j := range eps[n] {
				eps[n][s] = k(j)
			}
		}
	}

	r := mux.NewRouter()

	r.Handle("/group", kithttp.NewServer(
		adminEpsMap["post"],
		decodePostRequest,
		encode.EncodeResponse,
		opts...)).Methods("POST")
	r.Handle("/group", kithttp.NewServer(
		commonEpsMap["getAll"],
		decodeGetAllRequest,
		encode.EncodeResponse,
		opts...)).Methods("GET")

	r.Handle("/group/admin-add", kithttp.NewServer(
		adminEpsMap["adminAdd"],
		decodePostRequest,
		encode.EncodeResponse,
		opts...)).Methods("POST")
	r.Handle("/group/{groupId:[0-9]+}/admin-update", kithttp.NewServer(
		adminEpsMap["adminUpdate"],
		decodeAdminUpdateGroupRequest,
		encode.EncodeResponse,
		opts...)).Methods("PUT")
	r.Handle("/group/member-like", kithttp.NewServer(
		commonEpsMap["memberByEmailLike"],
		decodeMemberEmailLikeRequest,
		encode.EncodeResponse,
		opts...)).Methods("GET")
	r.Handle("/group/{groupId:[0-9]+}/admin-delete", kithttp.NewServer(
		adminEpsMap["adminDestroy"],
		decodeAdminDestroy,
		encode.EncodeResponse,
		opts...)).Methods("DELETE")
	r.Handle("/group/namespace/{ns}/project", kithttp.NewServer(
		commonEpsMap["nsProjectList"],
		decodeNsProjectList,
		encode.EncodeResponse,
		opts...)).Methods("GET")
	r.Handle("/group/namespace/{ns}/cronjob", kithttp.NewServer(
		commonEpsMap["nsCronjobList"],
		decodeNsCronjobList,
		encode.EncodeResponse,
		opts...)).Methods("GET")
	r.Handle("/group/{groupId:[0-9]+}/admin-add-project/{projectId:[0-9]+}", kithttp.NewServer(
		adminEpsMap["adminAddProject"],
		decodeAdminDoProject,
		encode.EncodeResponse,
		opts...)).Methods("POST")
	r.Handle("/group/{groupId:[0-9]+}/admin-add-cronjob/{cronjobId:[0-9]+}", kithttp.NewServer(
		adminEpsMap["adminAddCronjob"],
		decodeAdminDoCronjob,
		encode.EncodeResponse,
		opts...)).Methods("POST")
	r.Handle("/group/{groupId:[0-9]+}/admin-add-member/{memberId:[0-9]+}", kithttp.NewServer(
		adminEpsMap["adminAddMember"],
		decodeAdminDoMember,
		encode.EncodeResponse,
		opts...)).Methods("POST")
	r.Handle("/group/{groupId:[0-9]+}/admin-del-member/{memberId:[0-9]+}", kithttp.NewServer(
		adminEpsMap["adminDelMember"],
		decodeAdminDoMember,
		encode.EncodeResponse,
		opts...)).Methods("DELETE")
	r.Handle("/group/{groupId:[0-9]+}/admin-del-project/{projectId:[0-9]+}", kithttp.NewServer(
		adminEpsMap["adminDelProject"],
		decodeAdminDoProject,
		encode.EncodeResponse,
		opts...)).Methods("DELETE")
	r.Handle("/group/{groupId:[0-9]+}/admin-del-cronjob/{cronjobId:[0-9]+}", kithttp.NewServer(
		adminEpsMap["adminDelCronjob"],
		decodeAdminDoCronjob,
		encode.EncodeResponse,
		opts...)).Methods("DELETE")

	r.Handle("/group/owner-add-group", kithttp.NewServer(
		commonEpsMap["ownerAddGroup"],
		decodeOwnerAddGroup,
		encode.EncodeResponse,
		opts...)).Methods("POST")

	r.Handle("/group/{groupId:[0-9]+}/owner-update-group", kithttp.NewServer(
		ownerEpsMap["ownerUpdateGroup"],
		decodeOwnerUpdateGroup,
		encode.EncodeResponse,
		ownerOpts...)).Methods("PUT")

	r.Handle("/group/{groupId:[0-9]+}/owner-add-member/{memberId:[0-9]+}", kithttp.NewServer(
		ownerEpsMap["ownerAddMember"],
		decodeOwnerDoMember,
		encode.EncodeResponse,
		ownerOpts...)).Methods("POST")

	r.Handle("/group/{groupId:[0-9]+}/owner-del-member/{memberId:[0-9]+}", kithttp.NewServer(
		ownerEpsMap["ownerDelMember"],
		decodeOwnerDoMember,
		encode.EncodeResponse,
		ownerOpts...)).Methods("DELETE")

	r.Handle("/group/{groupId:[0-9]+}/owner-add-project/{projectId:[0-9]+}", kithttp.NewServer(
		ownerEpsMap["ownerAddProject"],
		decodeOwnerDoProject,
		encode.EncodeResponse,
		ownerOpts...)).Methods("POST")

	r.Handle("/group/{groupId:[0-9]+}/owner-del-project/{projectId:[0-9]+}", kithttp.NewServer(
		ownerEpsMap["ownerDelProject"],
		decodeOwnerDoProject,
		encode.EncodeResponse,
		ownerOpts...)).Methods("DELETE")

	r.Handle("/group/{groupId:[0-9]+}/owner-add-cronjob/{cronjobId:[0-9]+}", kithttp.NewServer(
		ownerEpsMap["ownerAddCronjob"],
		decodeOwnerDoCronjob,
		encode.EncodeResponse,
		ownerOpts...)).Methods("POST")

	r.Handle("/group/{groupId:[0-9]+}/owner-del-cronjob/{cronjobId:[0-9]+}", kithttp.NewServer(
		ownerEpsMap["ownerDelCronjob"],
		decodeOwnerDoCronjob,
		encode.EncodeResponse,
		ownerOpts...)).Methods("DELETE")

	r.Handle("/group/{groupId:[0-9]+}/owner-del-group", kithttp.NewServer(
		ownerEpsMap["ownerDelGroup"],
		decodeOwnerDelGroup,
		encode.EncodeResponse,
		ownerOpts...)).Methods("DELETE")

	r.Handle("/group/user-my-list", kithttp.NewServer(
		commonEpsMap["userMyList"],
		decodeUserMyList,
		encode.EncodeResponse,
		opts...)).Methods("GET")

	r.Handle("/group/user-ns-list", kithttp.NewServer(
		commonEpsMap["userNsList"],
		func(i context.Context, request2 *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...)).Methods("GET")

	r.Handle("/group/{groupId:[0-9]+}/rel", kithttp.NewServer(
		commonEpsMap["relDetail"],
		decodeRelDetail,
		encode.EncodeResponse,
		opts...)).Methods("GET")

	r.Handle("/group/name/exists", kithttp.NewServer(
		commonEpsMap["nameExists"],
		decodeGroupNameIsExists,
		encode.EncodeResponse,
		opts...)).Methods("GET")

	r.Handle("/group/display_name/exists", kithttp.NewServer(
		commonEpsMap["displayNameExists"],
		decodeGroupDisplayNameIsExists,
		encode.EncodeResponse,
		opts...)).Methods("GET")
	return r
}

func decodeGroupNameIsExists(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	return isExists{
		Name: name,
	}, nil
}

func decodeGroupDisplayNameIsExists(_ context.Context, r *http.Request) (interface{}, error) {
	displayName := r.URL.Query().Get("display_name")
	return isExists{
		DisplayName: displayName,
	}, nil
}

func decodeRelDetail(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadStrconv
	}

	return relDetailRequest{
		GroupId: int64(groupIdInt),
	}, nil
}

func decodeUserMyList(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	ns := r.URL.Query().Get("ns")
	return userGroupListRequest{
		Name: name,
		Ns:   ns,
	}, nil
}

func decodeOwnerDelGroup(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadStrconv
	}

	return ownerDoGroup{
		Id: int64(groupIdInt),
	}, nil
}

func decodeOwnerDoMember(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadStrconv
	}

	memberId, ok := vars["memberId"]
	if !ok {
		return nil, errBadRoute
	}

	memberIdInt, err := strconv.Atoi(memberId)
	if err != nil {
		return nil, errBadStrconv
	}

	return ownerDoMember{
		MemberId: int64(memberIdInt),
		GroupId:  int64(groupIdInt),
	}, nil
}

func decodeOwnerDoCronjob(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadStrconv
	}

	cronjobId, ok := vars["cronjobId"]
	if !ok {
		return nil, errBadRoute
	}

	cronjobIdInt, err := strconv.Atoi(cronjobId)
	if err != nil {
		return nil, errBadStrconv
	}

	return ownerDoCronjob{
		CronjobId: int64(cronjobIdInt),
		GroupId:   int64(groupIdInt),
	}, nil
}

func decodeOwnerDoProject(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadStrconv
	}

	projectId, ok := vars["projectId"]
	if !ok {
		return nil, errBadRoute
	}

	projectIdInt, err := strconv.Atoi(projectId)
	if err != nil {
		return nil, errBadStrconv
	}

	return ownerDoProject{
		ProjectId: int64(projectIdInt),
		GroupId:   int64(groupIdInt),
	}, nil
}

func decodeOwnerUpdateGroup(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req ownerDoGroup

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadStrconv
	}
	req.Id = int64(groupIdInt)
	return req, nil
}

func decodeOwnerAddGroup(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req ownerDoGroup

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeAdminDoMember(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadStrconv
	}

	memberId, ok := vars["memberId"]
	if !ok {
		return nil, errBadRoute
	}

	memberIdInt, err := strconv.Atoi(memberId)
	if err != nil {
		return nil, errBadStrconv
	}

	return adminDoMemberRequest{
		MemberId: int64(memberIdInt),
		GroupId:  int64(groupIdInt),
	}, nil
}

func decodeAdminDoCronjob(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadStrconv
	}

	cronjobId, ok := vars["cronjobId"]
	if !ok {
		return nil, errBadRoute
	}

	cronjobIdInt, err := strconv.Atoi(cronjobId)
	if err != nil {
		return nil, errBadStrconv
	}

	return adminDoCronjobRequest{
		CronjobId: int64(cronjobIdInt),
		GroupId:   int64(groupIdInt),
	}, nil
}

func decodeAdminDoProject(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadStrconv
	}

	projectId, ok := vars["projectId"]
	if !ok {
		return nil, errBadRoute
	}

	projectIdInt, err := strconv.Atoi(projectId)
	if err != nil {
		return nil, errBadStrconv
	}

	return adminDoProjectRequest{
		ProjectId: int64(projectIdInt),
		GroupId:   int64(groupIdInt),
	}, nil
}

func decodeNsCronjobList(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["ns"]
	if !ok {
		return nil, errBadRoute
	}
	name := r.URL.Query().Get("name")

	return nsListRequest{
		Ns:   ns,
		Name: name,
	}, nil
}

func decodeNsProjectList(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["ns"]
	if !ok {
		return nil, errBadRoute
	}
	name := r.URL.Query().Get("name")

	return nsListRequest{
		Ns:   ns,
		Name: name,
	}, nil
}

func decodeGetAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	ns := r.URL.Query().Get("ns")
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
	return getAllRequest{
		Name:  name,
		Limit: limitInt,
		Ns:    ns,
		Page:  pageInt,
	}, nil
}

func decodeAdminUpdateGroupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadRoute
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req gRequest

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	req.GroupId = int64(groupIdInt)
	return req, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req gRequest

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeMemberEmailLikeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	email := r.URL.Query().Get("email")
	ns := r.URL.Query().Get("ns")

	if email == "" {
		email = "@"
	}

	if ns == "" {
		ns = "default"
	}

	return memberLikeRequest{
		Email: email,
		Ns:    ns,
	}, nil
}

func decodeAdminDestroy(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		return nil, errBadRoute
	}

	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, errBadRoute
	}

	return adminDestroyRequest{
		GroupId: int64(groupIdInt),
	}, nil
}
