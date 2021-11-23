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

// Service 项目空间模块
// 目前该模块设计为所有成员都可以查看
// TODO: 考虑实现 同步空间下的Deployment、CronJob ConfigMap Service Ingress 或直接写一下Sync模块处理所有相关同步？
// TODO: 以上想法待考虑
type Service interface {
	// Sync 同步集群的namespace
	Sync(ctx context.Context, clusterId int64) error
	// Create 创建空间 如果有 imageSecrets 则下发,imageSecrets可以传多个
	Create(ctx context.Context, clusterId int64, name, alias, remark string, imageSecrets []string) (err error)
	// List 空间列表
	// names从中间件取得该用户有多少个空间的权限
	List(ctx context.Context, clusterId int64, names []string, query string, page, pageSize int) (res []result, total int, err error)
	// Delete 删除空间
	// 需要判断该空间下是否还有各个类型的资源，如果有并且force为false的话不给操作
	// 如果force为true的话,强制清理所有项目？这风险会不会太大了？待定
	// 直接中间件判定是否有该中间件的操作权限
	Delete(ctx context.Context, clusterId int64, name string, force bool) (err error)
	// Update 更新空间
	// 只允许更新别名和备注，其他的无法操作
	// 直接中间件判定是否有该中间件的操作权限
	Update(ctx context.Context, clusterId int64, name string, alias, remark, status string, imageSecrets []string) (err error)
	// Info 获取空间的基本信息
	Info(ctx context.Context, clusterId int64, name string) (res result, err error)
	// IssueSecret 生成镜像密钥，并下发至各个空间的Secrets
	IssueSecret(ctx context.Context, clusterId int64, name, regName string) (err error)
	// ReloadSecret 重新加载所有密钥
	// 清除之前的密钥信息，生成新的密钥并发布到Secrets
	ReloadSecret(ctx context.Context, clusterId int64, name, regName string) (err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func (s *service) IssueSecret(ctx context.Context, clusterId int64, name, regName string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	ns, err := s.repository.Namespace(ctx).FindByName(ctx, clusterId, name)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Namespace", "FindByName", "err", err.Error())
		err = encode.ErrNamespaceNotfound.Error()
		return
	}
	reg, err := s.repository.Registry(ctx).FindByName(ctx, regName)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Registry", "FindByName", "err", err.Error())
		err = encode.ErrRegistryNotfound.Error()
		return err
	}

	fmt.Println(reg.Host, reg.Name, ns.Name, ns.Alias)

	return
}

func (s *service) ReloadSecret(ctx context.Context, clusterId int64, name, regName string) (err error) {
	panic("implement me")
}

func (s *service) Info(ctx context.Context, clusterId int64, name string) (res result, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	ns, err := s.repository.Namespace(ctx).FindByName(ctx, clusterId, name)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Namespace", "FindByName", "err", err.Error())
		err = encode.ErrNamespaceNotfound.Error()
		return
	}

	// 获取空间证书
	regs, _, err := s.repository.Registry(ctx).List(ctx, "", 1, 50)
	if err != nil {
		_ = level.Error(logger).Log("repository.Registry", "List", "err", err.Error())
		return
	}
	var names, secrets []string
	for _, v := range regs {
		names = append(names, v.Name)
	}
	res.RegSecrets = names
	secretList, err := s.repository.Secrets(ctx).FindNsByNames(ctx, clusterId, ns.Name, names)
	if err != nil {
		_ = level.Error(logger).Log("repository.Registry", "List", "err", err.Error())
		return
	}
	for _, v := range secretList {
		secrets = append(secrets, v.Name)
	}
	res.ImageSecrets = secrets

	res.Name = ns.Name
	res.Alias = ns.Alias
	res.Remark = ns.Remark
	//res.Status = ns.Status
	res.UpdatedAt = ns.UpdatedAt
	res.CreatedAt = ns.CreatedAt

	info, err := s.k8sClient.Do(ctx).CoreV1().Namespaces().Get(ctx, ns.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Warn(logger).Log("k8sClient.Do.CoreV1.Namespaces", "Get", "err", err.Error())
		err = encode.ErrNamespaceNotfound.Error()
		return
	}

	res.Status = string(info.Status.Phase)
	return
}

func (s *service) Delete(ctx context.Context, clusterId int64, name string, force bool) (err error) {
	panic("implement me")
}

func (s *service) Update(ctx context.Context, clusterId int64, name string, alias, remark, status string, imageSecrets []string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	ns, err := s.repository.Namespace(ctx).FindByName(ctx, clusterId, name)
	if err != nil {
		_ = level.Error(logger).Log("repository.Namespace", "FindByName", "err", err.Error())
		err = encode.ErrNamespaceNotfound.Error()
		return
	}

	ns.Alias = alias
	ns.Remark = remark

	err = s.repository.Namespace(ctx).Save(ctx, &ns)
	if err != nil {
		_ = level.Error(logger).Log("repository.Namespace", "Save", "err", err.Error())
		err = encode.ErrNamespaceSave.Wrap(err)
		return
	}

	return
}

func (s *service) List(ctx context.Context, clusterId int64, names []string, query string, page, pageSize int) (res []result, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	list, total, err := s.repository.Namespace(ctx).List(ctx, clusterId, names, query, page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.Namespace", "List", "err", err.Error())
		err = encode.ErrNamespaceList.Wrap(err)
		return
	}

	// 获取所有的registry
	// 根据registryName和namespaces 获取secrets的 TODO: 似乎是太麻烦了

	for _, v := range list {
		res = append(res, result{
			Name:      v.Name,
			Alias:     v.Alias,
			Remark:    v.Remark,
			Status:    v.Status,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return
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
			return encode.ErrNamespaceSave.Wrap(err)
		}
		return nil
	}); err != nil {
		_ = level.Error(logger).Log("k8sClient.Do.CoreV1.Namespaces", "Create", "err", err.Error())
		return encode.ErrNamespaceSave.Wrap(err)
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
	// TODO: 需要给自己授权
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
