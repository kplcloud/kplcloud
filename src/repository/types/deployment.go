/**
 * @Time : 8/16/21 4:50 PM
 * @Author : solacowa@gmail.com
 * @File : deployment
 * @Software: GoLand
 */

package types

import "time"

type Deployment struct {
	Id        int64      `gorm:"column:id;rimary_key" json:"id"`
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	Data      string     `json:"data"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
