/**
 * @Time : 2019/7/1 5:48 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : types
 * @Software: GoLand
 */

package event

import (
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

// 接收数据
type WebhooksRequest struct {
	AppName   string                `json:"app_name"`
	Namespace string                `json:"namespace"`
	Event     repository.EventsKind `json:"event"`
	EventDesc string                `json:"event_desc"`
	MemberId  int64                 `json:"member_id"`
	Member    *types.Member         `json:"member"`
	Title     string                `json:"title"`
	Message   string                `json:"message"`
	Project   *types.Project        `json:"project"`
}

// 发送数据
type Params struct {
	AppName       string `json:"app_name"`
	Namespace     string `json:"namespace"`
	Event         string `json:"event"`
	Operator      string `json:"operator"`
	OperatorEmail string `json:"operator_email"`
	Title         string `json:"title"`
	Message       string `json:"message"`
	Project       struct {
		Name        string        `json:"name"`
		NameEn      string        `json:"name_en"`
		Namespace   string        `json:"namespace"`
		ProjectId   int64         `json:"project_id"`
		Member      string        `json:"member"`
		Email       string        `json:"email"`
		Description string        `json:"description"`
		Groups      []groupStruct `json:"groups"`
	} `json:"project"`
}

type groupStruct struct {
	GroupName   string `json:"group_name"`
	GroupNameEn string `json:"group_name_en"`
}
