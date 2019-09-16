/**
 * @Time : 2019-07-29 09:59
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package market

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"gopkg.in/guregu/null.v3"
)

type Service interface {
	// 创建dockerfile
	Post(ctx context.Context, name, language, version, detail, desc, dockerfile, fullPath string, status int64) (err error)

	// dockerfile详情
	Detail(ctx context.Context, id int64) (res *types.Dockerfile, err error) // download 可以调用这个方法

	// dockerfile列表
	List(ctx context.Context, page, limit int, language []string, status int, name string) (res []*types.Dockerfile, count int64, err error)

	// 更新dockerfile
	Put(ctx context.Context, id int64, name, language, version, detail, desc, dockerfile, fullPath string, status int64) (err error)

	// 删除dockerfile
	Delete(ctx context.Context, id int64) (err error)
}

var (
	ErrDockerfileCreate = errors.New("创建dockerfile错误")
	ErrDockerfileGet    = errors.New("获取dockerfile错误")
	ErrDockerfileUpdate = errors.New("dockerfile更新错误")
	//ErrDockerfileListKeyValue = errors.New("获取获取错误,查询的key value不正确")
)

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, store repository.Repository) Service {
	return &service{logger, store}
}

func (c *service) Post(ctx context.Context, name, language, version, detail, desc, dockerfile, fullPath string, status int64) (err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)

	if err = c.repository.Dockerfile().Create(&types.Dockerfile{
		Name:       name,
		Language:   language,
		Version:    version,
		Detail:     detail,
		Desc:       null.StringFrom(desc),
		Dockerfile: dockerfile,
		FullPath:   fullPath,
		Status:     null.IntFrom(status),
		AuthorID:   userId,
		Sha256:     null.StringFrom(encode.HashString([]byte(dockerfile))),
	}); err != nil {
		_ = level.Error(c.logger).Log("repository.Dockerfile().", "Create", "err", err.Error())
		return ErrDockerfileCreate
	}

	return
}

func (c *service) Detail(ctx context.Context, id int64) (res *types.Dockerfile, err error) {
	return c.repository.Dockerfile().FindById(id)
}

func (c *service) List(ctx context.Context, page, limit int, language []string, status int, name string) (res []*types.Dockerfile, count int64, err error) {

	return c.repository.Dockerfile().FindBy(language, status, name, (page-1)*limit, limit)
}

func (c *service) Put(ctx context.Context, id int64, name, language, version, detail, desc, dockerfile, fullPath string, status int64) (err error) {
	data, err := c.repository.Dockerfile().FindById(id)
	if err != nil {
		_ = level.Error(c.logger).Log("dockerfileRepository", "FindById", "err", err.Error())
		return ErrDockerfileGet
	}

	data.Name = name
	data.Language = language
	data.Version = version
	data.Desc = null.StringFrom(desc)
	data.Dockerfile = dockerfile
	data.FullPath = fullPath
	data.Detail = detail
	data.Sha256 = null.StringFrom(encode.HashString([]byte(dockerfile)))

	if err = c.repository.Dockerfile().Update(data); err != nil {
		_ = level.Error(c.logger).Log("dockerfileRepository", "Update", "err", err.Error())
		return ErrDockerfileUpdate
	}

	return
}

func (c *service) Delete(ctx context.Context, id int64) (err error) {
	return c.repository.Dockerfile().Delete(id)
}
