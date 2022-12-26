package http

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	hertzSentinel "github.com/hertz-contrib/opensergo/sentinel/adapter"
	"github.com/hertz-contrib/registry/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"hertz_demo/conf"
	"hertz_demo/service"
	"log"
)

var (
	svc *service.Service
	h   *server.Hertz
)

func Init(s *service.Service, config *conf.AppConfig) {
	svc = s
	addr := fmt.Sprintf("127.0.0.1:%d", config.Server.Port)
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public",
		Username:            "nacos",
		Password:            "nacos",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// nacos注册中心客户端
	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 服务注册
	r := nacos.NewNacosRegistry(cli)
	tracer, cfg := hertztracing.NewServerTracer()
	h = server.Default(
		server.WithHostPorts(addr),
		server.WithRegistry(r, &registry.Info{
			ServiceName: config.Server.Name,
			Addr:        utils.NewNetAddr("tcp", addr),
			Weight:      10,
			Tags:        nil,
		}),
		tracer,
	)
	// Tracing & Sentinel
	// todo cors
	h.Use(hertztracing.ServerMiddleware(cfg), hertzSentinel.SentinelServerMiddleware())

	// register handler with http route
	InitRouter(h)

	// start a http server
	go func() {
		h.Spin()
	}()
}
func Shutdown() {

	if svc != nil {
		svc.Close()
	}
}
