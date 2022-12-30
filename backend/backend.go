package backend

import (
	"go.backend/client"
	"go.backend/conf"
	"go.backend/dao"
	"go.backend/service"
)

type Backend struct {
	c   *conf.AppConfig
	srv *service.Service
	dao *dao.Dao
	mgr *client.Manager
}

func New(c *conf.AppConfig) *Backend {
	return &Backend{
		c:   c,
		srv: service.New(c),
		dao: dao.New(c),
		mgr: client.New(c),
	}
}

func (b *Backend) Start() {

}
