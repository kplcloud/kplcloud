/**
 * @Time : 2019-07-09 18:51
 * @Author : soupzhb@gmail.com
 * @File : service.go
 * @Software: GoLand
 */

package wechat

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/yijizhichang/wechat-sdk"
	"github.com/yijizhichang/wechat-sdk/mp/menu"
	"github.com/yijizhichang/wechat-sdk/mp/message"
	"github.com/yijizhichang/wechat-sdk/mp/message/callback/request"
	"github.com/yijizhichang/wechat-sdk/mp/message/callback/response"
	"github.com/yijizhichang/wechat-sdk/mp/message/template"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrCreateQr        = errors.New("微信生成二维码失败")
)

type Service interface {
	Receive(ctx context.Context) (str, contentType string, err error)
	GetQr(ctx context.Context, req qrRequest) (interface{}, error)
	TestSend(ctx context.Context) (interface{}, error)
	Menu(ctx context.Context) (interface{}, error)
}

type service struct {
	logger   log.Logger
	config   *config.Config
	wxClient *wechat.Wechat
	store    repository.Repository
}

func NewService(logger log.Logger, cf *config.Config, wx *wechat.Wechat, store repository.Repository) Service {
	return &service{
		logger:   logger,
		wxClient: wx,
		config:   cf,
		store:    store,
	}
}

func (c *service) Receive(ctx context.Context) (str, contentType string, err error) {
	req := ctx.Value(httpRequestContext).(*http.Request)

	// 传入request和responseWriter
	server := c.wxClient.GetResponseServer(req)

	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *response.Reply {
		var reStr interface{}
		var msgType message.MsgType
		var msgInfo interface{} //接收微信消息

		msgType = message.MsgTypeNothing
		reStr = ""

		//根据微信回调时的消息类型，来相应获取对应消息明细
		switch msg.MsgCommon.MsgType {
		//消息类型
		case "text":
			msgInfo = request.GetText(&msg)

			//根据业务需求,被动回复微信消息
			switch msg.Content {
			case "1":
				reStr = response.NewText("回复测试文件本消息")
				//reStr = models.ResponseWechatByKey("clickContactKefu")
				msgType = message.MsgTypeText
			case "2":
				reStr = response.NewImage("9999999999")
				msgType = message.MsgTypeImage
			case "3":
				reStr = response.NewVoice("9999999999")
				msgType = message.MsgTypeVoice
			case "4":
				reStr = response.NewVideo("999999999", "视频", "我是一条视频信息")
				msgType = message.MsgTypeVideo
			case "5":
				ar := response.NewArticle("图文消息", "我是一条图文消息", "https://www.baidu.com/img/bd_logo1.png", "https://www.baidu.com/")
				var newsList []*response.Article
				newsList = append(newsList, ar)
				reStr = response.NewNews(newsList)
				//fmt.Println("图文消息：", reStr)
				msgType = message.MsgTypeNews
			default:
				//reStr = ""
				//msgType = message.MsgTypeNothing
			}
		case "image":
			msgInfo = request.GetImage(&msg)
		case "voice":
			msgInfo = request.GetVoice(&msg)
		case "video":
			msgInfo = request.GetVideo(&msg)
		case "shortvideo":
			msgInfo = request.GetShortVideo(&msg)
		case "location":
			msgInfo = request.GetLocation(&msg)
		case "link":
			msgInfo = request.GetLink(&msg)
			//事件类型
		case "event":
			switch msg.Event {
			case "subscribe":
				InfoSubscribeEvent := request.GetSubscribeEvent(&msg)
				uTag := c.wxClient.GetUser()
				re, err := uTag.GetUserInfo(InfoSubscribeEvent.FromUserName, "zh_CN")
				if err != nil {
					_ = level.Error(c.logger).Log("wechat.Receive", "uTag.GetUserInfo", "err", err.Error(), "openid", InfoSubscribeEvent.FromUserName)
				} else {
					//wechatUser to db
					wu := new(types.WechatUser)
					wu.City = re.City
					wu.Country = re.Country
					wu.Headimgurl = re.Headimgurl
					wu.Nickname = re.Nickname
					wu.Openid = re.Openid
					wu.Province = re.Province
					wu.Remark = re.Remark
					wu.Sex = re.Sex
					wu.Subscribe = 1
					wu.SubscribeTime = time.Unix(int64(re.SubscribeTime), 0)

					err := c.store.WechatUser().CreateOrUpdate(wu)
					if err != nil {
						_ = level.Error(c.logger).Log("wechat.Receive", "c.wechatUser.Create", "err", err.Error())
					}
				}

				//默认关注欢迎语
				msgType, reStr = ResponseWechatByKey("welcomText")

				//绑定扫码关注
				openId := InfoSubscribeEvent.FromUserName
				eventKey := InfoSubscribeEvent.EventKey
				if strings.Index(eventKey, "bindEmail") != -1 {
					ekArr := strings.Split(eventKey, ":")
					email := ekArr[1]
					err := c.store.Member().BindWechat(email, openId)
					if err != nil {
						_ = level.Error(c.logger).Log("wechat.Receive", "c.member.BindWechat", "err", err.Error())
						msgType = message.MsgTypeText
						reStr = response.NewText(email + "\n 绑定微信失败")
					} else {
						msgType = message.MsgTypeText
						reStr = response.NewText(email + "\n 绑定微信成功")
					}
				}

			case "unsubscribe":
				InfoUnSubscribeEvent := request.GetUnsubscribeEvent(&msg)
				err := c.store.WechatUser().UnSubscribe(InfoUnSubscribeEvent.FromUserName)
				if err != nil {
					_ = level.Error(c.logger).Log("wechat.Receive", "c.wechatUser.UnSubscribe", "err", err.Error())
				}

			case "SCAN":
				InfoScanEvent := request.GetScanEvent(&msg)
				openId := InfoScanEvent.FromUserName
				eventKey := InfoScanEvent.EventKey
				if strings.Index(eventKey, "bindEmail") != -1 {
					ekArr := strings.Split(eventKey, ":")
					email := ekArr[1]
					err := c.store.Member().BindWechat(email, openId)
					if err != nil {
						_ = level.Error(c.logger).Log("wechat.Receive", "c.member.BindWechat", "err", err.Error())
						msgType = message.MsgTypeText
						reStr = response.NewText(email + "\n 绑定微信失败")
					} else {
						msgType = message.MsgTypeText
						reStr = response.NewText(email + "\n 绑定微信成功")
					}
				}

			case "CLICK", "VIEW":
				InfoMenuEvent := request.GetMenuEvent(&msg)
				key := InfoMenuEvent.EventKey
				msgInfo = InfoMenuEvent

				//回复文本/文章消息
				msgType, reStr = ResponseWechatByKey(key)

			case "TEMPLATESENDJOBFINISH":
				msgInfo = request.GetTemplateSendJobFinishEvent(&msg)
			}
		}
		fmt.Println("接收消息：", msgInfo, "回复类型：", msgType, "回复内容：", reStr)

		return &response.Reply{MsgType: msgType, MsgData: reStr}

	})

	//处理消息接收以及回复
	echostr, contentType, echostrExist, err := server.ResponseServe()
	if err != nil {
		_ = level.Error(c.logger).Log("wechat.Receive", "server.ResponseServe", "err", err.Error())
		return "", "", err
	}
	if echostrExist {
		return echostr, contentType, nil
	}
	//发送回复的消息
	data, dataContentType, err := server.ResponseSend()
	return data, dataContentType, nil
}

