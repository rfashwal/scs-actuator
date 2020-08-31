package internal

import (
	"github.com/rfashwal/scs-utilities/config"
	"github.com/rfashwal/scs-utilities/rabbit/publishing"
)

type Service interface {
	Ping() error
}

func NewService(conf config.Manager) (Service, error) {
	return service{
		conf: conf}, nil
}

type service struct {
	publisher *publishing.Publisher
	conf      config.Manager
}

func (s service) Ping() error {
	return nil
}
