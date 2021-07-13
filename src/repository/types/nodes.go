/**
 * @Time : 6/15/21 5:46 PM
 * @Author : solacowa@gmail.com
 * @File : nodes
 * @Software: GoLand
 */

package types

import "time"

type Node struct {
	Id               int64      `gorm:"column:id;primary_key" json:"id"`
	Name             string     `gorm:"column:name" json:"name"`
	InternalIP       string     `json:"internalIP"`
	Status           string     `json:"status"` // 状态
	CpuCapacity      int        `json:"cpuCapacity"`
	MemoryCapacity   float64    `json:"memoryCapacity"`   // 内存
	EphemeralStorage string     `json:"ephemeralStorage"` // 存储空容量
	PodCapacity      int        `json:"podCapacity"`      // POD 最容量
	AllocatedPods    int        `json:"allocatedPods"`    // 已分配
	ExternalIP       string     `json:"externalIP"`
	Hostname         string     `json:"hostname"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (e *Node) TableName() string {
	return "nodes"
}
