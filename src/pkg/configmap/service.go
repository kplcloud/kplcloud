/**
 * @Time: 2021/8/18 23:03
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package configmap

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type Middleware func(Service) Service

type Service interface {
	Sync(ctx context.Context, clusterId int64, ns string) (err error)

	// List 配置字典列表
	List(ctx context.Context, clusterId int64, ns, name string, page, pageSize int) (res []configMapResult, total int, err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) List(ctx context.Context, clusterId int64, ns, name string, page, pageSize int) (res []configMapResult, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	list, total, err := s.repository.ConfigMap(ctx).List(ctx, clusterId, ns, name, page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.ConfigMap", "List", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, configMapResult{
			Name:      v.Name,
			Namespace: v.Namespace,
			Desc:      v.Desc,
			Version:   v.ResourceVersion,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return
}

func (s *service) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	var nss []string
	if strings.EqualFold(ns, "") {
		if list, err := s.repository.Namespace(ctx).FindByCluster(ctx, clusterId); err == nil {
			for _, v := range list {
				nss = append(nss, v.Name)
			}
		} else {
			_ = level.Error(logger).Log("repository.Namespace", "FindByCluster", "err", err.Error())
			return err
		}
	} else {
		nss = append(nss, ns)
	}

	for _, n := range nss {
		list, err := s.k8sClient.Do(ctx).CoreV1().ConfigMaps(n).List(ctx, metav1.ListOptions{})
		if err != nil {
			_ = level.Error(logger).Log("k8sClient.Do.CoreV1.ConfigMaps", "List", "err", err.Error())
			//return encode.ErrConfigMapSyncList.Wrap(err)
			continue
		}

		for _, v := range list.Items {
			cfgMap := &types.ConfigMap{
				ClusterId:       clusterId,
				Name:            v.Name,
				Namespace:       v.Namespace,
				ResourceVersion: v.ResourceVersion,
			}
			var data []types.Data
			for key, val := range v.Data {
				data = append(data, types.Data{
					Style: types.DataStyleConfigMap,
					Key:   key,
					Value: val,
				})
			}
			if err = s.repository.ConfigMap(ctx).Save(ctx, cfgMap, data); err != nil {
				_ = level.Error(logger).Log("repository.ConfigMap", "Save", "err", err.Error())
				continue
			}
		}
	}

	return nil
}

func New(logger log.Logger, traceId string, repository repository.Repository, client kubernetes.K8sClient) Service {
	logger = log.With(logger, "configmap", "service")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
		k8sClient:  client,
	}
}
