/**
 * @Time : 2019/7/5 11:02 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : service
 * @Software: GoLand
 */
package configmap

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/jenkins"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/pkg/cronjob"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/configmapyaml"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	"gopkg.in/guregu/null.v3"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var (
	ErrConfigMapGetDB      = errors.New("configMap 数据获取失败")
	ErrConfigMapExist      = errors.New("configMap 已存在")
	ErrConfigMapNotExist   = errors.New("configMap 不存在")
	ErrConfigMapCreate     = errors.New("configMap 创建失败")
	ErrConfigMapUpdate     = errors.New("configMap 修改失败")
	ErrConfigMapDelete     = errors.New("configMap 删除失败")
	ErrConfigMapK8sGet     = errors.New("configMap 获取错误,可能不存在")
	ErrConfigMapDataGet    = errors.New("configMap Data 获取错误,可能不存在")
	ErrConfigMapK8sUpdate  = errors.New("configMap 更新错误")
	ErrConfigMapDataCount  = errors.New("configDataCount 不存在")
	ErrConfigMapDataList   = errors.New("configDataList 不存在")
	ErrSyncConfigMapYaml   = errors.New("syncConfigMapYaml 同步远程数据失败")
	ErrCreateConfigData    = errors.New("configData 创建失败")
	ErrUpdateConfigMapYaml = errors.New("updateConfigMapYaml 更新到远程失败")
	ErrUpdateConfigData    = errors.New("configData 更新失败")
	ErrDeleteConfigData    = errors.New("configData 删除失败")
	ErrConfigMapCreateYaml = errors.New("configMapYaml 创建失败")
	ErrConfigEnvFailed     = errors.New("configEnv 获取失败")
	ErrCreateConfEnvFailed = errors.New("coonfigEnv 创建失败")
	ErrUpdateConfEnvFailed = errors.New("coonfigEnv 更新失败")
	ErrExchangeCronJobTemp = errors.New("转换模板 失败")
)

type Service interface {
	// 获取configmap 详情
	GetOne(ctx context.Context, ns, name string) (res map[string]interface{}, err error)

	// 同步单个configmap
	GetOnePull(ctx context.Context, ns, name string) (res interface{}, err error)

	// configmap 列表
	List(ctx context.Context, req listRequest) (res map[string]interface{}, err error)

	// 创建configmap
	Post(ctx context.Context, req postRequest) error

	// 更新configmap
	Update(ctx context.Context, req postRequest) error

	// 删除configmap
	Delete(ctx context.Context, ns, name string) error

	// 同步空间的configmap
	Sync(ctx context.Context, ns string) error

	// 创建configmap
	CreateConfigMap(ctx context.Context, req createConfigMapRequest) error

	// 获取configmap
	GetConfigMap(ctx context.Context, ns, name string) (res interface{}, err error)

	// 获取configmap 的data 数据
	GetConfigMapData(ctx context.Context, ns, name string, page int, limit int) (res map[string]interface{}, err error)

	// 创建configmap 的data 数据
	CreateConfigMapData(ctx context.Context, req createConfigMapDataRequest) error

	// 更新configmap 的data 数据
	UpdateConfigMapData(ctx context.Context, req configMapDataRequest) error

	// 删除configmap 的data 数据
	DeleteConfigMapData(ctx context.Context, req configMapDataRequest) error

	// 获取configenv 的data 数据
	GetConfigEnv(ctx context.Context, name, ns string, page int, limit int) (res map[string]interface{}, err error)

	// 创建configenv 的data 数据
	CreateConfigEnv(ctx context.Context, req configEnvRequest) error

	// 更新configenv 的data 数据
	ConfigEnvUpdate(ctx context.Context, req configEnvRequest) error

	// 删除configenv 的data 数据
	ConfigEnvDel(ctx context.Context, req configEnvRequest) error
}

