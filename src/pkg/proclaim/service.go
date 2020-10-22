package proclaim

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	amqpClient "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	"github.com/kplcloud/kplcloud/src/util/upload"
	"strings"
)

var (
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrProclaimParamsRefused = errors.New("参数校验未通过.")
	ErrProclaimList          = errors.New("获取公告列表失败")
	ErrProclaimListCount     = errors.New("获取公告列表总数失败")
)

type Service interface {
	Get(ctx context.Context, id int) (resp *types.Notices, err error)
	Post(ctx context.Context, req proclaimRequest) (err error)
	List(ctx context.Context, name string, page int, limit int) (res map[string]interface{}, err error)
	ContentHtmlHandle(content string) (body string, err error)
}

type service struct {
	logger     log.Logger
	cf         *config.Config
	amqpClient amqpClient.AmqpClient
	store      repository.Repository
}

func NewService(logger log.Logger, cf *config.Config, amqpClient amqpClient.AmqpClient, store repository.Repository) Service {
	return &service{
		logger:     logger,
		cf:         cf,
		amqpClient: amqpClient,
		store:      store,
	}
}

func (c *service) Get(ctx context.Context, id int) (resp *types.Notices, err error) {
	return c.store.Proclaim().FindById(id)
}

func (c *service) Post(ctx context.Context, req proclaimRequest) (err error) {
	memberId := ctx.Value(middleware.UserIdContext).(int64)
	_ = level.Info(c.logger).Log("title", req.Title, "content", req.Content, "proclaim_type", req.ProclaimType)
	if req.Title == "" || req.Content == "" {
		return ErrProclaimParamsRefused
	}

	var proclaimReceive string
	if req.ProclaimType == "namespace" {
		for i := 0; i < len(req.Namespace); i++ {
			proclaimReceive += req.Namespace[i] + ","
		}
	} else if req.ProclaimType == "user" {
		for i := 0; i < len(req.UserList); i++ {
			proclaimReceive += req.UserList[i] + ","
		}
	}

	content, err := c.ContentHtmlHandle(req.Content)
	if err != nil {
		_ = level.Error(c.logger).Log("proclaim", "post", "c.ContentHtmlHandle.Err", err.Error())
		return
	}

	data := new(types.Notices)
	data.Title = req.Title
	data.Content = content
	data.Type = 1
	data.MemberID = int(memberId)
	data.ProclaimType = req.ProclaimType
	data.ProclaimReceive = proclaimReceive

	_ = c.store.Proclaim().CreateReturnId(data)

	b, _ := json.Marshal(data)

	//公告内容放队列，推送邮件用
	defer func() {
		if err := c.amqpClient.PublishOnQueue(amqpClient.ProclaimTopic, func() []byte {
			return []byte(b)
		}); err != nil {
			_ = level.Error(c.logger).Log("amqpClient", "PublicProclaimQueue", "err", err.Error())
		}
	}()

	return err
}

func (c *service) ContentHtmlHandle(content string) (body string, err error) {

	type uploadConfig struct {
		UploadPath   string `yaml:"upload_path"`
		DomainPrefix string `yaml:"domain_prefix"`
	}

	var upConf uploadConfig
	upConf.UploadPath = c.cf.GetString("server", "upload_path")
	upConf.DomainPrefix = c.cf.GetString("server", "domain")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	doc.Find(".image-wrap").Each(func(i int, s *goquery.Selection) {
		//获取src内容
		src, _ := s.Find("img").Attr("src")

		if strings.Contains(src, "base64") { //如果是引用网络图片，则不替换
			srcArr := strings.Split(src, ",")
			postfix := strings.Split(strings.Split(srcArr[0], "/")[1], ";")[0]
			base64Str := srcArr[1]

			var im upload.ImgInfo
			im.Base64 = base64Str
			im.Type = postfix
			im.Path = upConf.UploadPath + "images"

			re, err := im.Base64ToFile()
			if err != nil {
				_ = level.Error(c.logger).Log("ContentHtmlHandle.Base64ToFile.Err", err.Error())
				return
			}

			//图片url
			imgUrl := upConf.DomainPrefix + "images" + re.PathSmall
			s.Find("img").SetAttr("src", imgUrl) //替换src
		}

	})

	body, err = doc.Find("body").Html()
	if err != nil {
		_ = level.Error(c.logger).Log("doc", "Find", "Html", "body", "err", err.Error())
	}

	_ = level.Info(c.logger).Log("ContentHtmlHandle.ContentToHtml", body)

	return
}

func (c *service) List(ctx context.Context, name string, page int, limit int) (res map[string]interface{}, err error) {
	count, err := c.store.Proclaim().Count(name, 1)
	if err != nil {
		_ = level.Error(c.logger).Log("template", "Count", "err", err.Error())
		return nil, ErrProclaimListCount
	}

	p := paginator.NewPaginator(page, limit, count)

	list, err := c.store.Proclaim().FindOffsetLimit(name, 1, p.Offset(), limit)

	if err != nil {
		_ = level.Error(c.logger).Log("proclaim", "FindOffsetLimit", "err", err.Error())
		return nil, ErrProclaimList
	}

	//用户列表
	ml, err := c.store.Member().GetMembersAll()
	//业务线列表
	nsl, err := c.store.Namespace().FindAll()

	var listMap []interface{}
	for _, val := range list {
		var prText, memberName, conText string
		if val.ProclaimType == "namespace" {

			nsArr := strings.Split(val.ProclaimReceive, ",")
			for _, n := range nsl {
				for _, ns := range nsArr {
					if n.Name == ns {
						prText += n.Name + ","
					}
				}
			}
			//截取最后一个,
			if len(prText) > 0 {
				prText = prText[0 : len(prText)-1]
			}
		} else if val.ProclaimType == "user" {

			uidArr := strings.Split(val.ProclaimReceive, ",")

			for _, u := range ml {
				for _, uid := range uidArr {
					if u.Email == uid {
						prText += u.Username + ","
					}
				}
			}
			//截取最后一个,
			if len(prText) > 0 {
				prText = prText[0 : len(prText)-1]
			}
		} else {
			prText = "全部"
		}

		for _, u := range ml {
			if int(u.ID) == val.MemberID {
				memberName = u.Username
			}
		}

		c := val.Content
		if len(c) > 25 {
			conText = string([]rune(c)[:25]) + "..."
		} else {
			conText = c
		}

		data := map[string]interface{}{
			"id":                    val.ID,
			"title":                 val.Title,
			"content":               conText,
			"member_id":             val.MemberID,
			"member_name":           memberName,
			"proclaim_type":         val.ProclaimType,
			"proclaim_receive":      val.ProclaimReceive,
			"proclaim_receive_text": prText,
			"created_at":            val.CreatedAt.Time.Format("2006-01-02 15:04:05"),
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
