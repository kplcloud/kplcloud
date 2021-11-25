/**
 * @Time : 2021/11/25 4:27 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package application

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(service Service) Service

type Call func() error

type Service interface {
	Save(ctx context.Context, app *types.Application, call Call) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Save(ctx context.Context, app *types.Application, call Call) (err error) {
	return s.db.Model(app).Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(app).Error; err != nil {
			return err
		}
		if call != nil {
			err = call()
		}
		return err
	})
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
