/**
 * @Time : 2021/9/6 5:15 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package pvc

import "context"

type Service interface {
	FindByName(ctx context.Context, clusterId int64, ns, name string) (err error)
}
