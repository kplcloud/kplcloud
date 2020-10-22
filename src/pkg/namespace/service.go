package namespace

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrNamespaceIsExists = errors.New("空间已经存在.")
	ErrNamespaceCreate   = errors.New("空间创建失败.")
	ErrNamespaceList     = errors.New("空间获取获取失败.")
)

type Service interface {
	// 详情信息
	Get(ctx context.Context, name string) (resp *types.Namespace, err error)

	// 创建namespace
	Post(ctx context.Context, name, displayName string) error

	// 同步namespace
	Sync(ctx context.Context) error

	// 删除namespace
	Delete(ctx context.Context) error

	// 更新Namespaces
	Update(ctx context.Context, name, displayName string) error

	// 空间列表
	List(ctx context.Context) (res []*types.Namespace, err error)
}

type service struct {
	logger     log.Logger
	config     *config.Config
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func (c *service) List(ctx context.Context) (res []*types.Namespace, err error) {
	return c.repository.Namespace().FindAll()
}

func (c *service) Update(ctx context.Context, name, displayName string) error {

	return nil
}

func (c *service) Get(ctx context.Context, name string) (resp *types.Namespace, err error) {
	return c.repository.Namespace().Find(name)
}

func (c *service) Post(ctx context.Context, name, displayName string) error {
	res, err := c.repository.Namespace().Find(name)
	if err == nil || res != nil {
		_ = level.Error(c.logger).Log("displayName", name)
		return ErrNamespaceIsExists
	}

	namespace := new(v1.Namespace)
	namespace.Name = name

	if _, err := c.k8sClient.Do().CoreV1().Namespaces().Create(namespace); err != nil {
		_ = level.Error(c.logger).Log("k8s", "create", "err", err.Error())
		return ErrNamespaceCreate
	}

	if err = c.repository.Namespace().Create(&types.Namespace{
		DisplayName: displayName,
		Name:        name,
	}); err != nil {
		_ = level.Error(c.logger).Log("ns", "create", "err", err.Error())
		return ErrNamespaceCreate
	}

	// 是否需要创建secret
	if secrets := c.config.GetString("kubernetes", "image_pull_secrets"); secrets != "" {
		if _, err := c.k8sClient.Do().CoreV1().Secrets(name).Create(&v1.Secret{
			Type: v1.SecretTypeDockerConfigJson,
			ObjectMeta: metav1.ObjectMeta{
				Name:      secrets,
				Namespace: name,
			},
		}); err != nil {
			_ = level.Warn(c.logger).Log("k8s", "secrets", "err", err.Error())
		}
	}

	// 如果有配置jenkins的话，创建jenkins 视图

	return nil
}

func (c *service) Sync(ctx context.Context) error {

	nsList, err := c.k8sClient.Do().CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("namespace", "list", "err", err.Error())
		return ErrNamespaceList
	}

	for _, ns := range nsList.Items {
		if info, err := c.repository.Namespace().Find(ns.Name); err != nil || info == nil {
			if err = c.repository.Namespace().Create(&types.Namespace{
				DisplayName: ns.Name,
				Name:        ns.Name,
			}); err != nil {
				_ = level.Warn(c.logger).Log("namespace", "create", "err", err.Error())
			}
		}
	}

	return nil
}

func (c *service) Delete(ctx context.Context) error {

	return nil
}

func NewService(logger log.Logger, cf *config.Config, client kubernetes.K8sClient, store repository.Repository) Service {
	return &service{
		logger:     logger,
		config:     cf,
		k8sClient:  client,
		repository: store,
	}
}
