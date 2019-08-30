/**
 * Created by GoLand.
 * Email: xzghua@gmail.com
 * Date: 2019-08-09
 * Time: 19:04
 */
package configmapyaml

import (
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ErrConfigMapNotExist = errors.New("configMap 不存在")
	ErrConfigMapDataList = errors.New("configDataList 不存在")
)

//修改configMap前，先同步远程yaml 到 db, 再修改后，更新到远程
func SyncConfigMapYaml(ns, name string,
	logger log.Logger,
	k8sClient kubernetes.K8sClient, repository repository.Repository) (err error) {

	confMapYaml, err := k8sClient.Do().CoreV1().ConfigMaps(ns).Get(name, metav1.GetOptions{})
	confMap, notFound := repository.ConfigMap().Find(ns, name)
	if notFound == true {
		//add
		re, err := repository.ConfigMap().Create(&types.ConfigMap{
			Namespace: ns,
			Name:      name,
			Desc:      name,
		})

		if err == nil && re.ID > 0 {
			for key, val := range confMapYaml.Data {
				err := repository.ConfigData().Create(&types.ConfigData{
					ConfigMap: *re,
					Key:       key,
					Value:     val,
				})
				if err != nil {
					_ = logger.Log("configMap", "SyncConfigMapYaml", "err", err.Error())
					return err
				}
			}

		}
	} else {
		//update
		for key, val := range confMapYaml.Data {
			confData, notFound := repository.ConfigData().FindByConfMapIdAndKey(confMap.ID, key)
			if !notFound {
				//update
				go repository.ConfigData().Update(confData.ID, val, confData.Path)
			} else {
				//add
				d := &types.ConfigData{
					Key:         key,
					Value:       val,
					ConfigMapID: confMap.ID,
				}
				go repository.ConfigData().Create(d)
			}
		}
	}
	return
}

//configMapData key-value 变化时，更新configMap yaml
func UpdateConfigMapYaml(confMapId int64,
	logger log.Logger,
	k8sClient kubernetes.K8sClient, repository repository.Repository) (err error) {
	confMap, notFound := repository.ConfigMap().FindById(confMapId)
	if notFound == true {
		_ = logger.Log("configMap", "UpdateConfigMapYaml", "err", "data not found")
		return ErrConfigMapNotExist
	}

	confDataList, err := repository.ConfigData().Find(confMap.Namespace, confMap.Name)
	if err != nil {
		_ = logger.Log("configMap", "UpdateConfigMapYaml", "err", err.Error())
		return ErrConfigMapDataList
	}

	cdl := map[string]string{}
	for _, v := range confDataList {
		cdl[v.Key] = v.Value
	}

	//create configMap yaml
	conf := new(v1.ConfigMap)
	conf.Namespace = confMap.Namespace
	conf.Name = confMap.Name
	conf.Data = cdl

	//判断远程是否已存在
	_, err = k8sClient.Do().CoreV1().ConfigMaps(confMap.Namespace).Get(confMap.Name, metav1.GetOptions{})

	if err != nil {
		_, err = k8sClient.Do().CoreV1().ConfigMaps(confMap.Namespace).Create(conf)
		return
	} else {
		_, err = k8sClient.Do().CoreV1().ConfigMaps(confMap.Namespace).Update(conf)
		return
	}
}
