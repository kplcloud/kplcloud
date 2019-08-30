package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type TemplateRepository interface {
	Create(name, kind, detail string) error
	Update(id int, name, kind, detail string) error
	FindById(id int) (v *types.Template, err error)
	FindByKindType(kind TplKind) (v *types.Template, err error)
	DeleteById(id int) error
	Count(name string) (count int, err error)
	FindOffsetLimit(name string, offset, limit int) (res []*types.Template, err error)
	GetTemplateByKind(kind string) (*types.Template, error)
}

type template struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &template{db: db}
}

type TplKind string

const (
	FilebeatConfigMap         TplKind = "FilebeatConfigMap"
	FilebeatContainer         TplKind = "FilebeatContainer"
	ConfigMapKind             TplKind = "ConfigMap"
	StorageClassKind          TplKind = "StorageClass"
	GolangJenkinsKind         TplKind = "GolangJenkins"
	JenkinsCommand            TplKind = "JenkinsCommand"
	PythonKind                TplKind = "Python"
	RegcredSecretKind         TplKind = "RegcredSecret"
	JenkinsNodeCommandKind    TplKind = "JenkinsNodeCommand"
	InitContainersKind        TplKind = "InitContainers"
	IstioProxyKind            TplKind = "IstioProxy"
	ServiceKind               TplKind = "Service"
	DeploymentKind            TplKind = "Deployment"
	DeploymentNoConfigMap     TplKind = "DeploymentNoConfigMap"
	EndpointsKind             TplKind = "Endpoints"
	ServiceEntryKind          TplKind = "ServiceEntry"
	VirtualServiceKind        TplKind = "VirtualService"
	GatewayKind               TplKind = "Gateway"
	GatewayGrpcKind           TplKind = "GatewayGrpc"
	CronJobKind               TplKind = "CronJob"
	PersistentVolumeClaimKind TplKind = "PersistentVolumeClaim"
	IngressKind               TplKind = "Ingress"
	RedisClusterStatefulSet   TplKind = "RedisClusterStatefulSet"
	RedisClusterConfigMap     TplKind = "RedisClusterConfigMap"
	ZookeeperStatefulSet      TplKind = "ZookeeperStatefulSet"
	ZookeeperConfigMap        TplKind = "ZookeeperConfigMap"
	RocketMqStatefulSet       TplKind = "RocketMqStatefulSet"
	RabbitMqStatefulSet       TplKind = "RabbitMqStatefulSet"
	MysqlStatefulSet          TplKind = "MysqlStatefulSet"
	JenkinsJavaCommand        TplKind = "JenkinsJavaCommand"
	NewVirtualService         TplKind = "NewVirtualService"
	InitVirtualService        TplKind = "InitVirtualService"
	InitIngress               TplKind = "InitIngress"
	EmailAlarm                TplKind = "EmailAlarm"
	EmailNotice               TplKind = "EmailNotice"
	EmailProclaim             TplKind = "EmailProclaim"
)

func (c TplKind) ToString() string {
	return string(c)
}

func (c *template) Create(name, kind, detail string) error {
	return c.db.Create(&types.Template{
		Name:   name,
		Kind:   kind,
		Detail: detail,
	}).Error
}

func (c *template) Update(id int, name, kind, detail string) error {
	var temp types.Template
	return c.db.Model(&temp).Where("id=?", id).Update(&types.Template{Name: name, Kind: kind, Detail: detail}).Error
}

func (c *template) FindByKindType(kind TplKind) (v *types.Template, err error) {
	var temp types.Template
	if err = c.db.Where("kind = ?", kind.ToString()).First(&temp).Error; err != nil {
		return
	}
	return &temp, nil
}

func (c *template) FindById(id int) (v *types.Template, err error) {
	var temp types.Template
	if err = c.db.First(&temp, id).Error; err != nil {
		return
	}
	return &temp, nil
}

func (c *template) DeleteById(id int) error {
	return c.db.Where("id=?", id).Delete(types.Template{}).Error
}

func (c *template) Count(name string) (count int, err error) {
	var temp types.Template
	query := c.db.Model(&temp)
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err = query.Count(&count).Error
	return
}

func (c *template) FindOffsetLimit(name string, offset, limit int) (res []*types.Template, err error) {
	var list []*types.Template
	query := c.db
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err = query.Offset(offset).Limit(limit).Find(&list).Error
	return list, err
}

func (c *template) GetTemplateByKind(kind string) (*types.Template, error) {
	var tem types.Template
	err := c.db.Where("kind = ?", kind).First(&tem).Error
	return &tem, err
}