type service struct {
	logger     log.Logger
	config     *config.Config
	jenkins    jenkins.Jenkins
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func NewService(logger log.Logger, config *config.Config,
	jenkins jenkins.Jenkins,
	k8sClient kubernetes.K8sClient,
	repository repository.Repository) Service {
	return &service{logger, config,
		jenkins,
		k8sClient,
		repository}
}

/**
 * @Title 获取ConfigMap信息
 */
func (c *service) GetOne(ctx context.Context, ns, name string) (res map[string]interface{}, err error) {
	confData, err := c.repository.ConfigData().Find(ns, name)
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "GetOne", "err", err.Error())
		return nil, ErrConfigMapGetDB
	}
	var data []map[string]interface{}
	var confMap interface{}
	for _, v := range confData {
		if confMap == nil {
			confMap = v.ConfigMap
		}
		data = append(data, map[string]interface{}{
			"id":    v.ID,
			"key":   v.Key,
			"value": v.Value,
		})
	}
	res = map[string]interface{}{
		"confMap":  confMap,
		"confData": data,
	}
	return
}

/**
 * @Title 删除data
 */
func (c *service) DeleteData(ctx context.Context, ns, name string, id int64) (err error) {
	configMapData, err := c.repository.ConfigData().FindById(id)
	if err != nil {
		_ = level.Error(c.logger).Log("confDataRepository", "FindById", "err", err.Error())
		return ErrConfigMapDataGet
	}
	defer func() {
		if err == nil {
			if e := c.repository.ConfigData().Delete(id); e != nil {
				_ = level.Warn(c.logger).Log("confDataRepository", "Delete", "err", e.Error())
			}
		}
	}()

	configMap, err := c.k8sClient.Do().CoreV1().ConfigMaps(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMaps", "Get", "err", err.Error())
		return ErrConfigMapK8sGet
	}

	data := configMap.Data
	delete(data, configMapData.Key)
	configMap.Data = data

	configMap, err = c.k8sClient.Do().CoreV1().ConfigMaps(ns).Update(configMap)
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMaps", "Update", "err", err.Error())
		return ErrConfigMapK8sUpdate
	}

	return nil
}

/**
 * @Title 远程获取ConfigMap信息
 */
func (c *service) GetOnePull(ctx context.Context, ns, name string) (res interface{}, err error) {
	cf, err := c.k8sClient.Do().CoreV1().ConfigMaps(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "Get", "err", err.Error())
		return
	}
	cf.APIVersion = "v1"
	cf.Kind = "ConfigMap"
	return cf, nil
}

/**
 * @Title 获取ConfigMap列表
 */
func (c *service) List(ctx context.Context, req listRequest) (res map[string]interface{}, err error) {
	count, err := c.repository.ConfigMap().Count(req.Namespace, req.Name)
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "List Count", "err", err.Error())
		return nil, ErrConfigMapGetDB
	}

	p := paginator.NewPaginator(req.Page, req.Limit, count)

	list, err := c.repository.ConfigMap().FindOffsetLimit(req.Namespace, req.Name, p.Offset(), req.Limit)

	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "List", "err", err.Error())
		return nil, ErrConfigMapGetDB
	}
	res = map[string]interface{}{
		"list": list,
		"page": p.Result(),
	}
	return
}

/**
 * @Title 创建ConfigMap
 */
func (c *service) Post(ctx context.Context, req postRequest) error {
	if _, state := c.repository.ConfigMap().Find(req.Namespace, req.Name); state == false {
		_ = level.Error(c.logger).Log("ConfigMap", "Post", "Error", "exist", state, "state")
		return ErrConfigMapExist
	}

	//add database
	confMap, err := c.repository.ConfigMap().Create(&types.ConfigMap{
		Namespace: req.Namespace,
		Name:      req.Name,
		Desc:      req.Desc,
		Type:      null.IntFrom(1),
	})
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "Post", "Create Error", err.Error())
		return ErrConfigMapCreate
	}

	dat := map[string]string{}
	for _, v := range req.Data {
		if err = c.repository.ConfigData().Create(&types.ConfigData{
			Key:         v.Key,
			Value:       v.Value,
			ConfigMapID: confMap.ID,
		}); err != nil {
			_ = level.Error(c.logger).Log("ConfigMap", "Post", "ConfData Create Error", err.Error())
		}
		if v.Key != "" {
			dat[v.Key] = v.Value
		}

	}

	conf := new(v1.ConfigMap)
	conf.Namespace = req.Namespace
	conf.Name = req.Name
	conf.Data = dat

	cf, err := c.k8sClient.Do().CoreV1().ConfigMaps(req.Namespace).Create(conf)
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "Post", "ConfigMaps Create Error", err.Error())
		return ErrConfigMapCreate
	}
	cf.Kind = repository.ConfigMap.String()

	//校验项目Deployment是否需要创建ConfigMap
	if project, err := c.repository.Project().FindByNsName(req.Namespace, req.Name); err == nil && project.ID > 0 {
		if err = c.updateDeployment(req.Namespace, req.Name); err != nil {
			_ = level.Error(c.logger).Log("ConfigMap", "Post", "UpdateDeployment Error", err.Error())
		}
	}

	//@todo 写入操作记录

	return nil
}

