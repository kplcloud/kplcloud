/**
 * @Time : 2019/6/27 10:10 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : service
 * @Software: GoLand
 */

package hooks

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/event"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/paginator"
)

var (
	ErrHooksNameExits = errors.New("hook名称已存在")
	ErrEventsNotFound = errors.New("抱歉，您要绑定的事件不存在")
	ErrHooksCreate    = errors.New("创建Hook失败")
	ErrHooksCount     = errors.New("hook总数获取失败")
	ErrHooksList      = errors.New("hook列表获取失败")
	ErrHooksGet       = errors.New("hook获取失败")
)

type Service interface {
	// 获取单个hook信息
	Get(ctx context.Context, id int) (res *types.Webhook, err error)

	// 获取hook列表
	List(ctx context.Context, name, appName, namespace string, page, limit int) (res map[string]interface{}, err error)

	// 创建hook
	Post(ctx context.Context, req hookRequest) (err error)

	// 修改hook信息
	Update(ctx context.Context, req hookRequest) (err error)

	// 删除hook信息
	Delete(ctx context.Context, id int) (err error)

	// 发送测试信号
	TestSend(ctx context.Context, id int) error
}

type service struct {
	logger       log.Logger
	repository   repository.Repository
	hookQueueSvc ServiceHookQueue
}

func NewService(logger log.Logger, store repository.Repository, hookQueueSvc ServiceHookQueue) Service {
	return &service{
		logger,
		store,
		hookQueueSvc,
	}
}

func (c *service) Get(ctx context.Context, id int) (res *types.Webhook, err error) {
	return c.repository.Webhook().FindById(id)
}

func (c *service) TestSend(ctx context.Context, id int) error {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	webHook, err := c.repository.Webhook().FindById(id)
	if err != nil {
		_ = level.Error(c.logger).Log("webhook", "FindById", "err", err.Error())
		return ErrHooksGet
	}

	//发送测试数据
	go func() {
		if err := c.hookQueueSvc.TestHookQueue(&event.WebhooksRequest{
			AppName:   webHook.AppName,
			Namespace: webHook.Namespace,
			Event:     repository.TestEvent,
			Project:   nil,
			MemberId:  memberId,
			Message:   fmt.Sprintf("测试数据"),
		}, webHook); err != nil {
			_ = level.Error(c.logger).Log("hookQueueSvc", "HookQueue", "err", err.Error())
		}
	}()

	return nil
}

func (c *service) List(ctx context.Context, name, appName, namespace string, page, limit int) (res map[string]interface{}, err error) {
	count, err := c.repository.Webhook().Count(name, appName, namespace)
	if err != nil {
		_ = level.Error(c.logger).Log("hook", "Count", "err", err.Error())
		return nil, ErrHooksCount
	}

	p := paginator.NewPaginator(page, limit, count)

	list, err := c.repository.Webhook().FindOffsetLimit(name, appName, namespace, p.Offset(), limit)
	if err != nil {
		_ = level.Error(c.logger).Log("hook", "FindOffsetLimit", "err", err.Error())
		return nil, ErrHooksList
	}
	res = map[string]interface{}{
		"list": list,
		"page": p.Result(),
	}
	return res, nil
}

func (c *service) Post(ctx context.Context, req hookRequest) (err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)
	if hook, err := c.repository.Webhook().FindByName(req.Name); err == nil {
		if hook.Name != "" {
			_ = level.Error(c.logger).Log("hook", "create webhook", "err", "name is exist")
			return ErrHooksNameExits
		}
	}
	eventList, err := c.repository.Event().FindByIds(req.Events)
	if err != nil {
		_ = level.Error(c.logger).Log("hook", "create webhook. findEventsByIds", "err", err)
		return ErrEventsNotFound
	}
	err = c.repository.Webhook().Create(&types.Webhook{
		Name:      req.Name,
		AppName:   req.AppName,
		Namespace: req.Namespace,
		URL:       req.Url,
		Status:    int(req.Status),
		Target:    req.Target,
		Token:     req.Token,
		Events:    eventList,
		AutherID:  int(userId),
	})
	if err != nil {
		_ = level.Error(c.logger).Log("hook", "create webhook database ", "err", err)
		return ErrHooksCreate
	}
	return nil
}

func (c *service) Update(ctx context.Context, req hookRequest) (err error) {
	if hook, err := c.repository.Webhook().FindByName(req.Name); err == nil {
		if hook.Name != "" && hook.ID != req.Id {
			_ = level.Error(c.logger).Log("hook", "update webhook", "err", "name is exist", "id", hook.ID)
			return ErrHooksNameExits
		}
	}
	eventList, err := c.repository.Event().FindByIds(req.Events)
	if err != nil {
		_ = level.Error(c.logger).Log("hook", "update webhook. findEventsByIds", "err", err)
		return ErrEventsNotFound
	}

	//关联关系更新
	if err := c.repository.Webhook().DeleteEvents(&types.Webhook{ID: req.Id}); err != nil {
		_ = level.Error(c.logger).Log("hook", "delete hookEvents", "err", err.Error())
	}
	if err := c.repository.Webhook().CreateEvents(&types.Webhook{ID: req.Id}, eventList...); err != nil {
		_ = level.Error(c.logger).Log("hook", "update hookEvents", "err", err.Error())
	}
	err = c.repository.Webhook().UpdateById(&types.Webhook{
		ID:        req.Id,
		Name:      req.Name,
		AppName:   req.AppName,
		Namespace: req.Namespace,
		URL:       req.Url,
		Status:    int(req.Status),
		Target:    req.Target,
		Token:     req.Token,
	})
	if err != nil {
		_ = level.Error(c.logger).Log("hook", "update webhook database ", "err", err)
		return ErrHooksCreate
	}
	return nil
}

func (c *service) Delete(ctx context.Context, id int) (err error) {
	return c.repository.Webhook().Delete(id)
}
