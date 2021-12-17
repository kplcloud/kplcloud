/**
 * @Time : 2021/12/16 4:37 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package sysgroup

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

// Service 组管理模块
// 每个空间下都可以创建很多组
// 每个组可以添加app,定时任务,有状态服务
// 组里可以添加成员及权限 读/操作
// 全部都是多对多的关系
// 两种方案
//   1. 组的类型，只读或可操作 组里人跟随组权限
//   2. 组没有类型，组里的人设置读写权限
// 我决定采用1 方案，如果用户在多个组里并且组里都有项目的项目则取最大组权限
// 创建空间时默认给创建一个写组，添加新成员时必须给分配一个组
// TODO: 如何解决超管查询及更新处理等问题，能否在中间件完成？
type Service interface {
	// Create 创建组
	// userId 组管理员
	Create(ctx context.Context, sysUserId, clusterId int64, namespace, groupName, groupAlias, remark string, userId int64, onlyRead bool) (err error)
	// Update 更新组信息
	Update(ctx context.Context, sysUserId, clusterId int64, namespace, groupName, groupAlias, remark string, userId int64, onlyRead bool) (err error)
	// Delete 删除组
	// 删除组需要把他们关联关系全部删除
	// 只有组管理或超管可以删除
	Delete(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string) (err error)
	// List 组列表
	// groupName 可以为空
	List(ctx context.Context, clusterId int64, groupIds []int64, namespace, groupName string, page, pageSize int) (res []result, total int, err error)
	// AddUser 给组添加成员
	AddUser(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, userIds []int64) (err error)
	// AddApp 给组添加app
	AddApp(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error)
	// AddCronJob 给组添加定时任务
	AddCronJob(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error)
	// AddStatefulSet 给组添加有状态应用
	AddStatefulSet(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error)
	// DeleteApp 删除组里的应用
	DeleteApp(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error)
	// DeleteCronJob 删除组里的定时任务
	DeleteCronJob(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error)
	// DeleteStatefulSet 删除组里的有状态应用
	DeleteStatefulSet(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
}

func (s *service) Create(ctx context.Context, sysUserId, clusterId int64, namespace, groupName, groupAlias, remark string, userId int64, onlyRead bool) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	group, err := s.repository.SysGroup(ctx).FindByName(ctx, clusterId, namespace, groupName)
	if !gorm.IsRecordNotFoundError(err) {
		err = encode.ErrSysGroupExists.Error()
		_ = level.Warn(logger).Log("repository.SysGroup", "FindByName", "err", err.Error())
		return
	}
	user, err := s.repository.SysUser().Find(ctx, sysUserId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.SysUser", "Find", "err", err.Error())
		return
	}
	group.ClusterId = clusterId
	group.Namespace = namespace
	group.Name = groupName
	group.Alias = groupAlias
	group.Remark = remark
	group.OnlyRead = onlyRead
	group.UserId = userId
	group.Users = []types.SysUser{user}

	if err = s.repository.SysGroup(ctx).Save(ctx, &group, nil); err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "Save", "err", err.Error())
		err = encode.ErrSysGroupSave.Wrap(err)
		return
	}

	return
}

func (s *service) Update(ctx context.Context, sysUserId, clusterId int64, namespace, groupName, groupAlias, remark string, userId int64, onlyRead bool) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	group, err := s.repository.SysGroup(ctx).FindByName(ctx, clusterId, namespace, groupName)
	if gorm.IsRecordNotFoundError(err) {
		err = encode.ErrSysGroupNotfound.Error()
		_ = level.Warn(logger).Log("repository.SysGroup", "FindByName", "err", err.Error())
		return
	}
	if err != nil {
		_ = level.Warn(logger).Log("repository.SysGroup", "FindByName", "err", err.Error())
		err = encode.ErrSysGroupNotfound.Error()
		return
	}
	// TODO: 超管可以操作
	// 非组管理员无法操作
	if group.UserId != sysUserId {
		err = encode.ErrSysGroupNotPermission.Error()
		return
	}
	group.Alias = groupAlias
	group.Remark = remark
	group.UserId = userId
	group.OnlyRead = onlyRead
	if err = s.repository.SysGroup(ctx).Save(ctx, &group, nil); err != nil {
		_ = level.Error(logger).Log("repository.SysGroup", "Save", "err", err.Error())
		err = encode.ErrSysGroupSave.Wrap(err)
		return err
	}
	return
}

func (s *service) Delete(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	user, err := s.repository.SysUser().Find(ctx, sysUserId, "SysGroups")
	if err != nil {
		_ = level.Warn(logger).Log("repository.SysUser", "Find", "err", err.Error())
		return
	}

	fmt.Println(user.Id)

	s.repository.SysGroup(ctx)

	return
}

func (s *service) List(ctx context.Context, clusterId int64, groupIds []int64, namespace, groupName string, page, pageSize int) (res []result, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	list, total, err := s.repository.SysGroup(ctx).List(ctx, clusterId, groupIds, namespace, groupName, page, pageSize)
	if err != nil {
		_ = level.Warn(logger).Log("repository.SysGroup", "List", "err", err.Error())
		return
	}
	for _, v := range list {
		res = append(res, result{
			Alias:     v.Alias,
			Name:      v.Name,
			Namespace: v.Namespace,
			Remark:    v.Remark,
			User:      v.User.Username,
			OnlyRead:  v.OnlyRead,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return
}

func (s *service) AddUser(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, userIds []int64) (err error) {
	panic("implement me")
}

func (s *service) AddApp(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error) {
	panic("implement me")
}

func (s *service) AddCronJob(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error) {
	panic("implement me")
}

func (s *service) AddStatefulSet(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error) {
	panic("implement me")
}

func (s *service) DeleteApp(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error) {
	panic("implement me")
}

func (s *service) DeleteCronJob(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error) {
	panic("implement me")
}

func (s *service) DeleteStatefulSet(ctx context.Context, sysUserId, clusterId int64, namespace, groupName string, names []string) (err error) {
	panic("implement me")
}

func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "group", "service")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
	}
}
