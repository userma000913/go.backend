package http

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/cors"
	hertzlogrus "github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/registry/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.backend/conf"
	"go.backend/middleware"
	"go.backend/service"
	"os"
)

var (
	svc *service.Service
	h   *server.Hertz
)

func Init(s *service.Service, config *conf.AppConfig) {

	// init log
	hlog.SetLogger(hertzlogrus.NewLogger())
	f, err := os.OpenFile("./output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	hlog.SetOutput(f)
	hlog.SetLevel(hlog.LevelTrace)

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
		hlog.Fatal(err)
		return
	}
	// 服务注册
	r := nacos.NewNacosRegistry(cli)

	// 链路追踪
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.Server.Name),
		provider.WithExportEndpoint(addr),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
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

	// Tracing & cors
	h.Use(hertztracing.ServerMiddleware(cfg), cors.Default(), middleware.AccessLog())

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
