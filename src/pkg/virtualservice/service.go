package virtualservice

import "github.com/go-kit/kit/log"

type Service interface {
}

type service struct {
	logger log.Logger
}

func NewService() Service {
	return &service{}
}
