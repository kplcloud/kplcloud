/**
 * @Time : 2019-07-25 09:56
 * @Author : solacowa@gmail.com
 * @File : client
 * @Software: GoLand
 */

package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

type Client struct {
	*gitlab.Client
}

func NewGitlabClient(client *gitlab.Client) *Client {
	return &Client{client}
}
