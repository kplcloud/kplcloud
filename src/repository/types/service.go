/**
 * @Time : 2021/9/10 2:33 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package types

import "time"

type Port struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type Service struct {
	Id          int64      `gorm:"column:id;primary_key" json:"id"`
	ClusterId   int64      `gorm:"column:cluster_id;index;notnull;comment:'集群ID'" json:"cluster_id"`
	Namespace   string     `gorm:"column:namespace;index;size:32;notnull;comment:'空间'" json:"namespace"`
	Name        string     `gorm:"column:name;index;size:32;notnull;comment:'名称'" json:"name"`
	Ports       string     `gorm:"column:ports;size:1024;null;comment:'端口信息json'" json:"ports"`
	Selector    string     `gorm:"column:selector;size:2048;null;comment:'Selector选择器'" json:"selector"`
	ServiceType string     `gorm:"column:service_type;size:32;null;comment:'服务类型'" json:"service_type"`
	Detail      string     `gorm:"column:detail;size:2000;null;comment:'详情'" json:"detail"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName set table
func (*Service) TableName() string {
	return "services"
}