func (c *service) GetQr(ctx context.Context, req qrRequest) (res interface{}, err error) {
	account := c.wxClient.GetAccount()
	re, err := account.CreateQrCodeSceneStr(false, "bindEmail:"+req.Email, 7200)

	if err != nil {
		_ = level.Error(c.logger).Log("wechat.Receive", "server.GetQr", "err", err.Error())
		return nil, ErrCreateQr
	}

	qrUrl := account.GetQrCodeUrl(re.Ticket)

	type reQr struct {
		Ticket string `json:"ticket"`
		Url    string `json:"url"`
		Expire int32  `json:"expire"`
		QrUrl  string `json:"qr_url"`
	}

	res = reQr{re.Ticket, re.Url, re.ExpireSeconds, qrUrl}
	return
}

func (c *service) TestSend(ctx context.Context) (res interface{}, err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	memberInfo, err := c.store.Member().GetInfoById(memberId)

	if err != nil {
		return
	}

	tpl := c.wxClient.GetTemplate()

	tit := "Hi," + memberInfo.Username + "。这是一条来自开普勒云平台的测试消息"
	curTime := time.Now().Format("2006-01-02 15:04")

	msgText := new(template.Message)
	msgText.ToUser = memberInfo.Openid
	msgText.TemplateID = c.config.GetString("wechat", "tpl_test")
	msgText.URL = ""
	msgText.Data = make(map[string]*template.DataItem)
	msgText.Data["first"] = &template.DataItem{tit, "#00CD00"}
	msgText.Data["keyword1"] = &template.DataItem{curTime, "#000"}
	msgText.Data["keyword2"] = &template.DataItem{curTime, "#000"}
	msgText.Data["keyword3"] = &template.DataItem{"成功", "#000"}
	msgText.Data["remark"] = &template.DataItem{"您已成功绑定微信账号，可以通过微信接收订阅消息，感谢您的关注", "#FF8C00"}

	res, err = tpl.Send(msgText)

	if err != nil {
		return
	}

	return
}

func (c *service) Menu(ctx context.Context) (res interface{}, err error) {
	//获取菜单配置方法
	mu := c.wxClient.GetMenu()

	//二级菜单列表
	/*subMenuList1 := menu.SetButton(
		menu.WithClickButton("赞我们一下","V1001_GOOD"),    //不同的菜单类型，调用不用的menu.WithXXXButton()方法
		menu.WithViewButton("搜一下","http://www.soso.com/"),
		menu.WithLocationSelectButton("上报位置","wz2039_fdei"),
		menu.WithMiniprogramButton("跳转小程序","http://mp.weixin.qq.com","wx286b93c14bbf93aa","pages/lunar/index"),
	)*/
	subMenuList2 := menu.SetButton(
		menu.WithClickButton("联系方式", "clickContactUs"),
	)
	subMenuList3 := menu.SetButton(
		menu.WithViewButton("在线文档", "https://docs.nsini.com"),
	)

	//一级菜单列表
	parentMenu1 := menu.SetButton(
		menu.WithClickButton("关于开普勒", "clickAboutUsText"),
		menu.WithSubButton("联系我们", subMenuList2),
		menu.WithSubButton("帮助中心", subMenuList3),
	)

	//创建菜单
	res, err = mu.SetMenu(parentMenu1...)

	if err != nil {
		return
	}

	return

}
