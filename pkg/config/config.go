package config

import (
	"os"
)

// Config contains the environment specific configuration values needed by the
// application.
type Config struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	Environment      string
	Port             int
	SentryDSN        string
	StatsdHost       string
	StatsdPort       int
}

const environmentENV = "ENVIRONMENT"

// New returns an instance of Config based on the "ENVIRONMENT" environment
// variable.
func New() Config {
	cfg := Config{
		DatabasePort: 5432,
		Port:         3000,
		SentryDSN:    os.Getenv("SENTRY_DSN"),
		StatsdHost:   "127.0.0.1",
		StatsdPort:   8125,
	}

	switch os.Getenv(environmentENV) {
	case "development", "":
		loadDevelopmentConfig(&cfg)
	case "test":
		loadTestConfig(&cfg)
	}

	return cfg
}
