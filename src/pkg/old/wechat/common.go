/**
 * @Time : 2019-08-27 09:47
 * @Author : soupzhb@gmail.com
 * @File : common.go
 * @Software: GoLand
 */

package wechat

import (
	"github.com/kplcloud/kplcloud/src/util/mp"
	"github.com/yijizhichang/wechat-sdk/mp/message"
	"github.com/yijizhichang/wechat-sdk/mp/message/callback/response"
)

func ResponseWechatByKey(key string) (msgType message.MsgType, reStr interface{}) {
	//默认
	msgType = message.MsgTypeText
	reStr = ""

	//文本
	if v, ok := mp.TextConfig[key]; ok {
		msgType = message.MsgTypeText
		reStr = response.NewText(v.Content)
	}

	//文章
	if v, ok := mp.ArticleConfig[key]; ok {
		msgType = message.MsgTypeNews
		ar := response.NewArticle(v.Title, v.Description, v.PicURL, v.URL)
		var newsList []*response.Article
		newsList = append(newsList, ar)
		reStr = response.NewNews(newsList)
	}

	return
}
