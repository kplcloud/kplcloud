/**
 * @Time : 2019-07-29 15:33 
 * @Author : soupzhb@gmail.com
 * @File : groupsmemberss.go
 * @Software: GoLand
 */

package types

type GroupsMemberss struct {
	GroupsID  int64 `gorm:"column:groups_id"`
	ID        int64 `gorm:"column:id;primary_key"`
	MembersID int64 `gorm:"column:members_id"`
}

// TableName sets the insert table name for this struct type
func (g *GroupsMemberss) TableName() string {
	return "groups_memberss"
}
