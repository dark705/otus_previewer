package config

import (
	"github.com/caarlos0/env"
	"github.com/dark705/otus_previewer/internal/helpers"
)

type Config struct {
	LogLevel         string `env:"LOG_LEVEL" envDefault:"debug"`
	HttpListen       string `env:"HTTP_LISTEN" envDefault:":8013"`
	ImageMaxFileSize int    `env:"IMAGE_MAX_FILE_SIZE" envDefault:"1000000"`
	CacheSize        int    `env:"CACHE_SIZE" envDefault:"100000000"`
	CacheType        string `env:"CACHE_TYPE" envDefault:"inmemory"`
	CachePath        string `env:"CACHE_TYPE" envDefault:"./cache"`
}

func GetConfigFromEnv() Config {
	c := Config{}

	err := env.Parse(&c)
	helpers.FailOnError(err, "Fail get config from Env")

	return c
}
