/**
 * @Time : 3/5/21 2:41 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package namespace

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"strings"
)

type Middleware func(next Service) Service

type Call func() error

type Service interface {
	// Save 保存或更新
	Save(ctx context.Context, data *types.Namespace) (err error)
	SaveCall(ctx context.Context, data *types.Namespace, call Call) (err error)
	FindByIds(ctx context.Context, ids []int64) (res []types.Namespace, err error)
	FindByName(ctx context.Context, clusterId int64, name string) (res types.Namespace, err error)
	List(ctx context.Context, clusterId int64, names []string, query string, page, pageSize int) (res []types.Namespace, total int, err error)
	FindByNames(ctx context.Context, clusterId int64, names []string) (res []types.Namespace, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindByNames(ctx context.Context, clusterId int64, names []string) (res []types.Namespace, err error) {
	err = s.db.Model(&types.Namespace{}).
		Where("cluster_id = ? AND name IN (?)", clusterId, names).
		Find(&res).Error
	return
}

func (s *service) List(ctx context.Context, clusterId int64, names []string, query string, page, pageSize int) (res []types.Namespace, total int, err error) {
	q := s.db.Model(&types.Namespace{}).Where("cluster_id = ? AND name IN(?)", clusterId, names)
	if !strings.EqualFold(query, "") {
		q = q.Where("name LIKE ?", "'%"+query+"%'")
	}
	err = q.Count(&total).Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&res).Error
	return
}

func (s *service) SaveCall(ctx context.Context, data *types.Namespace, call Call) (err error) {
	tx := s.db.Model(data).Begin()
	if err = tx.Save(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = call(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *service) FindByName(ctx context.Context, clusterId int64, name string) (res types.Namespace, err error) {
	err = s.db.Model(&types.Namespace{}).Where("cluster_id = ? AND name = ?", clusterId, name).First(&res).Error
	return
}

func (s *service) FindByIds(ctx context.Context, ids []int64) (res []types.Namespace, err error) {
	err = s.db.Model(&types.Namespace{}).Where("id IN (?)", ids).Find(&res).Error
	return
}

func (s *service) Save(ctx context.Context, data *types.Namespace) (err error) {
	return s.db.Model(data).Save(data).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
