/**
 * @Time : 2019-07-15 10:56
 * @Author : soupzhb@gmail.com
 * @File : notice.go
 * @Software: GoLand
 */

package msgs

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	amqpClient "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/email"
	"github.com/kplcloud/kplcloud/src/event"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"html/template"
	"os"
	"strings"
)

//通知消息类型-MQ
type NoticeMqData struct {
	WebHooksReq event.WebhooksRequest
	Notice      types.Notices
	FromOs      string
}

type ServiceNotice interface {
	DistributeNotice(ctx context.Context, data string) error
	CreateNoticeMember(ml []types.Member, notice *NoticeMqData) (num int, err error)
	SendNoticeEmail(toUser, toCc []string, notice NoticeMqData) (err error)
}

type serviceNotice struct {
	config     *config.Config
	logger     log.Logger
	mailClient email.EmailInterface
	amqpClient amqpClient.AmqpClient
	store      repository.Repository
}

/**
 * @Title 处理消费出来的 notice 的数据
 */
func NewServiceNotice(logger log.Logger,
	cf *config.Config,
	mailClient email.EmailInterface,
	amqpClient amqpClient.AmqpClient,
	store repository.Repository) ServiceNotice {
	return &serviceNotice{
		logger:     logger,
		config:     cf,
		mailClient: mailClient,
		amqpClient: amqpClient,
		store:      store,
	}
}

//通知分发
func (c *serviceNotice) DistributeNotice(ctx context.Context, data string) (err error) {
	if len(data) <= 0 {
		return nil
	}

	//读取mq内容
	dat := new(NoticeMqData)
	err = json.Unmarshal([]byte(data), &dat)

	if err != nil {
		_ = level.Error(c.logger).Log("DistributeNotice", "json.Unmarshal", "err", err)
		return
	}

	//站内信分发
	//根据namespace和name获取成员id
	var ml []types.Member
	ml = c.store.Groups().GetMemberIdsByProjectNameAndNs(dat.Notice.Name, dat.Notice.Namespace)

	num, err := c.CreateNoticeMember(ml, dat)
	if err != nil {
		_ = level.Error(c.logger).Log("DistributeNotice", "c.CreateNoticeMember", "err", err)

	}
	_ = level.Info(c.logger).Log("DistributeNotice", "insert db", "num", num, "member-list", ml)

	return
}

//通知-用户关系  跟项目关联 （action-member除外）
func (c *serviceNotice) CreateNoticeMember(ml []types.Member, notice *NoticeMqData) (n int, err error) {
	data := []*types.NoticeMember{}

	//获取订阅此action通知的用户列表
	rl, err := c.store.NoticeReceive().GetNoticeReceiveByAction(notice.Notice.Action)

	var toCc, toUser []string

	//toUser := []string{notice.WebHooksReq.Member.Email}  //创建者为主接收者

	//to wx queue
	wechatQueueSvc := NewServiceWechatQueue(c.logger, c.config, c.amqpClient, nil, c.store)

	for _, v := range ml {
		if v.State == 2 { //如果用户状态已关闭，则不再发送消息
			continue
		}
		d := types.NoticeMember{}
		d.MemberID = v.ID
		d.NoticeID = int64(notice.Notice.ID)

		var isSite, isWechat, isEmail, isSms, isBee bool
		//接收人&&订阅关系
		for _, vr := range rl {
			if int(vr.MemberID) == int(v.ID) {
				if vr.Site == 1 { //站内信
					isSite = true
				}
				if vr.Wechat == 1 { //微信  增加一个通知模板
					isWechat = true
				}
				if vr.Email == 1 { //邮箱
					isEmail = true
				}
				if vr.Sms == 1 { //短信
					isSms = true
				}
				if vr.Bee == 1 { //蜜蜂
					isBee = true
				}
			}

			//默认值
			if vr.NoticeAction == "Build" || vr.NoticeAction == "Delete" { //站内信
				isSite = true
			}
			if vr.NoticeAction == "Delete" { //微信  增加一个通知模板
				isWechat = true
			}
			if vr.NoticeAction == "Build" || vr.NoticeAction == "Delete" { //邮箱
				isEmail = true
			}
		}

		//站内信
		if isSite == true { //站内信
			data = append(data, &d) //如果订阅了，则入库，否则不入库
		}

		if isWechat == true { //微信  增加一个通知模板
			err := wechatQueueSvc.PublicWechatQueue(v, notice.Notice, notice.WebHooksReq.Member.Username)
			if err != nil {
				_ = level.Error(c.logger).Log("DistributeNotice.amqpClient", "wechatQueueSvc.PublicWechatQueue", "err", err)
			}
		}
		if isEmail == true { //邮箱

			if v.Email != notice.WebHooksReq.Member.Email { //抄送时去除创建者，创建者为接收者
				toCc = append(toCc, v.Email)
			} else {
				toUser = append(toUser, v.Email) //如果创建者订阅了，则创建者为主接收
			}
		}
		if isSms == true { //短信

		}
		if isBee == true { //蜜蜂

		}
	}

	if len(toUser) == 0 {
		toUser = c.config.GetStrings("msg", "alarm_default_email")
	}

	//发送通知邮件至订阅人
	go func() {
		err = c.SendNoticeEmail(toUser, toCc, *notice)
		if err != nil {
			_ = level.Error(c.logger).Log("DistributeNotice", "c.SendNoticeEmail", "err", err.Error())
		}
	}()

	//用户-通知关系入库
	if len(data) > 0 {
		n, err = c.store.NoticeMember().InsertMulti(data)
		return
	} else {
		_ = level.Error(c.logger).Log("DistributeNotice", "no member subscribe", "action", notice.Notice.Action)
		return 0, nil
	}
}

//通知-邮件
func (c *serviceNotice) SendNoticeEmail(toUser, toCc []string, notice NoticeMqData) (err error) {

	notice.WebHooksReq.Message = strings.Replace(notice.WebHooksReq.Message, "\n", "@@@@@", -1)
	notice.FromOs = "hostname:" + os.Getenv("HOSTNAME") + "; env:" + os.Getenv("ENV")

	//获取模板
	temp, _ := c.store.Template().FindByKindType(repository.EmailNotice)
	tpl := temp.Detail
	tmpl := template.New("notice")
	_, err = tmpl.Parse(tpl)
	if err != nil {
		_ = level.Error(c.logger).Log("SendNoticeEmail", "tmpl.Parse", "err", err.Error())
		return
	}
	var body bytes.Buffer
	err = tmpl.Execute(&body, notice)

	if err != nil {
		_ = level.Error(c.logger).Log("SendNoticeEmail", "tmpl.Execute", "err", err.Error())
		return
	}

	_ = level.Info(c.logger).Log("sendNoticeEmail-toUser:", toUser, "toCc:", toCc, "HOSTNAME:", notice.FromOs, "string:", body.String())

	c.mailClient.SetTitle(notice.Notice.Title)
	c.mailClient.SetContent(body.String())
	c.mailClient.AddEmailAddress(toUser) //toUser
	c.mailClient.AddCcEmailAddress(toCc) //toCc
	c.mailClient.SetContentType("text/html;charset=utf-8")
	err = c.mailClient.Send()

	if err != nil {
		_ = level.Error(c.logger).Log("DistributeNotice", "SendNoticeEmail", "err", err.Error())
	}

	return
}
