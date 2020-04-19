package config

import (
	"github.com/caarlos0/env"
	"github.com/dark705/otus_previewer/internal/helpers"
)

type Config struct {
	LogLevel   string `env:"LOG_LEVEL" envDefault:"debug"`
	HttpListen string `env:"HTTP_LISTEN" envDefault:":8013"`
	CacheSize  int    `env:"CACHE_SIZE" envDefault:"7500"`
	CacheType  string `env:"CACHE_TYPE" envDefault:"disk"`
	CachePath  string `env:"CACHE_TYPE" envDefault:"./cache"`
}

func GetConfigFromEnv() Config {
	c := Config{}

	err := env.Parse(&c)
	helpers.FailOnError(err, "Fail get config from Env")

	return c
}
