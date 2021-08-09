/**
 * @Time : 8/9/21 6:20 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package cluster

import "context"

type Service interface {
	Add(ctx context.Context, data string) (err error)
}

type service struct {
}
