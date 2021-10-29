/**
 * @Time : 2019-07-02 10:41
 * @Author : soupzhb@gmail.com
 * @File : service.go
 * @Software: GoLand
 */

package notice

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	amqpClient "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/event"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/msgs"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	"strings"
)

var (
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrNoticeList            = errors.New("获取消息列表失败")
	ErrNoticeListCount       = errors.New("获取消息列表总数失败")
	ErrNoticeMemberListCount = errors.New("获取用户阅读消息列表总数失败")
	ErrNoticeNoAppName       = errors.New("项目名称不能为空")
	ErrNoticeNoNamespace     = errors.New("命名空间不能为空")
)

type Service interface {
	Create(ctx context.Context, req *event.WebhooksRequest) (err error)
	List(ctx context.Context, param map[string]string, page int, limit int) (res map[string]interface{}, err error)
	Tips(ctx context.Context) (res interface{}, err error)
	CountRead(ctx context.Context, param map[string]string) (res map[string]interface{}, err error)
	ClearAll(ctx context.Context, noticeType int) (err error)
	Detail(ctx context.Context, noticeId int) (res interface{}, err error)
}

type service struct {
	logger     log.Logger
	config     *config.Config
	amqpClient amqpClient.AmqpClient
	store      repository.Repository
}

func NewService(logger log.Logger, cf *config.Config, amqpClient amqpClient.AmqpClient, store repository.Repository) Service {
	return &service{
		logger:     logger,
		config:     cf,
		amqpClient: amqpClient,
		store:      store,
	}
}

func (c *service) Create(ctx context.Context, req *event.WebhooksRequest) (err error) {
	if req.AppName == "" {
		return ErrNoticeNoAppName
	}

	if req.Namespace == "" {
		return ErrNoticeNoNamespace
	}

	//消息组装
	var msgText string
	if req.Title == "" {
		msgText = "项目名称" + "【" + req.AppName + "】" + "命名空间" + "【" + req.Namespace + "】" + "已操作" + "【" + req.EventDesc + "】"
	} else {
		msgText = req.Title
	}

	//db
	data := new(types.Notices)
	data.Title = msgText
	data.Content = req.Message
	data.Type = 2
	data.MemberID = int(req.Member.ID)
	data.Action = req.Event.String()
	data.Namespace = req.Namespace
	data.Name = req.AppName

	_ = c.store.Notice().CreateReturnId(data)

	noticeMqData := new(msgs.NoticeMqData)
	noticeMqData.WebHooksReq = *req
	noticeMqData.Notice = *data

	b, _ := json.Marshal(noticeMqData)

	//通知内容放队列，推送邮件用
	defer func() {
		if err := c.amqpClient.PublishOnQueue(amqpClient.NoticeTopic, func() []byte {
			return []byte(b)
		}); err != nil {
			_ = level.Warn(c.logger).Log("amqpClient", "PublicNoticeQueue", "err", err.Error())
		}
	}()
	return
}

func (c *service) List(ctx context.Context, param map[string]string, page int, limit int) (res map[string]interface{}, err error) {

	memberId := ctx.Value(middleware.UserIdContext).(int64)
	count, err := c.store.NoticeMember().CountMessage(param, memberId)
	if err != nil {
		_ = level.Error(c.logger).Log("notice", "CountMessage", "err", err.Error())
		return nil, ErrNoticeListCount
	}

	p := paginator.NewPaginator(page, limit, count)

	list, err := c.store.NoticeMember().FindMessageLimit(param, memberId, p.Offset(), limit)
	if err != nil {
		_ = level.Error(c.logger).Log("notice", "FindMessageLimit", "err", err.Error())
		return nil, ErrNoticeList
	}

	//用户列表
	ml, err := c.store.Member().GetMembersAll()

	var listMap []interface{}
	for _, val := range list {

		var memberName string
		for _, u := range ml {
			if u.ID == val.MemberId {
				memberName = u.Username
			}
		}

		conText := val.Content
		if val.Type == 2 {
			conText = strings.Replace(conText, "\n", "<br/>", -1)

		}

		data := map[string]interface{}{
			"id":            val.Id,
			"title":         val.Title,
			"content":       conText,
			"action":        val.Action,
			"name":          val.Name,
			"namespace":     val.Namespace,
			"member_id":     val.MemberId,
			"member_name":   memberName,
			"proclaim_type": val.ProclaimType,
			"type":          val.Type,
			"is_read":       val.IsRead,
			"created_at":    val.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		}
		listMap = append(listMap, data)
	}

	res = map[string]interface{}{
		"list": listMap,
		"page": map[string]interface{}{
			"total":     count,
			"pageTotal": p.PageTotal(),
			"pageSize":  limit,
			"page":      p.Page(),
		},
	}
	return
}

