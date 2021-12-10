/**
 * @Time : 8/19/21 1:36 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package secret

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

type Middleware func(Service) Service

type Service interface {
	// Sync 同步所有Secret
	Sync(ctx context.Context, clusterId int64, ns string) (err error)
	// ImageSecret 添加更新拉取镜像的Secret
	ImageSecret(ctx context.Context, clusterId int64, ns, name, host, username, password string) (err error)
	// Delete 删除Secret
	Delete(ctx context.Context, clusterId int64, ns, name string) (err error)
	// List 列表查询
	List(ctx context.Context, clusterId int64, namespace, name string, page, pageSize int) (res []secretResult, total int, err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) List(ctx context.Context, clusterId int64, namespace, name string, page, pageSize int) (res []secretResult, total int, err error) {
	list, total, err := s.repository.Secrets(ctx).List(ctx, clusterId, namespace, name, page, pageSize)
	if err != nil {
		return
	}

	for _, v := range list {
		res = append(res, secretResult{
			Name:            v.Name,
			Namespace:       v.Namespace,
			ResourceVersion: v.ResourceVersion,
			CreatedAt:       v.CreatedAt,
			UpdatedAt:       v.UpdatedAt,
		})
	}

	return
}

func (s *service) Delete(ctx context.Context, clusterId int64, ns, name string) (err error) {
	err = s.k8sClient.Do(ctx).CoreV1().Secrets(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		err = errors.Wrap(err, "k8s.Do.CoreV1.Secrets.Delete")
		return encode.ErrSecretDelete.Wrap(err)
	}

	if err = s.repository.Secrets(ctx).Delete(ctx, clusterId, ns, name); err != nil {
		err = errors.Wrap(err, "repository.Secrets.Delete")
		return encode.ErrSecretDelete.Wrap(err)
	}

	return err
}

func (s *service) ImageSecret(ctx context.Context, clusterId int64, ns, name, host, username, password string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	marshal, err := json.Marshal(map[string]interface{}{
		host: map[string]string{
			"username": username,
			"password": password,
			"auth":     auth,
		},
	})
	if err != nil {
		return encode.ErrSecretMarshal.Wrap(err)
	}

	coreV1Secret, err := s.k8sClient.Do(ctx).CoreV1().Secrets(ns).Get(ctx, name, metav1.GetOptions{})
	coreV1Secret.Namespace = ns
	coreV1Secret.Name = name
	coreV1Secret.Type = corev1.SecretTypeDockercfg
	coreV1Secret.Data = map[string][]byte{
		corev1.DockerConfigKey: marshal,
	}
	if err != nil {
		_ = level.Warn(logger).Log("k8sClient.Do", "CoreV1", "Secrets", "Get", "err", err.Error())
		coreV1Secret, err = s.k8sClient.Do(ctx).CoreV1().Secrets(ns).Create(ctx, coreV1Secret, metav1.CreateOptions{})
		if err != nil {
			return encode.ErrSecretImageSave.Wrap(err)
		}
	} else {
		coreV1Secret, err = s.k8sClient.Do(ctx).CoreV1().Secrets(ns).Update(ctx, coreV1Secret, metav1.UpdateOptions{})
		if err != nil {
			return encode.ErrSecretImageSave.Wrap(err)
		}
	}

	var secret types.Secret
	secret.Namespace = ns
	secret.Name = name
	secret.ClusterId = clusterId
	secret.ResourceVersion = coreV1Secret.ResourceVersion
	data := types.Data{
		Style: types.DataStyleSecret,
		Key:   corev1.DockerConfigKey,
		Value: string(marshal),
	}

	if err = s.repository.Secrets(ctx).Save(ctx, &secret, []types.Data{data}); err != nil {
		return encode.ErrSecretImageSave.Wrap(err)
	}

	return nil
}

func (s *service) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	var items *corev1.SecretList
	if items, err = s.k8sClient.Do(ctx).CoreV1().Secrets(ns).List(ctx, metav1.ListOptions{}); err != nil {
		_ = level.Error(logger).Log("k8sClient.Do.AppsV1.Secrets", "List", "err", err.Error())
		return encode.ErrDeploymentSyncList.Wrap(err)
	}

	for _, v := range items.Items {
		var data []types.Data
		for key, val := range v.Data {
			data = append(data, types.Data{
				Style: types.DataStyleSecret,
				Key:   key,
				Value: string(val),
			})
		}

		if err = s.repository.Secrets(ctx).Save(ctx, &types.Secret{
			ClusterId:       clusterId,
			Name:            v.Name,
			Namespace:       v.Namespace,
			ResourceVersion: v.ResourceVersion,
		}, data); err != nil {
			_ = level.Error(logger).Log("repository.Secrets", "Save", "err", err.Error())
		}
	}

	return nil
}

func New(logger log.Logger, traceId string, repository repository.Repository, client kubernetes.K8sClient) Service {
	logger = log.With(logger, "secret", "service")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
		k8sClient:  client,
	}
}
