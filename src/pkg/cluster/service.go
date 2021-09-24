/**
 * @Time : 8/9/21 6:20 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package cluster

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"strings"
)

type Middleware func(Service) Service

// Service 集群模块
type Service interface {
	// Add 添加集群
	Add(ctx context.Context, name, alias, data string) (err error)
	// List 集群列表
	List(ctx context.Context, name string, page, pageSize int) (res []listResult, total int, err error)
	SyncRoles(ctx context.Context, clusterId int64) (err error)
}

type service struct {
	k8sClient  kubernetes.K8sClient
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) List(ctx context.Context, name string, page, pageSize int) (res []listResult, total int, err error) {
	list, total, err := s.repository.Cluster(ctx).List(ctx, name, 0, page, pageSize)
	if err != nil {
		err = encode.ErrClusterNotfound.Wrap(errors.Wrap(err, "repository.Cluster.List"))
		return
	}
	for _, v := range list {
		res = append(res, listResult{
			Name:      v.Name,
			Alias:     v.Alias,
			Remark:    v.Remark,
			Status:    v.Status,
			NodeNum:   v.NodeNum,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return
}

func (s *service) SyncRoles(ctx context.Context, clusterId int64) (err error) {
	//logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	list, err := s.k8sClient.Do(ctx).RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		err = errors.Wrap(err, "k8sClient.Do.RbacV1.ClusterRoles.List")
		return err
	}
	for _, v := range list.Items {
		var rules []types.PolicyRule
		for _, p := range v.Rules {
			rules = append(rules, types.PolicyRule{
				Kind:            v.Kind,
				Verbs:           strings.Join(p.Verbs, ","),
				APIGroups:       strings.Join(p.APIGroups, ","),
				Resources:       strings.Join(p.Resources, ","),
				ResourceNames:   strings.Join(p.ResourceNames, ","),
				NonResourceURLs: strings.Join(p.NonResourceURLs, ","),
			})
		}
		b, _ := json.Marshal(v)
		clusterRule := &types.ClusterRole{
			ClusterId: clusterId,
			Name:      v.Name,
			Data:      string(b),
			Rules:     rules,
		}
		err := s.repository.Cluster(ctx).SaveRole(ctx, clusterRule, rules)
		if err != nil {
			return err
		}
	}
	return
}

func (s *service) Add(ctx context.Context, name, alias, data string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	cluster := types.Cluster{
		Name:       name,
		Alias:      alias,
		Status:     2,
		ConfigData: data,
	}

	if err = s.repository.Cluster(ctx).Save(ctx, &cluster, func(tx *gorm.DB) error {
		if err = s.k8sClient.Connect(ctx, name, data); err != nil {
			_ = level.Error(logger).Log("k8sClient.Connect", "err", err.Error())
			return encode.ErrClusterConnect.Error()
		}
		return nil
	}); err != nil {
		_ = level.Error(logger).Log("repository.Cluster", "Save", "err", err.Error())
		return encode.ErrClusterAdd.Error()
	}

	return
}

func New(logger log.Logger, traceId string, repository repository.Repository, k8sClient kubernetes.K8sClient) Service {
	logger = log.With(logger, "cluster", "service")
	return &service{
		k8sClient:  k8sClient,
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}
