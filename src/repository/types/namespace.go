package types

import (
	"time"
)

type Namespace struct {
	Id        int64      `gorm:"column:id;primary_key" json:"id"`
	ClusterId int64      `gorm:"column:cluster_id;comment:'集群ID'" json:"cluster_id"`
	Alias     string     `gorm:"column:alias;comment:'别名'" json:"alias"`
	Name      string     `gorm:"column:name;index;comment:'名称'" json:"name"`
	Remark    string     `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	Status    string     `gorm:"column:status;comment:'状态'" json:"status"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (*Namespace) TableName() string {
	return "namespace"
}
