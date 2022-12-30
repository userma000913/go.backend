package service

import (
	"go.backend/client"
	"go.backend/conf"
	"go.backend/dao"
)

type Service struct {
	c *conf.AppConfig

	dao *dao.Dao
	mgr *client.Manager
}

func New(c *conf.AppConfig) *Service {
	return &Service{
		c:   c,
		dao: dao.New(c),
		mgr: client.New(c),
	}
}

func (s *Service) Close() {
	if s.dao != nil {
		s.dao.Close()
	}

	if s.mgr != nil {
		s.mgr.Close()
	}
}
