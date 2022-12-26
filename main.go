package main

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
	"hertz_demo/backend"
	"hertz_demo/conf"
	"hertz_demo/server/http"
	"hertz_demo/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	f, err := os.OpenFile("./output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	hlog.SetLogger(hertzlogrus.NewLogger())
	hlog.SetOutput(f)
	hlog.SetLevel(hlog.LevelDebug)

	config := conf.InitConfigWithNacos()
	if config == nil {
		panic("config is nil")
	}
	s := service.New(config)
	http.Init(s, config)
	defer http.Shutdown()

	// 启动后台进程
	backend.New(config).Start()

	// 优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-sigChan
		log.Printf("get a signal %s\n", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Println("room.event.adapter server exit now...")
			return
		case syscall.SIGHUP:
		default:
		}
	}

}
