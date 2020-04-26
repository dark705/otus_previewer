package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/storage"

	"github.com/dark705/otus_previewer/internal/config"
	"github.com/dark705/otus_previewer/internal/logger"
	"github.com/dark705/otus_previewer/internal/web"
)

func main() {
	conf := config.GetConfigFromEnv()
	fmt.Printf("Current settings: %+v\n", conf)
	log := logger.NewLogger(logger.Config{
		Level: conf.LogLevel,
	})
	st := storage.Create(conf.CacheType, conf.CachePath, &log)
	imageDispatcher := dispatcher.New(st, conf.CacheSize, &log)

	server := web.NewServer(web.Config{
		HttpListen:       conf.HttpListen,
		ImageMaxFileSize: conf.ImageMaxFileSize,
	}, &log, &imageDispatcher)

	server.RunServer()
	defer server.Shutdown()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	log.Infof("Got signal from OS: %v. Exit...", <-osSignals)
}
