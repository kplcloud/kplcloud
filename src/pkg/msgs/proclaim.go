/**
 * @Time : 2019-07-15 10:57
 * @Author : soupzhb@gmail.com
 * @File : proclaim.go
 * @Software: GoLand
 */

package msgs

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	amqpClient "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/email"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"strings"
)

//消息动作类型
type ProclaimType string

//不支持消息动作类型
var ErrProclaimType = errors.New("不支持的公告类型")

const (
	ProclaimAction = "Proclaim" //公告
)

const (
	ProclaimTypeAll       ProclaimType = "all"       //全部
	ProclaimTypeNamespace              = "namespace" //业务线
	ProclaimTypeUser                   = "user"      //指定用户
)

type ServiceProclaim interface {
	DistributeProclaim(ctx context.Context, data string) error
	CreateProclaimMember(ml []types.Member, notice *types.Notices) (num int, err error)
	SendProclaimEmail(toUser, toCc []string, notice types.Notices) (err error)
}

type serviceProclaim struct {
	config     *config.Config
	logger     log.Logger
	amqpClient amqpClient.AmqpClient
	mailClient email.EmailInterface
	store      repository.Repository
}

/**
 * @Title 处理消费出来的 alarm 的数据
 */
func NewServiceProclaim(logger log.Logger,
	cf *config.Config,
	mailClient email.EmailInterface,
	amqpClient amqpClient.AmqpClient,
	store repository.Repository) ServiceProclaim {
	return &serviceProclaim{
		logger:     logger,
		config:     cf,
		mailClient: mailClient,
		amqpClient: amqpClient,
		store:      store,
	}
}

//公告分发
func (c *serviceProclaim) DistributeProclaim(ctx context.Context, data string) (err error) {

	if len(data) <= 0 {
		return nil
	}

	//读取mq内容
	var dat *types.Notices
	err = json.Unmarshal([]byte(data), &dat)
	if err != nil {
		_ = level.Error(c.logger).Log("DistributeProclaim", "json.Unmarshal", "err", err.Error())
		return
	}

	//站内信分发
	switch dat.ProclaimType {
	case "all":
		ml, err := c.store.Member().GetMembersAll()
		if err != nil {
			_ = level.Error(c.logger).Log("DistributeProclaim", "c.member.GetMembersAll", "err", err.Error())
		}

		num, err := c.CreateProclaimMember(ml, dat)
		if err != nil {
			_ = level.Error(c.logger).Log("DistributeProclaim", "c.CreateProclaimMember", "err", err.Error())
		}
		_ = level.Info(c.logger).Log("DistributeProclaim", "c.CreateProclaimMember", "num", num, "case", "all")
	case "namespace":
		idStrArr := strings.Split(dat.ProclaimReceive, ",")
		ml, err := c.store.Member().GetMembersByNss(idStrArr)

		if err != nil {
			_ = level.Error(c.logger).Log("DistributeProclaim", "c.member.GetMembersByNss", "err", err.Error())
		}
		num, err := c.CreateProclaimMember(ml, dat)
		if err != nil {
			_ = level.Error(c.logger).Log("DistributeProclaim", "c.CreateProclaimMember", "err", err.Error())
		}
		_ = level.Info(c.logger).Log("DistributeProclaim", "c.CreateProclaimMember", "num", num, "case", "namespace")
	case "user":
		idStrArr := strings.Split(dat.ProclaimReceive, ",")
		ml, err := c.store.Member().GetMembersByEmails(idStrArr) //获取用户列表*/
		if err != nil {
			_ = level.Error(c.logger).Log("DistributeProclaim", "c.member.GetMembersByEmails", "err", err.Error())
		}
		num, err := c.CreateProclaimMember(ml, dat)
		if err != nil {
			_ = level.Error(c.logger).Log("DistributeProclaim", "c.CreateProclaimMember", "err", err.Error())
		}
		_ = level.Info(c.logger).Log("DistributeProclaim", "c.CreateProclaimMember", "num", num, "case", "user")
	}

	return
}

//公告-用户关系
func (c *serviceProclaim) CreateProclaimMember(ml []types.Member, notice *types.Notices) (n int, err error) {
	data := []*types.NoticeMember{}
	toCc := []string{}

	//rl,err := models.GetNoticeReceiveByAction(notice.Action)  //获取订阅此action通知的用户列表
	admin, _ := c.store.Member().GetInfoById(int64(notice.MemberID))
	toUser := []string{admin.Email} //创建者为主接收者

	for _, v := range ml {
		if v.State == 2 { //如果用户状态已关闭，则不再发送消息
			continue
		}
		d := types.NoticeMember{}
		d.MemberID = v.ID
		d.NoticeID = int64(notice.ID)

		data = append(data, &d) //公告默认选中站内信,邮件，其它方式不发送

		if admin.Email != v.Email {
			toCc = append(toCc, v.Email) //接收邮件的组
		}
	}

	//用户-公告关系入库
	n, err = c.store.NoticeMember().InsertMulti(data)

	//发送公告邮件至接收人
	go func() {
		err = c.SendProclaimEmail(toUser, toCc, *notice)
		if err != nil {
			_ = level.Error(c.logger).Log("DistributeProclaim", "c.SendProclaimEmail", "err", err.Error())
			return
		}
	}()
	return
}

//公告-邮件
func (c *serviceProclaim) SendProclaimEmail(toUser, toCc []string, notice types.Notices) (err error) {
	c.mailClient.SetTitle(notice.Title)
	c.mailClient.SetContent(notice.Content)
	c.mailClient.AddEmailAddress(toUser) //toUser
	c.mailClient.AddCcEmailAddress(toCc) //toCc
	c.mailClient.SetContentType("text/html;charset=utf-8")
	err = c.mailClient.Send()

	if err != nil {
		_ = level.Error(c.logger).Log("DistributeProclaim", "SendProclaimEmail", "err", err.Error())
	}

	return
}
