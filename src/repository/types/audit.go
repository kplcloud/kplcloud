/**
 * @Time : 2021/9/8 5:10 PM
 * @Author : solacowa@gmail.com
 * @File : audit
 * @Software: GoLand
 */

package types

import "time"

type AuditStatus string

const (
	AuditStatusSuccess AuditStatus = "SUCCESS"
	AuditStatusFailed  AuditStatus = "FAILED"
)

// Audit 审计表
type Audit struct {
	Id           int64       `gorm:"column:id;primary_key" json:"id"`
	ClusterId    int64       `gorm:"column:cluster_id;null;comment:'集群ID'" json:"cluster_id"`
	Namespace    string      `gorm:"column:namespace;null;comment:'空间'" json:"namespace"`
	Name         string      `gorm:"column:name;null;comment:'操作目标名称'" json:"name"`
	UserId       int64       `gorm:"column:user_id;notnull;index;comment:'Action'" json:"user_id"`
	PermissionId int64       `gorm:"column:permission_id;notnull;comment:'权限ID'" json:"permission_id"`
	Request      string      `gorm:"column:request;null;comment:'Request'" json:"request"`
	Response     string      `gorm:"column:response;null;comment:'Response'" json:"response"`
	Headers      string      `gorm:"column:headers;null;comment:'Headers'" json:"headers"`
	TimeSince    string      `gorm:"column:time_since;null;comment:'用时'" json:"time_since"`
	Status       AuditStatus `gorm:"column:status;null;comment:'状态'" json:"status"`
	Remark       string      `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	Url          string      `gorm:"column:url;null;comment:'URI'" json:"url"`
	TraceId      string      `gorm:"column:trace_id;null;comment:'TraceId'" json:"trace_id"`

	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName set table
func (*Audit) TableName() string {
	return "audits"
}
