/**
 * @Time : 2019-07-29 10:35
 * @Author : solacowa@gmail.com
 * @File : dockerfile
 * @Software: GoLand
 */

package types

import "gopkg.in/guregu/null.v3"

type Dockerfile struct {
	AuthorID   int64       `gorm:"column:author_id" json:"author_id"`
	CreatedAt  null.Time   `gorm:"column:created_at" json:"created_at"`
	Desc       null.String `gorm:"column:desc" json:"desc"`
	Detail     string      `gorm:"column:detail;type:text" json:"detail"`
	Dockerfile string      `gorm:"column:dockerfile;type:text" json:"dockerfile"`
	Download   null.Int    `gorm:"column:download" json:"download"`
	FullPath   string      `gorm:"column:full_path" json:"full_path"`
	ID         int64       `gorm:"column:id;primary_key" json:"id"`
	Language   string      `gorm:"column:language" json:"language"`
	Name       string      `gorm:"column:name" json:"name"`
	Score      null.Int    `gorm:"column:score" json:"score"`
	Sha256     null.String `gorm:"column:sha256" json:"sha_256"`
	Status     null.Int    `gorm:"column:status" json:"status"`
	UpdatedAt  null.Time   `gorm:"column:updated_at" json:"updated_at"`
	UploaderID int64       `gorm:"column:uploader_id" json:"uploader_id"`
	Version    string      `gorm:"column:version" json:"version"`

	Member Member `gorm:"ForeignKey:id;AssociationForeignKey:author_id" json:"member"`
}

// TableName sets the insert table name for this struct type
func (d *Dockerfile) TableName() string {
	return "dockerfiles"
}