/**
 * @Title 更新ConfigMap
 */
func (c *service) Update(ctx context.Context, req postRequest) error {
	confMap, state := c.repository.ConfigMap().Find(req.Namespace, req.Name)
	if state != false {
		_ = level.Error(c.logger).Log("ConfigMap", "Update Error. Not Exist")
		return ErrConfigMapNotExist
	}
	if confMap.ID <= 0 {
		return ErrConfigMapNotExist
	}

	if err := c.repository.ConfigMap().Update(req.Namespace, req.Name, req.Desc); err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "Update", "err", err.Error())
		return ErrConfigMapUpdate
	}

	if err := c.repository.ConfigData().Delete(confMap.ID); err != nil {
		_ = level.Error(c.logger).Log("ConfData", "Delete", "err", err.Error())
		return ErrConfigMapUpdate
	}

	dat := map[string]string{}
	for _, v := range req.Data {
		if err := c.repository.ConfigData().Create(&types.ConfigData{
			Key:         v.Key,
			Value:       v.Value,
			ConfigMapID: confMap.ID,
		}); err != nil {
			_ = level.Error(c.logger).Log("ConfigMap", "Update", "ConfData Create Error", err.Error())
		}
		dat[v.Key] = v.Value
	}
	conf := new(v1.ConfigMap)
	conf.Namespace = req.Namespace
	conf.Name = req.Name
	conf.Data = dat

	configInfo, err := c.k8sClient.Do().CoreV1().ConfigMaps(req.Namespace).Update(conf)
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "Update", "ConfigMaps Update Error", err.Error())
		return ErrConfigMapCreate
	}
	configInfo.Kind = repository.ConfigMap.String()

	//校验项目Deployment是否需要创建ConfigMap
	if project, err := c.repository.Project().FindByNsName(req.Namespace, req.Name); err == nil && project.ID > 0 {
		if err = c.updateDeployment(req.Namespace, req.Name); err != nil {
			_ = level.Error(c.logger).Log("ConfigMap", "Post", "UpdateDeployment Error", err.Error())
		}
	}

	// @todo 写入操作记录

	return nil
}

/**
 * @Title 删除ConfigMap
 */
func (c *service) Delete(ctx context.Context, ns, name string) error {
	confMap, state := c.repository.ConfigMap().Find(ns, name)
	if state == true {
		_ = level.Error(c.logger).Log("ConfigMap", "Delete", "Find", "Error")
	} else {
		if err := c.repository.ConfigMap().Delete(confMap.ID); err != nil {
			_ = level.Error(c.logger).Log("ConfigMap", "Delete", "Error", err.Error())
		}
		if err := c.repository.ConfigData().Delete(confMap.ID); err != nil {
			_ = level.Error(c.logger).Log("ConfigData", "Delete", "Error", err.Error())
		}
	}

	//远程删除数据
	err := c.k8sClient.Do().CoreV1().ConfigMaps(ns).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "Delete", "Error", err.Error())
		return ErrConfigMapDelete
	}
	return nil
}

/**
 * @Title 同步远程数据
 */
func (c *service) Sync(ctx context.Context, ns string) error {
	confMapList, err := c.k8sClient.Do().CoreV1().ConfigMaps(ns).List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigMap", "Sync", "Error", err.Error())
		return err
	}
	for _, v := range confMapList.Items {
		if err = c.updateOrCreateDB(v.Namespace, v.Name, v.Data); err != nil {
			_ = level.Error(c.logger).Log("ConfMap", "Sync", "UpdateOrCreateDB Error", err.Error())
		}
	}
	return nil
}

