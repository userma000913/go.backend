package http

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func InitRouter(h *server.Hertz) {
	v1 := h.Group("/api")
	v1.GET("/test", Test)
	v1.GET("/test/mgr", TestMgr)
	v1.Any("/ping", ping)
}
