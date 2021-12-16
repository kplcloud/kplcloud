/**
 * @Time : 2021/12/16 4:37 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package group

import "context"

// Service 组管理模块
// 每个空间下都可以创建很多组
// 每个组可以添加app,定时任务,有状态服务
// 组里可以添加成员及权限 读/操作
// 全部都是多对多的关系
// 两种方案
//   1. 组的类型，只读或可操作 组里人跟随组权限
//   2. 组没有类型，组里的人设置读写权限
// 我决定采用1 方案
// 创建空间时默认给创建一个写组，添加新成员时必须给分配一个组
type Service interface {
	// Create 创建组
	// userId 组管理员
	Create(ctx context.Context, clusterId, sysUserId int64, namespace, groupName, groupAlias string, userId int64, onlyRead bool) (err error)
	// Update 更新组信息
	Update(ctx context.Context, clusterId, sysUserId int64, namespace, groupName, groupAlias string, userId int64, onlyRead bool) (err error)
	// Delete 删除组
	// 删除组需要把他们关联关系全部删除
	Delete(ctx context.Context, clusterId, sysUserId int64, namespace, groupName string) (err error)
	// List 组列表
	// groupName 可以为空
	List(ctx context.Context, clusterId, sysUserId int64, namespace, groupName string, page, pageSize int) (res []result, total int, err error)
	// AddUser 给组添加成员
	AddUser(ctx context.Context, clusterId, sysUserId int64, namespace, groupName string, userIds []int64) (err error)
	// AddApp 给组添加app
	AddApp(ctx context.Context, clusterId, sysUserId int64, namespace, groupName string, names []string) (err error)
	// AddCronJob 给组添加定时任务
	AddCronJob(ctx context.Context, clusterId, sysUserId int64, namespace, groupName string, names []string) (err error)
	// AddStatefulSet 给组添加有状态应用
	AddStatefulSet(ctx context.Context, clusterId, sysUserId int64, namespace, groupName string, names []string) (err error)
	// DeleteApp 删除组里的应用
	DeleteApp(ctx context.Context, clusterId, sysUserId int64, namespace, groupName string, names []string) (err error)
	// DeleteCronJob 删除组里的定时任务
	DeleteCronJob(ctx context.Context, clusterId, sysUserId int64, namespace, groupName string, names []string) (err error)
	// DeleteStatefulSet 删除组里的有状态应用
	DeleteStatefulSet(ctx context.Context, clusterId, sysUserId int64, namespace, groupName string, names []string) (err error)
}
