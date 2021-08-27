/**
 * @Time: 2021/8/24 21:38
 * @Author: solacowa@gmail.com
 * @File: volume
 * @Software: GoLand
 */

package types

// DeploymentVolume 卷
type DeploymentVolume struct {
	Id           int64  `json:"id"`
	DeploymentId int64  `json:"deployment_id"`
	Name         string `json:"name"`
}
