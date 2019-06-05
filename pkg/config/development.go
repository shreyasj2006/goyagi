package config

import (
	"os"
	"strconv"
)

func loadDevelopmentConfig(cfg *Config) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		cfg.Port = port
	}
	cfg.Environment = "development"
}
