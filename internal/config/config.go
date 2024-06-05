package config

import (
	"os"
	"strconv"
)

const DefaultPort = 3000

type Config struct {
	Port int
}

func LoadConfigFromEnv() *Config {
	config := Config{
		Port: DefaultPort,
	}

	if val := os.Getenv("PORT"); val != "" {
		port, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		config.Port = port
	}

	return &config
}
