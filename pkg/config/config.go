package config

import (
	"os"
	"scrapping/pkg/logger"
	"strconv"

	env "github.com/joho/godotenv"
)

// Loader load config from reader into osiper


type Config struct {
	// db
	
	Debug     bool
	ApiServer ApiServer
}

type ENV interface {
	GetBool(string) bool
	GetEnv(string) string
}

type ApiServer struct {
	Port           string
	AllowedOrigins string
}

func LoadConfig() *Config {
	err := env.Load(".env")
	if err != nil {
		logger.L.Error(err, "Error loading .enos file")
	}
	envValue := os.Getenv("DEBUG")
	debug, err := strconv.ParseBool(envValue)
	if err != nil {
		logger.L.Error(err, "Error parsing boolean value from environment variable")
	}
	return &Config{
		Debug: debug,

		ApiServer: ApiServer{
			Port:           os.Getenv("PORT"),
			AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
		},
	}
}
