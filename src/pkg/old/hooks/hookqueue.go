/**
 * @Time : 2019/8/8 2:41 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : hookqueue
 * @Software: GoLand
 */

package hooks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/event"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/notice"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/request"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Data struct {
	Url      string
	Params   url.Values
	Headers  string
	ProxyUrl string //代理
	Count    int
}

type ServiceHookQueue interface {
	/**
	 * @Title hook数据写入队列及event
	 */
	SendHookQueue(ctx context.Context, kind repository.EventsKind, name, ns, msg string) error

	/**
	 * @Title 发送一条测试数据
	 */
	TestHookQueue(req *event.WebhooksRequest, webHook *types.Webhook) error

	/**
	 * @Title 处理消费出来的 hook数据
	 */
	HookReceiver(ctx context.Context, data string) error
}

type serviceHookQueue struct {
	logger     log.Logger
	amqpClient amqp.AmqpClient
	conf       *config.Config
	repository repository.Repository
	noticeSvc  notice.Service
}

func NewServiceHookQueue(logger log.Logger,
	amqpClient amqp.AmqpClient,
	conf *config.Config,
	repository repository.Repository,
	noticeSvc notice.Service) ServiceHookQueue {
	return &serviceHookQueue{
		logger,
		amqpClient,
		conf,
		repository,
		noticeSvc,
	}
}

func (c *serviceHookQueue) SendHookQueue(ctx context.Context, kind repository.EventsKind, name, ns, msg string) error {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)
	memberId := ctx.Value(middleware.UserIdContext).(int64)

	eventInfo, err := c.repository.Event().FindByKind(kind)
	if err != nil {
		_ = c.logger.Log("BindHooks", "FindEvents", "err", err.Error())
		return err
	}

	member, err := c.repository.Member().FindById(memberId)
	if err != nil {
		_ = c.logger.Log("BindHooks", "Find Member By Id", "err", err.Error())
		return err
	}

	// 调消息中心
	var req = &event.WebhooksRequest{
		AppName:   name,
		Namespace: ns,
		Member:    member,
		Project:   project,
		Event:     kind,
		EventDesc: eventInfo.Description.String,
		Message:   msg,
	}
	go func() {
		if err := c.noticeSvc.Create(ctx, req); err != nil {
			_ = c.logger.Log("BindHooks", "Create Notice", "err", err.Error())
		}
	}()

	for _, val := range eventInfo.WebHook {
		if val.Status != 1 {
			continue
		}
		if val.Target == repository.AppTarget {
			if project.Name != val.AppName || project.Namespace != val.Namespace {
				continue
			}
		}
		go c.to(req, val)
	}
	return nil
}

func (c *serviceHookQueue) TestHookQueue(req *event.WebhooksRequest, webHook *types.Webhook) error {
	if webHook.Status != 1 {
		return errors.New("请先激活当前Hook")
	}
	go c.to(req, webHook)
	return nil
}

func (c *serviceHookQueue) to(req *event.WebhooksRequest, webHook *types.Webhook) error {
	var httpProxy string
	params := url.Values{}

	var param = event.Params{
		AppName:       req.AppName,
		Namespace:     req.Namespace,
		Event:         req.Event.String(),
		Message:       req.Message,
		Operator:      req.Member.Username,
		OperatorEmail: req.Member.Email,
	}
	if req.Project != nil {
		param.Project.Name = req.Project.Name
		param.Project.NameEn = req.Project.Name
		param.Project.Namespace = req.Project.Namespace
		param.Project.Member = req.Project.Member.Username
		param.Project.Email = req.Project.Member.Email
		param.Project.ProjectId = req.Project.ID
		param.Project.Description = req.Project.Desc
	}

	dat, _ := json.Marshal(param)
	go func() {
		if err := c.repository.EventHistory().Create(&types.EventHistory{
			AppName:   req.AppName,
			Namespace: req.Namespace,
			Event:     req.Event.String(),
			Date:      string(dat),
		}); err != nil {
			_ = c.logger.Log("EventHistory", "Create", "err", err.Error())
		}
	}()

	if c.conf.GetString("server", "http_proxy") != "" {
		httpProxy = c.conf.GetString("server", "http_proxy")
	}

	var data = &Data{
		Url:      webHook.URL,
		Params:   params,
		Headers:  webHook.Token,
		ProxyUrl: httpProxy,
	}
	b, _ := json.Marshal(data)

	defer func() {
		if err := c.amqpClient.PublishOnQueue(amqp.HookTopic, func() []byte {
			return b
		}); err != nil {
			_ = c.logger.Log("amqpClient", "PublishOnQueue", "err", err.Error())
		}
		time.Sleep(time.Second * 2)
	}()
	return nil
}

func (c *serviceHookQueue) HookReceiver(ctx context.Context, data string) error {
	if len(data) <= 0 {
		return nil
	}
	//读取mq内容
	var dat Data
	if err := json.Unmarshal([]byte(data), &dat); err != nil {
		return err
	}
	_ = c.logger.Log("HookReceiver", "json.Unmarshal", "data", data)

	// request请求
	params := dat.Params.Encode()
	requestCli := request.NewRequest(dat.Url, "POST")

	if dat.Headers != "" {
		requestCli = requestCli.Header("X-Kpl-Token", dat.Headers)
	}

	if dat.ProxyUrl != "" {
		dialer := &net.Dialer{
			Timeout:   time.Duration(30 * time.Second),
			KeepAlive: time.Duration(30 * time.Second),
		}
		requestCli = requestCli.HttpClient(&http.Client{
			Transport: &http.Transport{
				Proxy: func(_ *http.Request) (*url.URL, error) {
					return url.Parse(dat.ProxyUrl)
				},
				DialContext: dialer.DialContext,
			},
		})
	}
	var out interface{}
	if err := requestCli.Body(strings.NewReader(params)).Do().Into(&out); err != nil {
		_ = c.logger.Log("HookReceiver", "Send Request", "err", err.Error())
	}
	return nil
}
