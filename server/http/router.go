package http

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func InitRouter(h *server.Hertz) {
	h.Any("/ping", ping)
	v1 := h.Group("/api")
	v1.GET("/test", Test)
	v1.GET("/test/mgr", TestMgr)

}
