package config

import (
	"os"
	"strconv"
)

const (
	DefaultPort      = 3000
	DefaultTMDBToken = ""
)

type Config struct {
	Port      int
	TMDBToken string
}

func LoadConfigFromEnv() *Config {
	config := Config{
		Port:      DefaultPort,
		TMDBToken: DefaultTMDBToken,
	}

	if val := os.Getenv("PORT"); val != "" {
		port, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		config.Port = port
	}

	if val := os.Getenv("TMDB_TOKEN"); val != "" {
		config.TMDBToken = val
	}

	return &config
}
