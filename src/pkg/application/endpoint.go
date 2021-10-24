/**
 * @Time: 2021/10/23 16:47
 * @Author: solacowa@gmail.com
 * @File: endpoint
 * @Software: GoLand
 */

package application

import "github.com/go-kit/kit/endpoint"

type (
	appResult struct {
	}
)

type Endpoints struct {
	ListEndpoint endpoint.Endpoint
}