// 更新 deployment
func (c *service) updateDeployment(ns, name string) error {
	deployment, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if deployment.Spec.Template.Spec.Containers[0].EnvFrom == nil || len(deployment.Spec.Template.Spec.Containers[0].EnvFrom) == 0 {
		var envFrom = new(v1.EnvFromSource)
		envFrom.ConfigMapRef = &v1.ConfigMapEnvSource{}
		envFrom.ConfigMapRef.LocalObjectReference.Name = name
		deployment.Spec.Template.Spec.Containers[0].EnvFrom = append(deployment.Spec.Template.Spec.Containers[0].EnvFrom, *envFrom)
		var volume = new(v1.Volume)
		volume.Name = name
		volume.ConfigMap = &v1.ConfigMapVolumeSource{}
		volume.ConfigMap.Name = name
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, *volume)
		_, err = c.k8sClient.Do().AppsV1().Deployments(ns).Update(deployment)
	}
	return err
}

// 更新或创建操作
func (c *service) updateOrCreateDB(ns, name string, data map[string]string) (err error) {
	confMap, state := c.repository.ConfigMap().Find(ns, name)
	if state == true {
		if confMap, err = c.repository.ConfigMap().Create(&types.ConfigMap{
			Namespace: ns,
			Name:      name,
			Desc:      name,
			Type:      null.IntFrom(1),
		}); err != nil {
			return err
		}
	}
	err = c.repository.ConfigData().Delete(confMap.ID)
	if err != nil {
		return err
	}

	for k, v := range data {
		if err := c.repository.ConfigData().Create(&types.ConfigData{
			Key:         k,
			Value:       v,
			ConfigMapID: confMap.ID,
		}); err != nil {
			_ = level.Error(c.logger).Log("ConfigMap", "Update", "ConfData Create Error", err.Error())
		}
	}

	return nil
}

//##配置字典和数据分开添加 【以下接口】
//创建config map
func (c *service) CreateConfigMap(ctx context.Context, req createConfigMapRequest) (err error) {
	//add database
	confMap, err := c.repository.ConfigMap().Create(&types.ConfigMap{
		Namespace: req.Namespace,
		Name:      req.Name,
		Desc:      req.Name,
		Type:      null.IntFrom(req.Type), //类型=2 定时任务 类型=1 项目
	})

	if err != nil {
		_ = level.Error(c.logger).Log("configmap", "Post", "CreateConfigMap", err.Error())
		return ErrConfigMapCreate
	}

	//create configMap yaml
	conf := new(v1.ConfigMap)
	conf.Namespace = confMap.Namespace
	conf.Name = confMap.Name

	_, err = c.k8sClient.Do().CoreV1().ConfigMaps(req.Namespace).Create(conf)
	if err != nil {
		_ = level.Error(c.logger).Log("CreateConfigMap", "Post", "ConfigMaps Create Error", err.Error())
		return ErrConfigMapCreateYaml
	}
	//cfg.Kind = repository.ConfigMap.String()

	if req.Type == 1 {
		//校验项目Deployment是否需要创建ConfigMap
		if project, err := c.repository.Project().FindByNsName(req.Namespace, req.Name); err == nil && project.ID > 0 {
			if err = c.updateDeployment(req.Namespace, req.Name); err != nil {
				_ = level.Error(c.logger).Log("CreateConfigMap", "Post", "UpdateDeployment Error", err.Error())
			}
		}
	}

	return
}

//获取config map
func (c *service) GetConfigMap(ctx context.Context, ns, name string) (res interface{}, err error) {

	res, notFound := c.repository.ConfigMap().Find(ns, name)
	if notFound == true {
		return nil, ErrConfigMapNotExist
	}
	return
}

//获取config map data list
func (c *service) GetConfigMapData(ctx context.Context, ns, name string, page int, limit int) (res map[string]interface{}, err error) {

	configMap, notFound := c.repository.ConfigMap().Find(ns, name)

	if notFound == true {
		return nil, ErrConfigMapNotExist
	}
	count, err := c.repository.ConfigData().Count(configMap.ID)
	if err != nil {
		_ = level.Error(c.logger).Log("GetConfigMapData", "Count", "err", err.Error())
		return nil, ErrConfigMapDataCount
	}
	p := paginator.NewPaginator(page, limit, count)

	list, err := c.repository.ConfigData().FindOffsetLimit(configMap.ID, p.Offset(), limit)

	if err != nil {
		_ = level.Error(c.logger).Log("GetConfigMapData", "FindOffsetLimit", "err", err.Error())
		return nil, ErrConfigMapDataList
	}
	res = map[string]interface{}{
		"list": list,
		"page": map[string]interface{}{
			"total":     count,
			"pageTotal": p.PageTotal(),
			"pageSize":  limit,
			"page":      p.Page(),
		},
	}
	return
}

