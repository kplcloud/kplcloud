/**
 * @Time : 2021/12/16 4:49 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package sysgroup

import "time"

type (
	result struct {
		Alias     string    `json:"alias"`
		Name      string    `json:"name"`
		Namespace string    `json:"namespace"`
		Remark    string    `json:"remark"`
		User      string    `json:"user"`
		OnlyRead  bool      `json:"onlyRead"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)
