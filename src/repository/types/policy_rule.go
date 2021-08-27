/**
 * @Time: 2021/8/24 22:14
 * @Author: solacowa@gmail.com
 * @File: policy_rule
 * @Software: GoLand
 */

package types

import "time"

type PolicyRule struct {
	Id              int64      `json:"id"`
	Kind            string     `json:"kind"`
	TargetId        int64      `json:"target_id"`
	Verbs           string     `json:"verbs"`
	APIGroups       string     `json:"api_groups"`
	Resources       string     `json:"resources"`
	ResourceNames   string     `json:"resource_names"`
	NonResourceURLs string     `json:"non_resource_ur_ls"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (c *PolicyRule) TableName() string {
	return "policy_rule"
}