//创建config map data
func (c *service) CreateConfigMapData(ctx context.Context, req createConfigMapDataRequest) (err error) {

	confMap, notFound := c.repository.ConfigMap().FindById(req.ConfigMapId)
	if notFound == true {
		return ErrConfigMapNotExist
	}

	//先同步远程configMapYaml 数据
	err = configmapyaml.SyncConfigMapYaml(confMap.Namespace, confMap.Name, c.logger, c.k8sClient, c.repository)

	if err != nil {
		return ErrSyncConfigMapYaml
	}

	//新增数据入库
	data := &types.ConfigData{
		Key:         req.Key,
		Value:       req.Value,
		ConfigMapID: req.ConfigMapId,
	}
	err = c.repository.ConfigData().Create(data)
	if err != nil {
		return ErrCreateConfigData
	}

	err = configmapyaml.UpdateConfigMapYaml(req.ConfigMapId, c.logger, c.k8sClient, c.repository)
	if err != nil {
		return ErrUpdateConfigMapYaml
	}

	return
}

//修改config map data
func (c *service) UpdateConfigMapData(ctx context.Context, req configMapDataRequest) (err error) {

	confMap, notFound := c.repository.ConfigMap().FindById(req.ConfigMapId)

	if notFound == true {
		return ErrConfigMapNotExist
	}

	//先同步远程configMapYaml 数据

	err = configmapyaml.SyncConfigMapYaml(confMap.Namespace, confMap.Name, c.logger, c.k8sClient, c.repository)

	if err != nil {
		return ErrSyncConfigMapYaml
	}

	//更新数据入库
	err = c.repository.ConfigData().Update(req.ConfigMapDataId, req.Value, req.Path)

	if err != nil {
		return ErrUpdateConfigData
	}

	//更新远程
	err = configmapyaml.UpdateConfigMapYaml(req.ConfigMapId, c.logger, c.k8sClient, c.repository)
	if err != nil {
		return ErrUpdateConfigMapYaml
	}

	if req.Path != "" {
		if deployment, err := c.k8sClient.Do().AppsV1().Deployments(confMap.Namespace).Get(confMap.Name, metav1.GetOptions{}); err == nil {
			for k, container := range deployment.Spec.Template.Spec.Containers {
				if container.Name == confMap.Name {
					for _, v := range container.VolumeMounts {
						if v.Name == confMap.Name {
							continue
						}
						reqPath := req.Path
						if !strings.Contains(req.Path, req.Key) {
							reqPath = req.Path + "/" + req.Key
						}
						deployment.Spec.Template.Spec.Containers[k].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[k].VolumeMounts,
							v1.VolumeMount{
								Name:      confMap.Name,
								ReadOnly:  true,
								MountPath: reqPath,
								SubPath:   req.Key,
							})
						break
					}
				}
			}
			if _, err = c.k8sClient.Do().AppsV1().Deployments(deployment.Namespace).Update(deployment); err != nil {
				_ = level.Error(c.logger).Log("Deployments", "Update", "err", err.Error())
			}
		}
	}

	return
}

//删除config map data
func (c *service) DeleteConfigMapData(ctx context.Context, req configMapDataRequest) (err error) {

	confData, _ := c.repository.ConfigData().FindById(req.ConfigMapDataId)
	confMap, notFound := c.repository.ConfigMap().FindById(confData.ConfigMapID)

	if notFound == true {
		return ErrConfigMapNotExist
	}

	//先同步远程configMapYaml 数据
	err = configmapyaml.SyncConfigMapYaml(confMap.Namespace, confMap.Name, c.logger, c.k8sClient, c.repository)

	if err != nil {
		return ErrSyncConfigMapYaml
	}

	//更新数据入库
	err = c.repository.ConfigData().DeleteById(req.ConfigMapDataId)
	if err != nil {
		return ErrDeleteConfigData
	}

	//更新远程
	err = configmapyaml.UpdateConfigMapYaml(confData.ConfigMapID, c.logger, c.k8sClient, c.repository)
	if err != nil {
		return ErrUpdateConfigMapYaml
	}

	return
}

