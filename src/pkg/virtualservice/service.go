package virtualservice

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Service interface {
	List(ctx context.Context)
}

type service struct {
	logger log.Logger
}

func NewService(logger log.Logger) Service {
	return &service{}
}

func (c *service) List(ctx context.Context) {

}
