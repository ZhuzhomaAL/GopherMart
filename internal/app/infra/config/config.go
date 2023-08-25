package config

import (
	"flag"
	"os"
)

type Config struct {
	DatabaseUri          string
	RunAddress           string
	AccrualSystemAddress string
}

func MakeConfig() Config {
	var config Config

	flag.StringVar(&config.RunAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.AccrualSystemAddress, "r", "/", "accrual system address")
	flag.StringVar(
		&config.DatabaseUri, "d", "postgres://postgres:345973@localhost:5432/postgres?sslmode=disable", "database connection",
	)
	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		config.RunAddress = envRunAddress
	}

	if envDatabaseUri := os.Getenv("DATABASE_URI"); envDatabaseUri != "" {
		config.DatabaseUri = envDatabaseUri
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		config.AccrualSystemAddress = envAccrualSystemAddress
	}

	return config
}
