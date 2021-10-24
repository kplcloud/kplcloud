/**
 * @Time : 8/11/21 11:25 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package nodes

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"strings"
)

type Middleware func(Service) Service

type Callback func() error

type Service interface {
	Save(ctx context.Context, data *types.Nodes) (err error)
	FindByName(ctx context.Context, clusterId int64, name string) (res types.Nodes, err error)
	List(ctx context.Context, clusterId int64, query string, page, pageSize int) (res []types.Nodes, total int, err error)
	Delete(ctx context.Context, clusterId int64, nodeName string, callback Callback) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Delete(ctx context.Context, clusterId int64, nodeName string, callback Callback) (err error) {
	tx := s.db.Model(&types.Nodes{}).Begin()
	if err = tx.Where("cluster_id = ? AND name = ?", clusterId, nodeName).
		Delete(&types.Nodes{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = callback(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *service) List(ctx context.Context, clusterId int64, query string, page, pageSize int) (res []types.Nodes, total int, err error) {
	q := s.db.Model(&types.Nodes{}).Where("cluster_id = ?", clusterId)
	if !strings.EqualFold(query, "") {
		q = q.Where("name LIKE ?", "%"+query+"%")
	}
	err = q.Count(&total).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&res).Error
	return
}

func (s *service) FindByName(ctx context.Context, clusterId int64, name string) (res types.Nodes, err error) {
	err = s.db.Model(&types.Nodes{}).Where("cluster_id = ?", clusterId).Where("name = ?", name).First(&res).Error
	return
}

func (s *service) Save(ctx context.Context, data *types.Nodes) (err error) {
	return s.db.Model(data).Save(data).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
