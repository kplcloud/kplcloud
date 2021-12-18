/**
 * @Time : 2019-06-26 15:14
 * @Author : solacowa@gmail.com
 * @File : persistentvolumeclaim
 * @Software: GoLand
 */

package types

import (
	"time"
)

// PersistentVolume 存储卷信息
type PersistentVolume struct {
	Id                            int64      `gorm:"column:id;primary_key" json:"id"`
	ClusterId                     int64      `gorm:"column:cluster_id;index;notnull;comment:'集群ID'" json:"cluster_id"`
	Name                          string     `gorm:"column:name;index;size:64;notnull;comment:'名称'" json:"name"`
	Namespace                     string     `gorm:"column:namespace;size:64;notnull;index;comment:'空间'" json:"namespace"`
	AccessModes                   string     `gorm:"column:access_modes;null;comment:'访问模式'" json:"access_modes"`
	Remark                        string     `gorm:"column:remark;null;size:1000;comment:'备注'" json:"desc"`
	PersistentVolumeSource        string     `gorm:"column:persistent_volume_source;notnull;comment:'pv源,NFS,cephFS,GFS'" json:"persistent_volume_source"`
	PvcId                         int64      `gorm:"column:pvc_id;index;notnull;comment:'PvcId'" json:"pvc_id"`
	PersistentVolumeReclaimPolicy string     `gorm:"column:persistent_volume_reclaim_policy;null;comment:'生命周期结束维护的策略'" json:"persistent_volume_reclaim_policy"`
	PersistentVolumeMode          string     `gorm:"column:persistent_volume_mode;null;comment:'描述卷'" json:"persistent_volume_mode"`
	StorageClassId                int64      `gorm:"column:storage_class_id;index;notnull;comment:'存储类ID'" json:"storage_class_id"`
	Storage                       string     `gorm:"column:storage;notnull;comment:'容量'" json:"storage"`
	Status                        string     `gorm:"column:status;size:14;null;comment:'状态'" json:"status"`
	Detail                        string     `gorm:"column:detail;size:5000;null;comment:'详情'" json:"detail"`
	CreatedAt                     time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt                     time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt                     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (p *PersistentVolume) TableName() string {
	return "persistent_volume"
}

// PersistentVolumeClaim 存储卷声明信息
type PersistentVolumeClaim struct {
	Id             int64      `gorm:"column:id;primary_key" json:"id"`
	ClusterId      int64      `gorm:"column:cluster_id;index;notnull;comment:'集群ID'" json:"cluster_id"`
	Name           string     `gorm:"column:name;index;notnull;comment:'名称'" json:"name"`
	Namespace      string     `gorm:"column:namespace;size:64;notnull;index;comment:'空间'" json:"namespace"`
	AccessModes    string     `gorm:"column:access_modes;null;comment:'访问模式'" json:"access_modes"`
	Remark         string     `gorm:"column:remark;null;size:1000;comment:'备注'" json:"desc"`
	Labels         string     `gorm:"column:labels;null;size:1000;comment:'label信息'" json:"labels"`
	RequestStorage string     `gorm:"column:request_storage;notnull;size:255;comment:'请求容量'" json:"request_storage"`
	LimitStorage   string     `gorm:"column:limit_storage;notnull;size:255;comment:'最大容量'" json:"limit_storage"`
	StorageClassId int64      `gorm:"column:storage_class_id;index;notnull;comment:'存储类ID'" json:"storage_class_id"`
	Status         string     `gorm:"column:status;null;comment:'状态'" json:"status"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	StorageClass StorageClass `gorm:"foreignkey:StorageClassId;references:Id"`
	Cluster      Cluster      `gorm:"foreignkey:ClusterId;references:Id"`
}

// TableName sets the insert table name for this struct type
func (p *PersistentVolumeClaim) TableName() string {
	return "persistent_volume_claim"
}
