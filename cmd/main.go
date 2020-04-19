package main

import (
	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/storage/inmemory"
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

	stor := inmemory.New()
	//stor := disk.New(conf.CachePath)
	storDis := dispatcher.New(&stor, conf.CacheSize, &log)

	server := web.NewServer(web.Config{
		HttpListen: conf.HttpListen,
	}, &log, &storDis)

	server.RunServer()
	defer server.Shutdown()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	log.Infof("Got signal from OS: %v. Exit...", <-osSignals)

}
