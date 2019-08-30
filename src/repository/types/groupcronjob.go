/**
 * Created by GoLand.
 * Email: xzghua@gmail.com
 * Date: 2019-07-12
 * Time: 16:44
 */
package types

type GroupsCronjobs struct {
	CronjobsID int64 `gorm:"column:cronjobs_id"`
	GroupsID   int64 `gorm:"column:groups_id"`
	ID         int64 `gorm:"column:id;primary_key"`
}

// TableName sets the insert table name for this struct type
func (g *GroupsCronjobs) TableName() string {
	return "groups_cronjobss"
}
