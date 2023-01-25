package config

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

var EnvConfig *envConfig

// Struct with environment variables used in app
type envConfig struct {
	ServerPort        string
	RabbitUrl         string
	DbHost            string
	DbPort            string
	DbUser            string
	DbPass            string
	DbName            string
	SslMode           string
	NotificationQueue string
}

// Load the env variables from .env file. Defined variables
// in docker-compose file will override those from .env file
func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	EnvConfig = &envConfig{
		ServerPort:        os.Getenv("SERVER_PORT"),
		RabbitUrl:         os.Getenv("RABBIT_URL"),
		DbHost:            os.Getenv("DB_HOST"),
		DbPort:            os.Getenv("DB_PORT"),
		DbUser:            os.Getenv("DB_USER"),
		DbPass:            os.Getenv("DB_PASS"),
		DbName:            os.Getenv("DB_NAME"),
		SslMode:           os.Getenv("SSL_MODE"),
		NotificationQueue: os.Getenv("NOTIFICATION_QUEUE"),
	}
}
