package initialization

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/registry/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type HTTP struct {
	ServiceName string
}

func InitHTTP(serviceName string) *HTTP {
	return &HTTP{ServiceName: serviceName}
}

func (h *HTTP) JsonPost(ctx context.Context, url string, data []byte) (int, []byte, error) {
	url = fmt.Sprintf("server://%s%s", h.ServiceName, url)
	req := protocol.AcquireRequest()
	req.SetOptions(config.WithSD(true))
	req.SetMethod(consts.MethodPost)
	req.SetRequestURI(url)
	req.SetBody(data)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	resp := protocol.AcquireResponse()
	return h.call(ctx, req, resp)

}

func (h *HTTP) FormPost(ctx context.Context, url string, data map[string]string) (int, []byte, error) {
	url = fmt.Sprintf("server://%s%s", h.ServiceName, url)
	req := protocol.AcquireRequest()
	req.SetOptions(config.WithSD(true))
	req.SetMethod(consts.MethodPost)
	req.SetRequestURI(url)
	req.Header.SetContentTypeBytes([]byte("application/x-www-form-urlencoded"))
	req.SetFormData(data)
	resp := protocol.AcquireResponse()
	return h.call(ctx, req, resp)
}

func (h *HTTP) Get(ctx context.Context, url string) (int, []byte, error) {
	url = fmt.Sprintf("server://%s%s", h.ServiceName, url)
	req := protocol.AcquireRequest()
	req.SetOptions(config.WithSD(true))
	req.SetMethod(consts.MethodGet)
	req.SetRequestURI(url)

	return h.call(ctx, req, nil)
}

func (h *HTTP) call(ctx context.Context, req *protocol.Request, resp *protocol.Response) (int, []byte, error) {
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
	}

	nacosCli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		})
	if err != nil {
		panic(err)
	}
	r := nacos.NewNacosResolver(nacosCli)
	cli.Use(sd.Discovery(r))
	// 客户端Tracing
	cli.Use(hertztracing.ClientMiddleware())
	err = cli.Do(ctx, req, resp)
	if err != nil {
		hlog.Fatal(err)
		return resp.StatusCode(), nil, err
	}
	return resp.StatusCode(), resp.Body(), nil
}
