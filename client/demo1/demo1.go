package demo1

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/json"
)

func (m *Manager) Func1(c context.Context) (int, []byte, error) {
	type Test struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	t := Test{A: 11, B: 22}
	b, err := json.Marshal(t)
	if err != nil {

	}
	return m.httpClient.JsonPost(c, m.Url1, b)

}
