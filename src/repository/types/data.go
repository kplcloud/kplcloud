/**
 * @Time : 8/19/21 1:55 PM
 * @Author : solacowa@gmail.com
 * @File : data
 * @Software: GoLand
 */

package types

import "time"

type DataStyle int

const (
	DataStyleConfigMap DataStyle = 1
	DataStyleSecret    DataStyle = 2
)

type Data struct {
	Id        int64      `gorm:"column:id;primary_key" json:"id"`
	TargetId  int64      `gorm:"column:target_id;notnull;index;comment:'上级ID'" json:"target_id"`
	Style     DataStyle  `gorm:"column:style;notnull;index;default:1;comment:'类型:1:ConfigMap,2:Secret'" json:"style"`
	Key       string     `gorm:"column:key;notnull;comment:'Key'" json:"key"`
	Value     string     `gorm:"column:value;null;type:text;comment:'Value'" json:"value"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (c *Data) TableName() string {
	return "data"
}
