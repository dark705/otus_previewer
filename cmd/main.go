package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dark705/otus_previewer/internal/config"
	"github.com/dark705/otus_previewer/internal/logger"
	"github.com/dark705/otus_previewer/internal/web"
)

func main() {
	conf := config.GetConfigFromEnv()

	log := logger.NewLogger(logger.Config{
		Level: conf.LogLevel,
	})

	server := web.NewServer(web.Config{
		HttpListen: conf.HttpListen,
	}, &log)

	server.RunServer()
	defer server.Shutdown()

	fmt.Println(conf)

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	log.Infof("Got signal from OS: %v. Exit...", <-osSignals)

}
