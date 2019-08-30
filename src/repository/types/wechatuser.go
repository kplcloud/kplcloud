/**
 * @Time : 2019-07-12 10:34 
 * @Author : soupzhb@gmail.com
 * @File : wechatuser.go
 * @Software: GoLand
 */

package types

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type WechatUser struct {
	City          string    `gorm:"column:city" json:"city"`
	Country       string    `gorm:"column:country" json:"country"`
	CreatedAt     null.Time `gorm:"column:created_at" json:"created_at"`
	Headimgurl    string    `gorm:"column:headimgurl" json:"headimgurl"`
	ID            int       `gorm:"column:id;primary_key" json:"id"`
	Nickname      string    `gorm:"column:nickname" json:"nickname"`
	Openid        string    `gorm:"column:openid" json:"openid"`
	Province      string    `gorm:"column:province" json:"province"`
	Remark        string    `gorm:"column:remark" json:"remark"`
	Sex           int       `gorm:"column:sex" json:"sex"`
	Subscribe     int       `gorm:"column:subscribe" json:"subscribe"`
	SubscribeTime time.Time `gorm:"column:subscribe_time" json:"subscribe_time"`
	UpdatedAt     null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (w *WechatUser) TableName() string {
	return "wechat_users"
}
