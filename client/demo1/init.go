package demo1

import (
	"hertz_demo/conf"
	"hertz_demo/initialization"
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
