/**
 * @Time : 2019-07-15 10:55
 * @Author : soupzhb@gmail.com
 * @File : alarm.go
 * @Software: GoLand
 */

package msgs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	amqpClient "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/email"
	"github.com/kplcloud/kplcloud/src/pkg/public"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"html/template"
	"os"
	"time"
)

type alarm struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	AlarmAction = "Alarm" //告警
)

type ServiceAlarm interface {
	DistributeAlarm(ctx context.Context, data string) (err error)
	CreateAlarmMember(ctx context.Context, ml []types.Member, notice types.Notices, al alarm) (num int, err error)
	SendAlarmEmail(toUser, toCc []string, notice types.Notices) (err error)
}

type serviceAlarm struct {
	config     *config.Config
	logger     log.Logger
	mailClient email.EmailInterface
	amqpClient amqpClient.AmqpClient
	store      repository.Repository
}

/**
 * @Title 处理消费出来的 alarm 的数据
 */
func NewServiceAlarm(logger log.Logger,
	cf *config.Config,
	mailClient email.EmailInterface,
	amqpClient amqpClient.AmqpClient,
	store repository.Repository) ServiceAlarm {
	return &serviceAlarm{
		logger:     logger,
		config:     cf,
		mailClient: mailClient,
		amqpClient: amqpClient,
		store:      store,
	}
}

func (c *serviceAlarm) DistributeAlarm(ctx context.Context, data string) (err error) {

	if len(data) <= 0 {
		return nil
	}

	//读取mq内容
	var al alarm
	err = json.Unmarshal([]byte(data), &al)

	if err != nil {
		_ = level.Error(c.logger).Log("DistributeAlarm", "json.Unmarshal", "data", al)
		return
	}

	//写DB
	var v types.Notices
	v.Title = al.Title
	v.Name = al.Name
	v.Content = al.Content
	v.Namespace = al.Namespace
	v.Action = AlarmAction
	v.Type = 3

	_ = c.store.Notice().CreateReturnId(&v)

	//站内信分发
	//根据namespace和name获取成员id
	var ml []types.Member
	ml = c.store.Groups().GetMemberIdsByProjectNameAndNs(v.Name, v.Namespace)
	//如果获取不到分发给配置中默认接收者
	if len(ml) == 0 {
		kaEmailArr := c.config.GetStrings("msg", "alarm_default_email")
		ml, err = c.store.Member().GetMembersByEmails(kaEmailArr)
		if err != nil {
			_ = level.Error(c.logger).Log("DistributeAlarm", "c.store.Member().GetMembersByEmails", "err", err)
		}
	}

	num, err := c.CreateAlarmMember(ctx, ml, v, al)
	if err != nil {
		_ = level.Error(c.logger).Log("DistributeAlarm", "c.CreateAlarmMember", "err", err)
	}
	_ = level.Info(c.logger).Log("DistributeAlarm", "c.CreateAlarmMember", "num", num)

	return

}

//告警-用户关系
func (c *serviceAlarm) CreateAlarmMember(ctx context.Context, ml []types.Member, notice types.Notices, al alarm) (n int, err error) {
	data := []*types.NoticeMember{}

	//获取订阅此action通知的用户列表
	rl, err := c.store.NoticeReceive().GetNoticeReceiveByAction(notice.Action)

	//获取项目创建人
	var adminEmail string
	pro, notExist := c.store.Project().FindByNsNameExist(notice.Namespace, notice.Name)

	if notExist == true {
		admin, _ := c.store.Member().GetInfoById(pro.MemberID)
		adminEmail = admin.Email
	}

	var toUser, toCc []string

	//to wx queue
	wechatQueueSvc := NewServiceWechatQueue(c.logger, c.config, c.amqpClient, nil, c.store)

	for _, v := range ml {
		if v.State == 2 { //如果用户状态已关闭，则不再发送消息
			continue
		}
		d := types.NoticeMember{}
		d.MemberID = v.ID
		d.NoticeID = notice.ID

		var isSite, isWechat, isEmail, isSms, isBee bool

		//接收人&&订阅关系
		for _, vr := range rl {
			if int(vr.MemberID) == int(v.ID) {
				if vr.Site == 1 { //站内信
					isSite = true
				}
				if vr.Wechat == 1 { //微信
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
			isSite = true
		}

		//站内信
		if isSite == true { //站内信
			data = append(data, &d) //如果订阅了，则入库，否则不入库
		}

		if isWechat == true { //微信  增加一个通知模板
			//微信公众号信息
			wxNotice := notice
			wxNotice.Content = al.Desc
			err := wechatQueueSvc.PublicWechatQueue(v, wxNotice, "")
			if err != nil {
				_ = level.Error(c.logger).Log("amqpClient", "wechatQueueSvc.PublicWechatQueue", "err", err)
			}
		}
		if isEmail == true { //邮箱
			if v.Email != adminEmail {
				if v.State != 2 {
					toCc = append(toCc, v.Email)
				}
			} else {
				if v.State != 2 {
					toUser = append(toUser, v.Email)
				}

			}
		}
		if isSms == true { //短信

		}
		if isBee == true { //蜜蜂

		}
	}

	//用户-报警关系入库
	n, err = c.store.NoticeMember().InsertMulti(data)

	if len(toUser) == 0 {
		toUser = c.config.GetStrings("msg", "alarm_default_email")
	}

	go func() {
		err = c.SendAlarmEmail(toUser, toCc, notice)
		if err != nil {
			_ = level.Error(c.logger).Log("amqpClient", "c.SendAlarmEmail", "Err", err)
		}
	}()

	return
}

type alarmNotice struct {
	Notice  types.Notices
	Content public.Prom
	FromOs  string
}

//告警-邮件
func (c *serviceAlarm) SendAlarmEmail(toUser, toCc []string, notice types.Notices) (err error) {
	an := new(alarmNotice)
	an.Notice = notice
	an.FromOs = "hostname:" + os.Getenv("HOSTNAME") + "; env:" + os.Getenv("ENV")

	err = json.Unmarshal([]byte(notice.Content), &an.Content)

	if err != nil {
		_ = level.Error(c.logger).Log("DistributeAlarm", "SendAlarmEmail", "json.Unmarshal.Err", err.Error())
		return
	}

	//获取模板
	temp, _ := c.store.Template().FindByKindType(repository.EmailAlarm)
	tpl := temp.Detail
	tmpl := template.New("alarm")
	_, err = tmpl.Parse(tpl)

	if err != nil {
		_ = level.Error(c.logger).Log("DistributeAlarm", "SendAlarmEmail", "tmpl.Parse.Err", err.Error())
		return
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, an)

	if err != nil {
		_ = level.Error(c.logger).Log("DistributeAlarm", "SendAlarmEmail", "tmpl.Execute.Err", err.Error())
		return
	}

	fmt.Println("sendAlarmEmail-toUser:", toUser, "toCc:", toCc, "HOSTNAME:", an.FromOs, "string:", body.String())

	c.mailClient.SetTitle(notice.Title)
	c.mailClient.SetContent(body.String())
	c.mailClient.AddEmailAddress(toUser) //toUser
	c.mailClient.AddCcEmailAddress(toCc) //toCc
	c.mailClient.SetContentType("text/html;charset=utf-8")
	err = c.mailClient.Send()

	if err != nil {
		_ = level.Error(c.logger).Log("DistributeAlarm", "SendAlarmEmail", "err", err.Error())
	}

	return
}
