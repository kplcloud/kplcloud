/**
 * @Time: 2020/3/27 17:34
 * @Author: solacowa@gmail.com
 * @File: responsestatus
 * @Software: GoLand
 */

package encode

import (
	"github.com/pkg/errors"
)

type ResStatus string

var ResponseMessage = map[ResStatus]int{
	Invalid:        400,
	InvalidParams:  400,
	ErrParamsPhone: 401,
	ErrBadRoute:    401,
	ErrSystem:      500,
	ErrNotfound:    404,
	ErrLimiter:     429,

	ErrAccountNotFound:    404,
	ErrAccountLogin:       1002,
	ErrAccountLoginIsNull: 1003,
	ErrAccountNotLogin:    501,
	ErrAccountASD:         1004,
	ErrAccountLocked:      1005,
	ErrAuthNotLogin:       501,

	// 系统API
	ErrSysRoleNotfound:     2001,
	ErrSysRoleSave:         2002,
	ErrSysRoleUserNotfound: 2003,
	ErrSysRoleUser:         2004,
	ErrSysRoleUserLen:      2005,
	ErrSysRoleUserDelete:   2006,
	ErrSysUserNotfound:     2007,

	ErrInstallDbConnect:  901,
	ErrInstallDbDrive:    902,
	ErrInstallUploadPath: 903,
	ErrInstallWebPath:    904,
	ErrInstallUpload:     905,

	ErrClusterAdd:      920,
	ErrClusterConnect:  921,
	ErrClusterNotfound: 922,
	ErrClusterParams:   923,

	ErrNamespaceNotfound: 940,
	ErrNameNotfound:      941,

	ErrDeploymentSyncList:    960,
	ErrDeploymentGetNotfound: 961,

	ErrConfigMapSyncList: 980,

	ErrStorageClassSync:     1100,
	ErrStorageClassNotfound: 1101,
	ErrStorageClassSyncPv:   1102,
}

const (
	// 公共错误信息
	Invalid        ResStatus = "invalid"
	InvalidParams  ResStatus = "请求参数错误"
	ErrNotfound    ResStatus = "不存在"
	ErrBadRoute    ResStatus = "请求路由错误"
	ErrParamsPhone ResStatus = "手机格式不正确"
	ErrLimiter     ResStatus = "太快了,等我一会儿..."

	ErrInstallDbConnect  ResStatus = "数据库连接失败"
	ErrInstallDbDrive    ResStatus = "暂不支持其他数据库"
	ErrInstallUploadPath ResStatus = "文件目录未配置"
	ErrInstallWebPath    ResStatus = "Web目录未配置"
	ErrInstallUpload     ResStatus = "上传文件失败"

	ErrClusterAdd      ResStatus = "集群添加错误"
	ErrClusterConnect  ResStatus = "集群链接错误"
	ErrClusterNotfound ResStatus = "集群不存在"
	ErrClusterParams   ResStatus = "集群参数错误"

	ErrNamespaceNotfound ResStatus = "空间不存在"
	ErrNameNotfound      ResStatus = "名称不存在"

	ErrDeploymentSyncList    ResStatus = "同步失败"
	ErrDeploymentGetNotfound ResStatus = "项目不存在"

	ErrConfigMapSyncList ResStatus = "ConfigMap同步失败"

	ErrStorageClassSync     ResStatus = "StorageClass同步失败"
	ErrStorageClassNotfound ResStatus = "StorageClass不存在"
	ErrStorageClassSyncPv   ResStatus = "StorageClass Pv同步失败"

	// 中间件错误信息
	ErrSystem             ResStatus = "系统错误"
	ErrAccountNotLogin    ResStatus = "用户没登录"
	ErrAuthNotLogin       ResStatus = "请先登录"
	ErrAccountLoginIsNull ResStatus = "用户名和密码不能为空"
	ErrAccountLogin       ResStatus = "用户名或密码错误"
	ErrAccountNotFound    ResStatus = "账号不存在"
	ErrAccountASD         ResStatus = "权限验证失败"
	ErrAccountLocked      ResStatus = "用户已被锁定"

	// 系统API
	ErrSysRoleNotfound     ResStatus = "角色不存在"
	ErrSysRoleSave         ResStatus = "角色保证错误"
	ErrSysRoleUserNotfound ResStatus = "角色用户不存在"
	ErrSysRoleUser         ResStatus = "用户角色配置失败"
	ErrSysRoleUserLen      ResStatus = "请选择用户"
	ErrSysRoleUserDelete   ResStatus = "角色删除失败"
	ErrSysUserNotfound     ResStatus = "用户不存在"
)

func (c ResStatus) String() string {
	return string(c)
}

func (c ResStatus) Error() error {
	return errors.New(string(c))
}

func (c ResStatus) Wrap(err error) error {
	return errors.Wrap(err, string(c))
}
