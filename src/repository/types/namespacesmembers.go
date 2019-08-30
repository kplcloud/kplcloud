/**
 * Created by GoLand.
 * User: zghua
 * Email: xzghua@gmail.com
 * Date: 2019-07-01
 * Time: 15:36
 */
package types

type NamespacesMembers struct {
	ID           int64 `gorm:"column:id;primary_key"`
	MembersID    int64 `gorm:"column:members_id"`
	NamespacesID int64 `gorm:"column:namespaces_id"`
}

// TableName sets the insert table name for this struct type
func (n *NamespacesMembers) TableName() string {
	return "namespaces_memberss"
}