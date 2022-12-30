package main

import (
	"go.backend/backend"
	"go.backend/conf"
	"go.backend/server/http"
	"go.backend/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	config := conf.InitConfigWithFile("./conf/config.yml")
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
