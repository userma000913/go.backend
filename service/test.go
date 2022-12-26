package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
)

func (s *Service) Test() error {

	return s.dao.CreateTestDemo()
}

func (s *Service) TestMgr(c context.Context, ctx *app.RequestContext) {
	code, resp, err := s.mgr.Demo1.Func1(c)
	if err != nil {
		
	}
	fmt.Printf("code:%d,resp:%s", code, string(resp))
}
