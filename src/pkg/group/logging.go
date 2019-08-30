/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-06-25
 * Time: 11:30
 */
package group

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) Post(ctx context.Context, gr gRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "Post",
			"name", gr.Name,
			"displayName", gr.DisplayName,
			"namespace", gr.NameSpace,
			"memberId", gr.MemberId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, gr)
}

func (s *loggingService) GetAll(ctx context.Context, request getAllRequest) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "GetAll",
			"name", request.Name,
			"limit", request.Limit,
			"page", request.Page,
			"ns", request.Ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetAll(ctx, request)
}

func (s *loggingService) AdminUpdateGroup(ctx context.Context, groupId int64, gr gRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AdminUpdateGroup",
			"name", gr.Name,
			"name", gr.Name,
			"namespace", gr.NameSpace,
			"memberId", gr.MemberId,
			"groupId", groupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminUpdateGroup(ctx, gr.GroupId, gr)
}

func (s *loggingService) AdminAddGroup(ctx context.Context, gr gRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AdminAddGroup",
			"name", gr.Name,
			"name", gr.Name,
			"namespace", gr.NameSpace,
			"memberId", gr.MemberId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminAddGroup(ctx, gr)
}

func (s *loggingService) GetMemberByEmailLike(ctx context.Context, email string, ns string) (res []types.Member, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "GetMemberByEmailLike",
			"email", email,
			"namespace", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetMemberByEmailLike(ctx, email, ns)
}

func (s *loggingService) AdminDestroy(ctx context.Context, groupId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AdminDestroy",
			"groupId", groupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminDestroy(ctx, groupId)
}

func (s *loggingService) NamespaceProjectList(ctx context.Context, nq nsListRequest) (res []*types.Project, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "NamespaceProjectList",
			"name", nq.Name,
			"ns", nq.Ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.NamespaceProjectList(ctx, nq)
}

func (s *loggingService) NamespaceCronjobList(ctx context.Context, nq nsListRequest) (res []*types.Cronjob, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "NamespaceCronjobList",
			"name", nq.Name,
			"ns", nq.Ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.NamespaceCronjobList(ctx, nq)
}

func (s *loggingService) AdminAddProject(ctx context.Context, aq adminDoProjectRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AdminAddProject",
			"projectId", aq.ProjectId,
			"groupId", aq.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminAddProject(ctx, aq)
}

func (s *loggingService) AdminAddCronjob(ctx context.Context, ac adminDoCronjobRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AdminAddCronjob",
			"cronjobId", ac.CronjobId,
			"groupId", ac.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminAddCronjob(ctx, ac)
}

func (s *loggingService) AdminAddMember(ctx context.Context, am adminDoMemberRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AdminAddMember",
			"memeberId", am.MemberId,
			"groupId", am.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminAddMember(ctx, am)
}

func (s *loggingService) AdminDelMember(ctx context.Context, am adminDoMemberRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AdminDelMember",
			"memeberId", am.MemberId,
			"groupId", am.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminDelMember(ctx, am)
}

func (s *loggingService) AdminDelCronjob(ctx context.Context, ac adminDoCronjobRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AdminDelCronjob",
			"cronjobId", ac.CronjobId,
			"groupId", ac.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminDelCronjob(ctx, ac)
}

func (s *loggingService) AdminDelProject(ctx context.Context, aq adminDoProjectRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "AdminDelProject",
			"projectId", aq.ProjectId,
			"groupId", aq.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminDelProject(ctx, aq)
}

func (s *loggingService) OwnerAddGroup(ctx context.Context, oq ownerDoGroup) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "OwnerAddGroup",
			"name", oq.Name,
			"displayName", oq.DisplayName,
			"namespace", oq.Namespace,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.OwnerAddGroup(ctx, oq)
}

func (s *loggingService) OwnerUpdateGroup(ctx context.Context, oq ownerDoGroup) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "OwnerUpdateGroup",
			"name", oq.Name,
			"displayName", oq.DisplayName,
			"namespace", oq.Namespace,
			"groupId", oq.Id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.OwnerUpdateGroup(ctx, oq)
}

func (s *loggingService) OwnerAddMember(ctx context.Context, om ownerDoMember) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "OwnerAddMember",
			"memeberId", om.MemberId,
			"groupId", om.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.OwnerAddMember(ctx, om)
}

func (s *loggingService) OwnerDelMember(ctx context.Context, om ownerDoMember) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "OwnerDelMember",
			"memeberId", om.MemberId,
			"groupId", om.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.OwnerDelMember(ctx, om)
}

func (s *loggingService) OwnerAddProject(ctx context.Context, op ownerDoProject) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "OwnerAddProject",
			"projectId", op.ProjectId,
			"groupId", op.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.OwnerAddProject(ctx, op)
}

func (s *loggingService) OwnerDelProject(ctx context.Context, op ownerDoProject) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "OwnerDelProject",
			"projectId", op.ProjectId,
			"groupId", op.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.OwnerDelProject(ctx, op)
}

func (s *loggingService) OwnerAddCronjob(ctx context.Context, oc ownerDoCronjob) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "OwnerAddCronjob",
			"cronjobId", oc.CronjobId,
			"groupId", oc.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.OwnerAddCronjob(ctx, oc)
}

func (s *loggingService) OwnerDelCronjob(ctx context.Context, oc ownerDoCronjob) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "OwnerDelCronjob",
			"cronjobId", oc.CronjobId,
			"groupId", oc.GroupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.OwnerDelCronjob(ctx, oc)
}

func (s *loggingService) OwnerDelGroup(ctx context.Context, oq ownerDoGroup) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "OwnerDelGroup",
			"groupId", oq.Id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.OwnerDelGroup(ctx, oq)
}

func (s *loggingService) UserMyList(ctx context.Context, ugr userGroupListRequest) (res []*types.Groups, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "UserMyList",
			"namespace", ugr.Ns,
			"name", ugr.Name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UserMyList(ctx, ugr)
}

func (s *loggingService) NsMyList(ctx context.Context) (res []types.Namespace, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "NsMyList",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.NsMyList(ctx)
}

func (s *loggingService) RelDetail(ctx context.Context, groupId int64) (res *types.Groups, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "RelDetail",
			"groupId", groupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.RelDetail(ctx, groupId)
}

func (s *loggingService) GroupNameIsExists(ctx context.Context, name string) (notFound bool, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "GroupNameExists",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GroupNameExists(ctx, name)
}

func (s *loggingService) GroupDisplayNameIsExists(ctx context.Context, displayName string) (notFound bool, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(http.ContextKeyRequestURI),
			"method", "GroupDisplayNameExists",
			"name", displayName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GroupDisplayNameExists(ctx, displayName)
}
