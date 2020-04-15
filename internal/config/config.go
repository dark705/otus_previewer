package config

import (
	"fmt"
	"os"
)

type Config struct {
	LogLevel   string
	HttpListen string
	CacheSize  int
}

func GetConfigFromEnv() Config {
	c := Config{
		LogLevel:   "debug",
		HttpListen: ":8013",
		CacheSize:  10 * 1024 * 1024,
	}

	val := os.Getenv("LOG_LEVEL")
	if val == "" {
		fmt.Println("env LOG_LEVEL not defined, use default value:", c.LogLevel)
	} else {
		c.LogLevel = val
	}

	val = os.Getenv("HTTP_LISTEN")
	if val == "" {
		fmt.Println("env HTTP_LISTEN not defined, use default value:", c.HttpListen)
	} else {
		c.HttpListen = val
	}

	return c
}
