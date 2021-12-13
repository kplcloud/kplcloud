/**
 * @Time : 2021/8/25 9:56 AM
 * @Author : solacowa@gmail.com
 * @File : k8s_template
 * @Software: GoLand
 */

package types

import "time"

type Kind string

const (
	KindClusterRole             Kind = "ClusterRole"
	KindSecret                  Kind = "Secret"
	KindPersistentVolumeClaim   Kind = "PersistentVolumeClaim"
	KindHorizontalPodAutoscaler Kind = "HorizontalPodAutoscaler"
)

func (s Kind) String() string {
	return string(s)
}

type K8sTemplate struct {
	Id int64 `gorm:"column:id;primary_key" json:"id"`
	//ClusterId int64      `gorm:"column:cluster_id;notnull;index;comment:'集群ID'" json:"cluster_id"`
	Kind      Kind       `gorm:"column:kind;notnull;unique;index;size:32;comment:'资源类型'" json:"kind"`
	Alias     string     `gorm:"column:alias;index;size:32;null;comment:'资源名称'" json:"alias"`
	Rules     string     `gorm:"column:rules;null;comment:'规则'" json:"rules"`
	Content   string     `gorm:"column:content;null;type:text;comment:'模版内容'" json:"content"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName set table
func (*K8sTemplate) TableName() string {
	return "k8s_template"
}
