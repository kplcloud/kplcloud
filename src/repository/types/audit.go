/**
 * @Time : 2021/9/8 5:10 PM
 * @Author : solacowa@gmail.com
 * @File : audit
 * @Software: GoLand
 */

package types

import "time"

// Audit is the operation recoder
type Audit struct {
	Id        int64  `gorm:"column:id;primary_key" json:"id"`
	ClusterId string `gorm:"column:cluster_id;comment:'集群ID'" json:"cluster_id"`
	Namespace string `gorm:"column:namespace;null;comment:'空间'" json:"namespace"`
	Name      string `gorm:"column:name;null;comment:'操作目标名称'" json:"name"`
	Method    string `gorm:"column:method;null;comment:'Method'" json:"method"`
	Action    string `gorm:"column:action;null;comment:'Action'" json:"action"`
	UserId    int64  `gorm:"column:user_id;notnull;comment:'Action'" json:"user_id"`
	TimeSince string `gorm:"column:time_since;null;comment:'用时'" json:"time_since"`
	Status    int    `gorm:"column:status;null;comment:'状态'" json:"status"` //1 success 2 faild
	Message   string `gorm:"column:message;null;comment:'消息'" json:"message"`
	Request   string `gorm:"column:request;null;comment:'Request'" json:"request"`
	Response  string `gorm:"column:response;null;comment:'Response'" json:"response"`
	Headers   string `gorm:"column:headers;null;comment:'Headers'" json:"headers"`
	Url       string `gorm:"column:url;null;comment:'URI'" json:"url"`

	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName set table
func (*Audit) TableName() string {
	return "audits"
}
