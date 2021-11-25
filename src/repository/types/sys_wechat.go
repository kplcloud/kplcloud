/**
 * @Time : 2019-07-12 10:34
 * @Author : soupzhb@gmail.com
 * @File : wechatuser.go
 * @Software: GoLand
 */

package types

import (
	"time"
)

type SysWechat struct {
	ID            int        `gorm:"column:id;primary_key" json:"id"`
	UserId        int64      `json:"user_id"`
	City          string     `gorm:"column:city" json:"city"`
	Country       string     `gorm:"column:country" json:"country"`
	HeadImgUrl    string     `gorm:"column:head_img_url" json:"headimgurl"`
	Nickname      string     `gorm:"column:nickname" json:"nickname"`
	Openid        string     `gorm:"column:openid" json:"openid"`
	Province      string     `gorm:"column:province" json:"province"`
	Remark        string     `gorm:"column:remark" json:"remark"`
	Sex           int        `gorm:"column:sex" json:"sex"`
	Subscribe     int        `gorm:"column:subscribe" json:"subscribe"`
	SubscribeTime time.Time  `gorm:"column:subscribe_time" json:"subscribe_time"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (w *SysWechat) TableName() string {
	return "sys_wechat"
}
