/**
 * @Time: 2021/8/24 21:40
 * @Author: solacowa@gmail.com
 * @File: deployment_container
 * @Software: GoLand
 */

package types

type DeploymentContainer struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Data  string `json:"data"`
}
