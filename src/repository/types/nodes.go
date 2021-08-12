/**
 * @Time : 8/11/21 11:17 AM
 * @Author : solacowa@gmail.com
 * @File : nodes
 * @Software: GoLand
 */

package types

import "time"

type Nodes struct {
	Id               int64      `gorm:"column:id;rimary_key" json:"id"`
	ClusterId        int64      `gorm:"column:cluster_id;notnull;index;comment:'集群ID'" json:"cluster_id"`
	Name             string     `gorm:"column:name;notnull;index;comment:'节点名称'" json:"name"`
	Memory           int64      `gorm:"column:memory;null;comment:'节点内存单位字节'" json:"memory"`
	Cpu              int64      `gorm:"column:cpu;null;default:1;comment:'节点CPU核数'" json:"cpu"`
	EphemeralStorage int64      `gorm:"column:ephemeral_storage;null;comment:'临时存储单位字节'" json:"ephemeral_storage"`
	InternalIp       string     `gorm:"column:internal_ip;null;comment:'节点内部IP'" json:"internal_ip"`
	ExternalIp       string     `gorm:"column:external_ip;null;comment:'节点外部IP'" json:"external_ip"`
	KubeletVersion   string     `gorm:"column:kubelet_version;null;comment:'kubelet版本'" json:"kubelet_version"`
	KubeProxyVersion string     `gorm:"column:kube_proxy_version;null;comment:'kubeproxy版本'" json:"kube_proxy_version"`
	ContainerVersion string     `gorm:"column:container_version;null;comment:'容器版本'" json:"container_version"`
	OsImage          string     `gorm:"column:os_image;null;comment:'系统镜像'" json:"os_image"`
	Status           string     `gorm:"column:status;null;comment:'状态'" json:"status"`
	Scheduled        bool       `gorm:"column:scheduled;null;default:false;comment:'是否调度'" json:"scheduled"`
	Remark           string     `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	Labels []Label `gorm:"many2many:node_label;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:label_id;jointable_foreignkey:node_id;"`
}

// TableName set table
func (*Nodes) TableName() string {
	return "nodes"
}
