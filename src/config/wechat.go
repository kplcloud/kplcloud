/**
 * @Time : 2019-07-09 18:29
 * @Author : soupzhb@gmail.com
 * @File : wechat.go
 * @Software: GoLand
 */

package config

import (
	"github.com/icowan/config"
	"github.com/yijizhichang/wechat-sdk"
)

var WechatClient *wechat.Wechat

func GetWechatConfig(cf *config.Config) *wechat.Config {
	//微信配置
	configWechat := &wechat.Config{
		AppID:            cf.GetString("wechat", "app_id"),           //开发者ID(AppID)
		AppSecret:        cf.GetString("wechat", "app_secret"),       //开发者密码AppSecret
		Token:            cf.GetString("wechat", "token"),            //令牌(Token)
		EncodingAESKey:   cf.GetString("wechat", "encoding_aes_key"), //消息加解密密钥 EncodingAESKey
		PayMchId:         "",                                         //支付 - 商户 ID
		PayNotifyUrl:     "",                                         //支付 - 接受微信支付结果通知的接口地址
		PayKey:           "",                                         //支付 - 商户后台设置的支付 key
		Cache:            nil,                                        //缓存方式 默认为file，可选 file,redis,redisCluster
		ThirdAccessToken: false,                                      //是否使用第三方accessToken
		ProxyUrl:         "",                                         //代理地址
	}

	return configWechat
}