func (c *service) Tips(ctx context.Context) (res interface{}, err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	//类型=1公告；=2通知；=3告警  每组取前100条，用于用户中心右上角通知提醒
	var myNotices [][]*repository.MyNotice
	myNotices1, err := c.store.NoticeMember().FindMessageLimit(map[string]string{"type": "1", "is_read": "0"}, memberId, 0, 100)
	myNotices2, err := c.store.NoticeMember().FindMessageLimit(map[string]string{"type": "2", "is_read": "0"}, memberId, 0, 100)
	myNotices3, err := c.store.NoticeMember().FindMessageLimit(map[string]string{"type": "3", "is_read": "0"}, memberId, 0, 100)

	myNotices = append(myNotices, myNotices1)
	myNotices = append(myNotices, myNotices2)
	myNotices = append(myNotices, myNotices3)

	if err != nil {
		_ = level.Error(c.logger).Log("notice.Tips", "FindMessageLimit", "err", err.Error())
		return nil, ErrNoticeList
	}

	icons := map[string]string{
		"Alarm":          "http://source.qiniu.cnd.nsini.com/images/2019/08/f7/8c/61/20190827-5ee43b38986a144d6b5022ea8c8f748f.jpeg",
		"Mail":           "https://gw.alipayobjects.com/zos/rmsportal/ThXAXghbEsBCCSDihZxY.png",
		"Audit":          "https://gw.alipayobjects.com/zos/rmsportal/kISTdvpyTAhtGxpovNWd.png",
		"Delete":         "https://gw.alipayobjects.com/zos/rmsportal/GvqBnKhFgObvnSGkDsje.png",
		"Build":          "https://jenkins.io/sites/default/files/jenkins_favicon.ico",
		"Apply":          "https://gw.alipayobjects.com/zos/rmsportal/ThXAXghbEsBCCSDihZxY.png",
		"Member":         "https://niu.yirendai.com/kpl-logo-blue.png",
		"Rollback":       "http://niu.yirendai.com/clock-event-history-schedule-time-icon--19.png",
		"Logging":        "http://niu.yirendai.com/kpl-logging.png",
		"Reboot":         "http://niu.yirendai.com/kpl-reboot.png",
		"Command":        "http://niu.yirendai.com/kpl-command.png",
		"Storage":        "http://niu.yirendai.com/kpl-storage.png",
		"Gateway":        "http://niu.yirendai.com/kpl-gateway.png",
		"Expansion":      "http://niu.yirendai.com/kpl-expansion.png",      //扩容
		"Extend":         "http://niu.yirendai.com/kpl-extend.png",         //伸缩
		"SwitchModel":    "http://niu.yirendai.com/kpl-switchmodel.png",    //调整服务模式
		"ReadinessProbe": "http://niu.yirendai.com/kpl-readinessprobe.png", //修改探针
		"Test":           "http://niu.yirendai.com/kpl-test.png",
	}

	var notices []*noticeInfo

	for _, curNotices := range myNotices {
		for _, v := range curNotices {
			var t = "公告"
			if v.Type == 2 {
				t = "通知"
			}
			if v.Type == 3 {
				t = "告警"
			}
			notices = append(notices, &noticeInfo{
				Id:       int64(v.Id),
				Avatar:   icons[v.Action],
				Title:    v.Title,
				Datetime: v.CreatedAt,
				Type:     t,
			})
		}
	}

	res = notices
	return
}

func (c *service) CountRead(ctx context.Context, param map[string]string) (res map[string]interface{}, err error) {

	memberId := ctx.Value(middleware.UserIdContext).(int64)
	param["is_read"] = "0"
	countUnRead, err := c.store.NoticeMember().CountRead(param, memberId)
	if err != nil {
		_ = level.Error(c.logger).Log("notice", "CountRead", "err", err.Error())
		return nil, ErrNoticeMemberListCount
	}
	param["is_read"] = "1"
	countRead, err := c.store.NoticeMember().CountRead(param, memberId)
	if err != nil {
		_ = level.Error(c.logger).Log("notice", "CountRead", "err", err.Error())
		return nil, ErrNoticeMemberListCount
	}

	res = map[string]interface{}{
		"read":   countRead,
		"unread": countUnRead,
	}

	return
}

func (c *service) ClearAll(ctx context.Context, noticeType int) (err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	err = c.store.NoticeMember().ClearAll(noticeType, memberId)
	return
}

func (c *service) Detail(ctx context.Context, noticeMemberId int) (res interface{}, err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	res, err = c.store.NoticeMember().Detail(noticeMemberId, memberId)
	_ = c.store.NoticeMember().HasRead(int64(noticeMemberId))
	return
}
