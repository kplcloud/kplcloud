/**
 * @Time : 8/16/21 6:17 PM
 * @Author : solacowa@gmail.com
 * @File : reject
 * @Software: GoLand
 */

package types

// 项目拒绝表
type ProjectReject struct {
	Id        int64  `json:"id"`
	ProjectId int64  `json:"project_id"`
	Desc      string `json:"desc"`
}
