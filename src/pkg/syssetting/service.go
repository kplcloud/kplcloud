/**
 * @Time : 7/20/21 5:41 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package syssetting

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/repository"
)

type Middleware func(Service) Service

type Service interface {
	// List 列表
	List(ctx context.Context, key string, page, pageSize int) (res []listResult, total int, err error)
	// Delete 删除设置
	Delete(ctx context.Context, id int64) (err error)
	// Add 添加设置
	Add(ctx context.Context, section, key, value, remark string, enable bool) (err error)
	// Update 更新设置
	Update(ctx context.Context, id int64, section, key, value, remark string, enable bool) (err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) Add(ctx context.Context, section, key, value, remark string, enable bool) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	_, err = s.repository.SysSetting().Find(ctx, section, key)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			_ = level.Error(logger).Log("repository.SysSetting", "Find", "err", err.Error())
			return encode.ErrSystem.Wrap(errors.New("repository.SysSetting.Find"))
		}
	} else {
		return encode.ErrSysSettingExists.Error()
	}

	if err = s.repository.SysSetting().Add(ctx, section, key, value, remark); err != nil {
		_ = level.Error(logger).Log("repository.SysSetting", "Add", "err", err.Error())
		return encode.ErrSysSettingSave.Error()
	}

	return nil
}

func (s *service) Update(ctx context.Context, id int64, section, key, value, remark string, enable bool) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	setting, err := s.repository.SysSetting().FindById(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysSetting", "Find", "err", err.Error())
		if !gorm.IsRecordNotFoundError(err) {
			return encode.ErrSystem.Wrap(errors.New("repository.SysSetting.Find"))
		}
		return encode.ErrSysSettingNotfound.Error()
	}

	setting.Section = section
	setting.Key = key
	setting.Value = value
	setting.Description = remark
	setting.Enable = enable
	if err = s.repository.SysSetting().Update(ctx, &setting); err != nil {
		_ = level.Error(logger).Log("repository.SysSetting", "Add", "err", err.Error())
		return encode.ErrSysSettingSave.Error()
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	setting, err := s.repository.SysSetting().FindById(ctx, id)
	if err != nil {
		return encode.ErrSysSettingNotfound.Error()
	}
	err = s.repository.SysSetting().Delete(ctx, setting.Section, setting.Key)
	if err != nil {
		_ = level.Error(logger).Log("repository.SysSetting", "Delete", "err", err.Error())
		return encode.ErrSysSettingDelete.Error()
	}

	return nil
}

func (s *service) List(ctx context.Context, key string, page, pageSize int) (res []listResult, total int, err error) {
	list, total, err := s.repository.SysSetting().List(ctx, key, page, pageSize)
	if err != nil {
		return
	}
	for _, v := range list {
		res = append(res, listResult{
			Section:   v.Section,
			Key:       v.Key,
			Value:     v.Value,
			Id:        v.Id,
			Remark:    v.Description,
			Enable:    v.Enable,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return
}

func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "syssetting", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}
