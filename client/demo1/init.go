package demo1

import (
	"go.backend/conf"
	"go.backend/initialization"
)

type Manager struct {
	c          *conf.AppConfig
	httpClient *initialization.HTTP
	Url1       string
}

func NewManager(c *conf.AppConfig) *Manager {
	return &Manager{
		c:          c,
		httpClient: initialization.InitHTTP("demo1"),
		Url1:       "/ping",
	}
}
