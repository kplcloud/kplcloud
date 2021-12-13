/**
 * @Time : 2021/12/13 5:41 PM
 * @Author : solacowa@gmail.com
 * @File : horizontal_pod_autoscaler
 * @Software: GoLand
 */

package types

import "time"

// HorizontalPodAutoscaler 水平自动伸缩配置
type HorizontalPodAutoscaler struct {
	Id                       int64      `gorm:"column:id;primary_key" json:"id"`
	ClusterId                int64      `gorm:"column:cluster_id;notnull;index;comment:'集群ID'" json:"cluster_id"`
	Namespace                string     `gorm:"column:namespace;notnull;index;comment:'空间'" json:"namespace"`
	Name                     string     `gorm:"column:name;notnull;index;comment:'名称'" json:"name"`
	AppName                  string     `gorm:"column:app_name;notnull;index;comment:'服务名称'" json:"app_name"`
	ApiVersion               string     `gorm:"column:api_version;notnull;comment:'API VERSION'" json:"api_version"`
	MinReplicas              int        `gorm:"column:min_replicas;null;default:1;comment:'最小实例数'" json:"min_replicas"`
	MaxReplicas              int        `gorm:"column:max_replicas;null;default:3;comment:'最大实例数'" json:"max_replicas"`
	ResourceName             string     `gorm:"column:resource_name;null;default:'cpu';comment:'集群ID'" json:"resource_name"`
	TargetAverageUtilization int        `gorm:"column:target_average_utilization;null;default:50;comment:'cpu平均负载'" json:"target_average_utilization"`
	Kind                     Kind       `gorm:"column:kind;notnull;size:32;comment:'资源类型'" json:"kind"`
	Remark                   string     `gorm:"column:remark;size:255;null;comment:'备注'" json:"remark"`
	Detail                   string     `gorm:"column:detail;null;type:text;comment:'源文'" json:"detail"`
	CreatedAt                time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt                time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt                *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName set table
func (*HorizontalPodAutoscaler) TableName() string {
	return "horizontal_pod_autoscaler"
}
