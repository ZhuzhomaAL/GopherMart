package config

import (
	"flag"
	"os"
)

type Config struct {
	DatabaseURI          string
	RunAddress           string
	AccrualSystemAddress string
	LogLevel             string
} //

func MakeConfig() Config {
	var config Config

	flag.StringVar(&config.RunAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.AccrualSystemAddress, "r", "/", "accrual system address")
	flag.StringVar(
		&config.DatabaseURI, "d", "", "database connection",
	)
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		config.RunAddress = envRunAddress
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		config.DatabaseURI = envDatabaseURI
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		config.AccrualSystemAddress = envAccrualSystemAddress
	}

	return config
}
