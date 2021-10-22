package account

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrMemberInfo      = errors.New("获取用户信息失败")
)

type Service interface {
	Detail(ctx context.Context) (res map[string]interface{}, err error)
	GetReceive(ctx context.Context) (res interface{}, err error)
	UpdateReceive(ctx context.Context, req accountReceiveRequest) (err error)
	UpdateBase(ctx context.Context, req accountBaseRequest) (err error)
	UnWechatBind(ctx context.Context) (err error)
	GetProject(ctx context.Context) (res map[string]interface{}, err error)
}

type service struct {
	logger log.Logger
	config *config.Config
	store  repository.Repository
}

func NewService(logger log.Logger, cf *config.Config, store repository.Repository) Service {
	return &service{
		logger: logger,
		config: cf,
		store:  store,
	}
}

func (c *service) Detail(ctx context.Context) (res map[string]interface{}, err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)

	memberInfo, err := c.store.Member().GetInfoById(memberId)
	memberNamespace, err := c.store.Member().GetNssByMemberId(memberId)
	memberRoles, err := c.store.Member().GetRolesByMemberId(memberId)

	param := map[string]string{
		"is_read": "0",
	}
	unRendNum, err := c.store.NoticeMember().CountRead(param, memberId)

	if err != nil {
		_ = level.Error(c.logger).Log("account.Detail", "member.GetInfoById", "err", err.Error())
		return nil, ErrMemberInfo
	}

	res = map[string]interface{}{
		"id":          memberInfo.ID,
		"email":       memberInfo.Email,
		"username":    memberInfo.Username,
		"state":       memberInfo.State,
		"notifyCount": unRendNum,
		"openid":      memberInfo.Openid,
		"phone":       memberInfo.Phone,
		"city":        memberInfo.City,
		"department":  memberInfo.Department,
		"roles":       memberRoles,
		"namespaces":  memberNamespace,
		"attrs":       "",
	}
	return
}

func (c *service) GetReceive(ctx context.Context) (res interface{}, err error) {
	//获取所有通知事件
	eventList, err := c.store.Event().FindAllEvents()
	if err != nil {
		_ = level.Error(c.logger).Log("GetReceive.eventList", "event.FindAllEvents", "err", err.Error())
		return
	}

	type event struct {
		ID          int
		Name        string
		Description string
	}

	DefaultEventList := []*event{}

	alarmEvent := event{0, "Alarm", "告警"}
	proclaimEvent := event{0, "Proclaim", "公告"}

	DefaultEventList = append(DefaultEventList, &proclaimEvent, &alarmEvent)

	memberId := ctx.Value(middleware.UserIdContext).(int64)
	//获取用户配置
	memberReceive, err := c.store.NoticeReceive().FindListByMid(memberId)

	//组合用户配置
	data := []*receiveResponse{}

	//公告/告警
	for _, de := range DefaultEventList {
		d := new(receiveResponse)
		d.Action = de.Name
		d.ActionDesc = de.Description

		for _, mr := range memberReceive {
			if de.Name == mr.NoticeAction {
				//if mr.Site == 1{ //站内信 公告，告警默认选中站内信
				d.Site = 1
				//}
				//if mr.NoticeAction =="Alarm"{ //微信 告警默认选中微信
				//	d.Wechat = 1
				//}
				if mr.Wechat == 1 { //微信
					d.Wechat = 1
				}
				if mr.Email == 1 { //邮箱
					d.Email = 1
				}
				if mr.Sms == 1 { //短信
					d.Sms = 1
				}
				if mr.Bee == 1 { //蜜蜂
					d.Bee = 1
				}

			}
		}
		data = append(data, d)
	}

	//通知
	for _, e := range eventList {
		d := new(receiveResponse)
		d.Action = e.Name.String
		d.ActionDesc = e.Description.String

		for _, mr := range memberReceive {
			if e.Name.String == mr.NoticeAction {
				if mr.Site == 1 { //站内信
					d.Site = 1
				}
				if mr.Wechat == 1 { //微信
					d.Wechat = 1
				}
				if mr.Email == 1 { //邮箱
					d.Email = 1
				}
				if mr.Sms == 1 { //短信
					d.Sms = 1
				}
				if mr.Bee == 1 { //蜜蜂
					d.Bee = 1
				}

			}
		}

		data = append(data, d)
	}

	res = data
	return
}

func (c *service) UpdateReceive(ctx context.Context, req accountReceiveRequest) (err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)

	var data = new(types.NoticeReceive)
	data.Email = req.Email
	data.Wechat = req.Wechat
	data.Sms = req.Sms
	data.Bee = req.Bee
	data.Site = req.Site
	data.MemberID = int(memberId)
	data.NoticeAction = req.Action

	//当前信息
	curAction, _ := c.store.NoticeReceive().GetNoticeReceiveByMidAction(memberId, req.Action)

	if curAction != nil {
		//update
		data.ID = curAction.ID
		err := c.store.NoticeReceive().Update(data)
		if err != nil {
			_ = level.Error(c.logger).Log("UpdateReceive", "c.noticeReceive.Update", "err", err.Error())
		}
	} else {
		//insert
		err := c.store.NoticeReceive().Create(data)
		if err != nil {
			_ = level.Error(c.logger).Log("UpdateReceive", "c.noticeReceive.Create", "err", err.Error())
		}
	}
	return
}

func (c *service) UpdateBase(ctx context.Context, req accountBaseRequest) (err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)

	var data = new(types.Member)
	data.Username = req.Name
	data.City = req.City
	data.Department = req.Department
	data.Phone = req.Phone
	data.ID = memberId

	err = c.store.Member().Update(data)
	if err != nil {
		_ = level.Error(c.logger).Log("account.UpdateBase", "c.member.Update", "err", err.Error())
	}

	return
}

func (c *service) UnWechatBind(ctx context.Context) (err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	err = c.store.Member().UnBindWechat(memberId)
	return
}

func (c *service) GetProject(ctx context.Context) (res map[string]interface{}, err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	re, err := c.store.Project().GetProjectByMid(memberId)
	res = map[string]interface{}{
		"list": re,
	}
	return
}
