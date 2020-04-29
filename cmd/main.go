package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/http"
	"github.com/dark705/otus_previewer/internal/storage"

	"github.com/dark705/otus_previewer/internal/config"
	"github.com/dark705/otus_previewer/internal/logger"
)

func main() {
	conf := config.GetConfigFromEnv()
	fmt.Printf("current settings: %+v\n", conf)
	log := logger.NewLogger(logger.Config{
		Level: conf.LogLevel,
	})
	st := storage.Create(conf.CacheType, conf.CachePath, log)
	imageDispatcher := dispatcher.New(st, conf.CacheSize, log)

	server := http.NewServer(http.Config{
		HTTPListen:       conf.HTTPListen,
		ImageMaxFileSize: conf.ImageMaxFileSize,
		ImageGetTimeout:  conf.ImageGetTimeout,
	}, log, &imageDispatcher)

	server.RunServer()
	defer server.Shutdown()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	log.Infof("got signal from OS: %v. Exit...", <-osSignals)
}
