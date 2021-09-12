package namespace

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

type Middleware func(Service) Service

type Service interface {
	// Sync 同步namespace
	Sync(ctx context.Context, clusterId int64) error
	// Create 创建空间 如果有 imageSecrets 则下发
	Create(ctx context.Context, clusterId int64, name, alias, remark string, imageSecrets []string) (err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func (s *service) Create(ctx context.Context, clusterId int64, name, alias, remark string, imageSecrets []string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	byName, err := s.repository.Namespace(ctx).FindByName(ctx, clusterId, name)
	if err == nil {
		_ = level.Error(logger).Log("repository.Namespace", "FindByName", "err", err.Error())
		return encode.ErrNamespaceExists.Error()
	}
	if !gorm.IsRecordNotFoundError(err) {
		return encode.ErrNamespaceExists.Wrap(err)
	}
	byName.Name = name
	byName.Alias = alias
	byName.Remark = remark
	byName.ClusterId = clusterId
	if err = s.repository.Namespace(ctx).SaveCall(ctx, &byName, func() error {
		_, err = s.k8sClient.Do(ctx).CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			_ = level.Error(logger).Log("k8sClient.Do.CoreV1.Namespaces", "Create", "err", err.Error())
			return encode.ErrNamespaceCreate.Wrap(err)
		}
		return nil
	}); err != nil {
		_ = level.Error(logger).Log("k8sClient.Do.CoreV1.Namespaces", "Create", "err", err.Error())
		return encode.ErrNamespaceCreate.Wrap(err)
	}
	if len(imageSecrets) < 1 {
		return nil
	}
	// 下发获取镜像的Secrets
	names, err := s.repository.Registry(ctx).FindByNames(ctx, imageSecrets)
	if err != nil {
		_ = level.Error(logger).Log("repository.Registry", "FindByNames", "err", err.Error())
		return nil
	}
	for _, v := range names {
		auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", v.Username, v.Password)))
		marshal, err := json.Marshal(map[string]interface{}{
			v.Host: map[string]string{
				"username": v.Username,
				"password": v.Password,
				"auth":     auth,
			},
		})
		if err != nil {
			_ = level.Error(logger).Log("json", "Marshal", "err", err.Error())
			continue
		}
		coreV1Secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: name,
				Name:      v.Name,
			},
			Data: map[string][]byte{
				corev1.DockerConfigKey: marshal,
			},
			Type: corev1.SecretTypeDockercfg,
		}
		_, err = s.k8sClient.Do(ctx).CoreV1().Secrets(name).Create(ctx, &coreV1Secret, metav1.CreateOptions{})
		if err != nil {
			_ = level.Error(logger).Log("k8sClient.Do", "CoreV1", "Secrets", "Create", "err", err.Error())
		} else {
			var secret types.Secret
			secret.Namespace = name
			secret.Name = v.Name
			secret.ClusterId = clusterId
			secret.ResourceVersion = coreV1Secret.ResourceVersion
			data := types.Data{
				Style: types.DataStyleSecret,
				Key:   corev1.DockerConfigKey,
				Value: string(marshal),
			}
			if err = s.repository.Secrets(ctx).Save(ctx, &secret, []types.Data{data}); err != nil {
				err = encode.ErrSecretImageSave.Wrap(err)
				_ = level.Error(logger).Log("repository.Secrets", "Save", "err", err.Error())
			}
		}
	}
	return nil
}

func (s *service) Sync(ctx context.Context, clusterId int64) error {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	if namespaces, err := s.k8sClient.Do(ctx).CoreV1().Namespaces().List(ctx, metav1.ListOptions{}); err == nil {
		for _, v := range namespaces.Items {
			ns, err := s.repository.Namespace(ctx).FindByName(ctx, clusterId, v.Name)
			if err != nil {
				if !gorm.IsRecordNotFoundError(err) {
					_ = level.Error(logger).Log("repository.Namespace", "FindByName", "err", err.Error())
					err = encode.ErrNamespaceNotfound.Error()
					return err
				}
				ns = types.Namespace{}
			}
			ns.ClusterId = clusterId
			ns.Name = v.Name
			ns.Alias = v.Name
			ns.Status = string(v.Status.Phase)
			if err = s.repository.Namespace(ctx).Save(ctx, &ns); err != nil {
				_ = level.Error(logger).Log("repository.Namespace", "Save", "err", err.Error())
			}
		}
	} else {
		_ = level.Error(logger).Log("k8sClient.Do.CoreV1.Namespaces", "List", "err", err.Error())
	}

	return nil
}

func New(logger log.Logger, traceId string, client kubernetes.K8sClient, store repository.Repository) Service {
	return &service{
		traceId:    traceId,
		logger:     logger,
		k8sClient:  client,
		repository: store,
	}
}
