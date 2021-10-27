/**
 * @Time : 2021/10/27 5:17 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package audits

import "context"

// Service 审核日志模块
// 所有非Get请求的日志都会在这展示
// 可根据类型、路由、空间、集群、服务等过滤
// 只有查询功能不提供其他操作
type Service interface {
	// List 获取审计日志列表
	// TODO: 应该会有很多查询条件，具体的后面再定
	List(ctx context.Context, query string, page, pageSize int) (res []auditResult, total int, err error)
}
