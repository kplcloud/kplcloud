/**
 * @Time: 2021/8/24 22:09
 * @Author: solacowa@gmail.com
 * @File: cluster_role
 * @Software: GoLand
 */

package types

import "time"

type Kind string

const (
	KindClusterRole Kind = "ClusterRole"
)

func (s Kind) String() string {
	return string(s)
}

type ClusterRole struct {
	Id        int64      `json:"id"`
	ClusterId int64      `json:"cluster_id"`
	Name      string     `json:"name"`
	Data      string     `json:"data"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	Rules []PolicyRule
}

// TableName sets the insert table name for this struct type
func (c *ClusterRole) TableName() string {
	return "cluster_role"
}
