package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

type Middleware func(Service) Service

// Service 镜像仓库管理模块
type Service interface {
	// Create 创建空间
	Create(ctx context.Context, name, host, username, password, remark string) (err error)
	// List 仓库列表
	List(ctx context.Context, query string, page, pageSize int) (res []result, total int, err error)
	// Update 更新仓库中信息
	// name允许更新
	// 需要更新相应的secrets，需要便利每个namespace 并更新相应的secrets
	Update(ctx context.Context, name, host, username, password, remark string) (err error)
	// Delete 删除仓库
	// 需要便利每个namespace 并删除相应的secrets
	Delete(ctx context.Context, name string) (err error)
	// Password 获取仓库密码 只有管理员可以查看，在中间件处理就行了
	Password(ctx context.Context, name string) (res string, err error)
	// Info 获取仓库详情，如果仓库类型是harbor，可以调用api获取更多的信息
	Info(ctx context.Context, name string) (res result, err error)
	// Secret 下发取各个空间下的Secrets
	Secret(ctx context.Context, name string) (err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) Secret(ctx context.Context, name string) (err error) {
	panic("implement me")
}

func (s *service) Update(ctx context.Context, name, host, username, password, remark string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	reg, err := s.repository.Registry(ctx).FindByName(ctx, name)
	if err != nil {
		err = encode.ErrRegistryNotfound.Error()
		return
	}
	var oldReg types.Registry
	oldReg = reg
	reg.Host = host
	reg.Username = username
	reg.Password = password
	reg.Remark = remark
	if err = s.repository.Registry(ctx).SaveCall(ctx, &reg, func() error {
		secrets, err := s.repository.Secrets(ctx).FindByName(ctx, oldReg.Name)
		if err != nil {
			return err
		}
		for _, v := range secrets {
			auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
			marshal, _ := json.Marshal(map[string]interface{}{
				host: map[string]string{
					"username": username,
					"password": password,
					"auth":     auth,
				},
			})

			var dataList []types.Data
			for _, vv := range v.Data {
				if vv.Style != types.DataStyleSecret {
					continue
				}
				if !strings.EqualFold(vv.Key, corev1.DockerConfigKey) {
					continue
				}
				vv.Value = string(marshal)
				dataList = append(dataList, vv)
			}
			// 有可能其中几个会保存失败，先不管它
			if err := s.repository.Secrets(ctx).Save(ctx, &v, dataList); err != nil {
				_ = level.Error(logger).Log("repository.Secrets", "Save", "err", err.Error())
			}
		}
		return nil
	}); err != nil {
		_ = level.Error(logger).Log("repository.Registry", "SaveCall", "err", err.Error())
		return encode.ErrRegistryUpdate.Wrap(err)
	}
	return
}

func (s *service) Delete(ctx context.Context, name string) (err error) {
	reg, err := s.repository.Registry(ctx).FindByName(ctx, name)
	if err != nil {
		err = encode.ErrRegistryNotfound.Error()
		return
	}

	// 遍历所有空间，删除对应的secret
	fmt.Println(reg.Name)

	return
}

func (s *service) Password(ctx context.Context, name string) (res string, err error) {
	reg, err := s.repository.Registry(ctx).FindByName(ctx, name)
	if err != nil {
		err = encode.ErrRegistryNotfound.Error()
		return
	}
	return reg.Password, nil
}

func (s *service) Info(ctx context.Context, name string) (res result, err error) {
	reg, err := s.repository.Registry(ctx).FindByName(ctx, name)
	if err != nil {
		err = encode.ErrRegistryNotfound.Error()
		return
	}
	fmt.Println(reg.Host)
	return
}

func (s *service) List(ctx context.Context, query string, page, pageSize int) (res []result, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	list, total, err := s.repository.Registry(ctx).List(ctx, query, page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.Registry", "List", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, result{
			Name:      v.Name,
			Host:      v.Host,
			Username:  v.Username,
			Password:  "",
			Remark:    v.Remark,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return
}

func (s *service) Create(ctx context.Context, name, host, username, password, remark string) (err error) {
	return s.repository.Registry(ctx).Save(ctx, &types.Registry{
		Name:     name,
		Host:     host,
		Username: username,
		Password: password,
		Remark:   remark,
	})
}

func New(logger log.Logger, traceId string, store repository.Repository) Service {
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: store,
	}
}
