package main

import (
	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/storage"
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
	stor := storage.CreateStorage(conf.CacheType, conf.CachePath, &log)
	storDis := dispatcher.New(stor, conf.CacheSize, &log)

	server := web.NewServer(web.Config{
		HttpListen:       conf.HttpListen,
		ImageMaxFileSize: conf.ImageMaxFileSize,
	}, &log, &storDis)

	server.RunServer()
	defer server.Shutdown()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	log.Infof("Got signal from OS: %v. Exit...", <-osSignals)

}
