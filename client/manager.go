package client

import (
	"go.backend/client/demo1"
	"go.backend/conf"
)

// 服务发现
var mgr *Manager

type Manager struct {
	c     *conf.AppConfig
	Demo1 *demo1.Manager
}

func New(c *conf.AppConfig) *Manager {
	if mgr == nil {
		mgr = &Manager{
			c:     c,
			Demo1: demo1.NewManager(c),
		}
	}
	return mgr
}
func (m *Manager) Close() error {
	return nil
}