// 获取config env data
func (c *service) GetConfigEnv(ctx context.Context, name string, ns string, page int, limit int) (res map[string]interface{}, err error) {
	cnt, err := c.repository.ConfigEnv().GetConfigEnvCountByNameNs(name, ns)
	if err != nil {
		_ = level.Error(c.logger).Log("GetConfigEnv", "GetConfigEnvCountByNameNs", "err", err.Error())
		return nil, ErrConfigEnvFailed
	}
	p := paginator.NewPaginator(page, limit, int(cnt))

	list, err := c.repository.ConfigEnv().GetConfigEnvPaginate(name, ns, p.Offset(), limit)

	if err != nil {
		_ = level.Error(c.logger).Log("GetConfigEnv", "GetConfigEnvPaginate", "err", err.Error())
		return nil, ErrConfigEnvFailed
	}
	res = map[string]interface{}{
		"list": list,
		"page": map[string]interface{}{
			"total":     cnt,
			"pageTotal": p.PageTotal(),
			"pageSize":  limit,
			"page":      p.Page(),
		},
	}
	return
}

func (c *service) CreateConfigEnv(ctx context.Context, req configEnvRequest) error {
	err := c.repository.ConfigEnv().CreateConfEnv(req.Name, req.Namespace, req.EnvKey, req.EnvVar, req.EnvDesc)
	if err != nil {
		_ = level.Error(c.logger).Log("CreateConfigEnv", "CreateConfEnv", "err", err.Error())
		return ErrCreateConfEnvFailed
	}

	if _, err = cronjob.ExchangeCronJobTemp(req.Name, req.Namespace, c.k8sClient, c.config, c.repository); err != nil {
		_ = level.Error(c.logger).Log("CreateConfigEnv", "ExchangeCronJobTemp", "err", err.Error())
		return ErrExchangeCronJobTemp
	}

	// TODO:: 通知
	//c.Notice("cronjob", cronjob)

	return nil
}

func (c *service) ConfigEnvUpdate(ctx context.Context, req configEnvRequest) error {
	confEnv, notFound := c.repository.ConfigEnv().FindById(req.Id)
	if notFound {
		_ = level.Error(c.logger).Log("ConfigEnvUpdate", "FindById", "err", "data not found")
		return ErrConfigEnvFailed
	}

	confEnv.EnvDesc = req.EnvDesc
	confEnv.EnvVar = req.EnvVar
	err := c.repository.ConfigEnv().Update(req.Id, confEnv)
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigEnvUpdate", "Update", "err", err.Error())
		return ErrUpdateConfEnvFailed
	}
	if _, err := cronjob.ExchangeCronJobTemp(req.Name, req.Namespace, c.k8sClient, c.config, c.repository); err != nil {
		_ = level.Error(c.logger).Log("ConfigEnvUpdate", "ExchangeCronJobTemp", "err", err.Error())
		return ErrExchangeCronJobTemp
	}

	// TODO:: 通知
	//c.Notice("cronjob", cronjob)

	return nil
}

func (c *service) ConfigEnvDel(ctx context.Context, req configEnvRequest) error {
	confEnv, notFound := c.repository.ConfigEnv().FindById(req.Id)
	if notFound {
		_ = level.Error(c.logger).Log("ConfigEnvDel", "FindById", "err", "data not found")
		return ErrConfigEnvFailed
	}

	err := c.repository.ConfigEnv().Delete(req.Id)
	if err != nil {
		_ = level.Error(c.logger).Log("ConfigEnvDel", "Delete", "err", err.Error())
		return ErrConfigEnvFailed
	}

	if _, err = cronjob.ExchangeCronJobTemp(confEnv.Name, confEnv.Namespace, c.k8sClient, c.config, c.repository); err != nil {
		_ = level.Error(c.logger).Log("ConfigEnvDel", "ExchangeCronJobTemp", "err", err.Error())
		return ErrExchangeCronJobTemp
	}

	// TODO:: 通知
	//c.Notice("cronjob", cronjob)

	return nil
}
