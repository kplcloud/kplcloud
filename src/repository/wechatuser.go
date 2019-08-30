/**
 * @Time : 2019-07-12 10:39
 * @Author : soupzhb@gmail.com
 * @File : wechatuser.go
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type WechatUserRepository interface {
	FindByOpenid(openid string) (v *types.WechatUser, err error)
	CreateOrUpdate(n *types.WechatUser) error
	UnSubscribe(openid string) error
}

type wechatUser struct {
	db *gorm.DB
}

func NewWechatUserRepository(db *gorm.DB) WechatUserRepository {
	return &wechatUser{db: db}
}

/**
 * @Title 创建或更新微信用户信息
 */
func (c *wechatUser) CreateOrUpdate(wu *types.WechatUser) error {
	wxInfo, err := c.FindByOpenid(wu.Openid)
	if err != nil {
		//insert
		return c.db.Save(wu).Error
	}
	if wxInfo.ID > 0 {
		//update
		var w types.WechatUser
		return c.db.Model(&w).Where("openid = ?", wu.Openid).Updates(wu).Error
	}
	return err
}

/**
 * @Title 获取详情
 */
func (c *wechatUser) FindByOpenid(openid string) (v *types.WechatUser, err error) {
	var temp types.WechatUser
	if err = c.db.Where("openid = ?", openid).First(&temp).Error; err != nil {
		return
	}
	return &temp, nil
}

/**
 * @Title 取消关注
 */
func (c *wechatUser) UnSubscribe(openid string) error {
	var w types.WechatUser
	return c.db.Model(&w).Where("openid = ?", openid).Update("subscribe", 0).Error
}
