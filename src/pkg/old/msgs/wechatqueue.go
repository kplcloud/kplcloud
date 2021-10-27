/**
 * @Time : 2019-07-16 14:58
 * @Author : soupzhb@gmail.com
 * @File : wechatqueue.go
 * @Software: GoLand
 */

package msgs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	amqpClient "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/event"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/yijizhichang/wechat-sdk"
	"github.com/yijizhichang/wechat-sdk/mp/message/template"
	"time"
)

var ErrWechatSend = errors.New("微信消息推送失败")

type TemplateMsg struct {
	FirstValue    string
	FirstColor    string
	Keyword1Value string
	Keyword1Color string
	Keyword2Value string
	Keyword2Color string
	Keyword3Value string
	Keyword3Color string
	RemarkValue   string
	RemarkColor   string
	Url           string
	TemplateID    string
	ToUser        string
}

type WechatNotice struct {
	Id          int64     `json:"id"`
	Type        int       `json:"type"`
	Action      string    `json:"action"`
	Name        string    `json:"name"`
	Namespace   string    `json:"namespace"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedName string    `json:"created_name"`
	ToUser      string    `json:"to_user"`
}

type ServiceWechatQueue interface {
	PublicWechatQueue(m types.Member, notice types.Notices, createdName string) (err error)
	PublicNoticeQueue(req *event.WebhooksRequest) (err error)
	PublicProclaimQueue(proclaim *types.Notices) (err error)
	DistributeMsgWechat(ctx context.Context, data string) (err error)
}

type serviceWechatQueue struct {
	config     *config.Config
	logger     log.Logger
	amqpClient amqpClient.AmqpClient
	wxClient   *wechat.Wechat
	store      repository.Repository
}

/**
 * @Title 将待发送至微信的消息放到队列
 */
func NewServiceWechatQueue(logger log.Logger,
	cf *config.Config,
	amqpClient amqpClient.AmqpClient,
	wxClient *wechat.Wechat,
	store repository.Repository) ServiceWechatQueue {
	return &serviceWechatQueue{
		logger:     logger,
		config:     cf,
		amqpClient: amqpClient,
		wxClient:   wxClient,
		store:      store,
	}
}

//消息分发给微信
func (c *serviceWechatQueue) DistributeMsgWechat(ctx context.Context, data string) (err error) {
	//读取mq内容
	var wxMsg WechatNotice
	err = json.Unmarshal([]byte(data), &wxMsg)

	if err != nil {
		_ = level.Error(c.logger).Log("DistributeMsgWechat", "json.Unmarshal", "data", wxMsg)
		return
	}

	tpl := c.wxClient.GetTemplate()

	msgText := c.SendTemplateWechatNotice(wxMsg)

	fmt.Println("msgText", msgText)

	re, _ := tpl.Send(msgText)
	if re != nil && re.ErrCode == 0 {
		return nil
	} else {
		return ErrWechatSend
	}

}

//待推微信的消息存入MQ
func (c *serviceWechatQueue) PublicWechatQueue(m types.Member, notice types.Notices, createdName string) (err error) {

	data := new(WechatNotice)
	data.Id = notice.ID
	data.Type = int(notice.Type)
	data.Action = notice.Action
	data.Name = notice.Name
	data.Namespace = notice.Namespace
	data.Title = notice.Title
	data.Content = notice.Content
	data.CreatedAt = notice.CreatedAt.Time
	data.CreatedName = createdName
	data.ToUser = m.Openid

	b, err := json.Marshal(data)
	_ = level.Debug(c.logger).Log("amqpClient", "PublishOnQueue", "data-json", b)

	defer func() {

		if err := c.amqpClient.PublishOnQueue(amqpClient.MsgWechatTopic, func() []byte {
			return []byte(b)
		}); err != nil {
			_ = level.Error(c.logger).Log("amqpClient", "PublishOnQueue", "err", err.Error())
		}
		time.Sleep(time.Second * 2)
	}()

	return
}

func (c *serviceWechatQueue) PublicNoticeQueue(req *event.WebhooksRequest) (err error) {
	if req.AppName == "" {
		return errors.New("项目名称不能为空")
	}

	if req.Namespace == "" {
		return errors.New("命名空间不能为空")
	}

	//消息组装
	var msgText string
	if req.Title == "" {
		msgText = "项目名称" + "【" + req.AppName + "】" + "命名空间" + "【" + req.Namespace + "】" + "已操作" + "【" + req.EventDesc + "】"
	} else {
		msgText = req.Title
	}

	//写DB
	var v types.Notices
	v.Title = msgText
	v.Content = req.Message
	v.Name = req.AppName
	v.Namespace = req.Namespace
	v.Action = string(req.Event)
	v.Type = 2
	v.MemberID = int(req.Member.ID)

	id := c.store.Notice().CreateReturnId(&v)

	//写MQ
	notice := v
	notice.ID = id

	noticeMqData := new(NoticeMqData)
	noticeMqData.WebHooksReq = *req
	noticeMqData.Notice = notice

	b, _ := json.Marshal(noticeMqData)

	defer func() {
		if err := c.amqpClient.PublishOnQueue(amqpClient.NoticeTopic, func() []byte {
			return []byte(b)
		}); err != nil {
			_ = level.Error(c.logger).Log("amqpClient", "PublicNoticeQueue", "err", err.Error())
		}
	}()

	return
}

func (c *serviceWechatQueue) PublicProclaimQueue(proclaim *types.Notices) (err error) {
	b, err := json.Marshal(proclaim)
	_ = level.Debug(c.logger).Log("amqpClient", "PublicProclaimQueue", "data-json", b)

	defer func() {
		if err := c.amqpClient.PublishOnQueue(amqpClient.ProclaimTopic, func() []byte {
			return []byte(b)
		}); err != nil {
			_ = level.Error(c.logger).Log("amqpClient", "PublicProclaimQueue", "err", err.Error())
		}
	}()

	return
}

//wechatNotice 不同的类型，不同的模板
func (c *serviceWechatQueue) SendTemplateWechatNotice(wechatNotce WechatNotice) *template.Message {
	curTemplateMsg := new(TemplateMsg)

	switch wechatNotce.Type {
	case 1:
		//公告
		//curTemplateMsg.FirstValue = wechatNotce.Title
		//curTemplateMsg.FirstColor = "#000"
		//curTemplateMsg.Keyword1Value = "proclaim~"
		//curTemplateMsg.Keyword1Color = "#EEB422"
		//curTemplateMsg.Keyword2Value = ""
		//curTemplateMsg.Keyword2Color = "#FF4500"
		//curTemplateMsg.Keyword3Value = wechatNotce.CreatedAt.String()
		//curTemplateMsg.Keyword3Color = "#000"
		//curTemplateMsg.RemarkValue = "公告内容：" + wechatNotce.Content
		//curTemplateMsg.RemarkColor = "#358DE1"
		//curTemplateMsg.TemplateID = c.config.GetString("wechat", "tpl_proclaim")
		//curTemplateMsg.ToUser = wechatNotce.ToUser
	case 2:
		//通知
		if len(wechatNotce.Content) > 100 {
			wechatNotce.Content = string([]rune(wechatNotce.Content)[:100]) + "..."
		}
		curTemplateMsg.FirstValue = "事件类型：" + wechatNotce.Action
		curTemplateMsg.FirstColor = "#358DE1"
		curTemplateMsg.Keyword1Value = wechatNotce.Title
		curTemplateMsg.Keyword1Color = "#000"
		curTemplateMsg.Keyword2Value = wechatNotce.CreatedAt.Format("2006-01-02 15:04:05")
		curTemplateMsg.Keyword2Color = "#000"
		curTemplateMsg.RemarkValue = wechatNotce.Content
		curTemplateMsg.RemarkColor = "#FF4500"
		curTemplateMsg.TemplateID = c.config.GetString("wechat", "tpl_notice")
		curTemplateMsg.ToUser = wechatNotce.ToUser
	case 3:
		//告警
		if len(wechatNotce.Content) > 100 {
			wechatNotce.Content = string([]rune(wechatNotce.Content)[:100]) + "..."
		}
		curTemplateMsg.FirstValue = "项目名称：" + wechatNotce.Name + "\n" + "业务线：" + wechatNotce.Namespace
		curTemplateMsg.FirstColor = "#358DE1"
		curTemplateMsg.Keyword1Value = wechatNotce.CreatedAt.Format("2006-01-02 15:04:05")
		curTemplateMsg.Keyword1Color = "#000"
		curTemplateMsg.Keyword2Value = wechatNotce.Title
		curTemplateMsg.Keyword2Color = "#000"
		curTemplateMsg.RemarkValue = "报警摘要：" + wechatNotce.Content
		curTemplateMsg.RemarkColor = "#FF4500"
		curTemplateMsg.Url = ""
		curTemplateMsg.TemplateID = c.config.GetString("wechat", "tpl_alarm")
		curTemplateMsg.ToUser = wechatNotce.ToUser

	}

	//fmt.Println("接收消息内容：", wechatNotce, "模板内容日志：", curTemplateMsg)

	msgText := new(template.Message)
	msgText.ToUser = curTemplateMsg.ToUser
	msgText.TemplateID = curTemplateMsg.TemplateID
	msgText.URL = curTemplateMsg.Url
	msgText.Data = make(map[string]*template.DataItem)
	msgText.Data["first"] = &template.DataItem{curTemplateMsg.FirstValue, curTemplateMsg.FirstColor}
	msgText.Data["keyword1"] = &template.DataItem{curTemplateMsg.Keyword1Value, curTemplateMsg.Keyword1Color}
	msgText.Data["keyword2"] = &template.DataItem{curTemplateMsg.Keyword2Value, curTemplateMsg.Keyword2Color}
	msgText.Data["keyword3"] = &template.DataItem{curTemplateMsg.Keyword3Value, curTemplateMsg.Keyword3Color}
	msgText.Data["remark"] = &template.DataItem{curTemplateMsg.RemarkValue, curTemplateMsg.RemarkColor}

	return msgText
}
