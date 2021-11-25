/**
 * @Time : 2021/11/25 9:46 AM
 * @Author : solacowa@gmail.com
 * @File : sys_notice
 * @Software: GoLand
 */

package types

import "time"

type PushType string
type NoticeType string
type SubType string

const (
	PushTypeAnnouncement PushType = "announcement" // 公告
)

// Subscription 通知、告警订阅表
type Subscription struct {
	Id        int64      `gorm:"column:id;primary_key" json:"id"`
	UserId    int64      `json:"user_id"`
	SubType   SubType    `json:"sub_type"`
	Email     bool       `json:"email"`
	Site      bool       `json:"site"`
	Wechat    bool       `json:"wechat"`
	Sms       bool       `json:"sms"`
	DingDing  bool       `json:"ding_ding"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// NotifyType 通知类型
type NotifyType string

const (
	NotifyTypeDelete    NotifyType = "Delete"    // 删除
	NotifyTypeApply     NotifyType = "Apply"     // 审核通过
	NotifyTypeBuild     NotifyType = "Build"     // 构建
	NotifyTypeAudit     NotifyType = "Audit"     // 提交审核
	NotifyTypeRollback  NotifyType = "Rollback"  // 回滚
	NotifyTypeReboot    NotifyType = "Reboot"    // 重启
	NotifyTypeExpansion NotifyType = "Expansion" // 扩容
	NotifyTypeExtend    NotifyType = "Extend"    // 伸缩
	// 存储
	// 日志
	// 网关
	// 调整启动命令
	// ReadinessProbe 探针
)

// Notify 通知表
// 通知可以是定身的,如果选择定身类型则需要填useId或ns或name等等
type Notify struct {
	Id         int64 `gorm:"column:id;primary_key" json:"id"`
	Title      string
	Content    string
	NotifyType NotifyType `json:"notify_type"`
	Namespace  string
	Name       string
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// Alarm 告警
type Alarm struct {
	Id        int64
	Title     string
	Content   string
	Namespace string
	Name      string
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// Information 消息关联关系表
type Information struct {
	Id          int64     `gorm:"column:id;primary_key" json:"id"`
	UserId      int64     // 用户ID
	RelationKey string    // 关联类型 公告、通知、告警
	RelationId  int64     // 关联id
	Readed      bool      // 是否已读
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
}

// SysAnnouncement 系统公告结构
// 默认information表没有该关联，也就是用户未读
// 用户登陆后默认取相应时间内的所有公告
// 然后用noticeId取当前用户的information表数据取交集就能得到哪些是已读
// 已读不弹窗，未读弹窗
type SysAnnouncement struct {
	Id             int64      `gorm:"column:id;primary_key" json:"id"`
	Title          string     `gorm:"column:title;size:255;null;comment:'标题'" json:"title"`
	Content        int64      `gorm:"column:content;null;type:text;size:10000;comment:'内容'" json:"content"`
	PublisherId    int64      `gorm:"column:publisher_id;notnull;comment:'发布人ID'" json:"publisher_id"`
	PushType       PushType   // 所有
	NoticeType     NoticeType // 公告类型, 如平台、资源、安全、财务、告警
	Namespace      string     `json:"namespace"`
	Name           string     `json:"name"`
	PublicityBegin *time.Time `gorm:"column:publicity_begin;comment:'公示开始时间'" json:"publicity_begin"`
	PublicityEnd   *time.Time `gorm:"column:publicity_end;comment:'公示结束时间'" json:"publicity_end"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	Publisher SysUser `gorm:"foreignkey:Id;references:PublisherId"`
}

// TableName sets the insert table name for this struct type
func (p *SysAnnouncement) TableName() string {
	return "sys_announcement"
}
