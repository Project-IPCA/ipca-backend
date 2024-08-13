package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
	Auth AuthConfig
	Redis RedisConfig
	RabbitMQ RabbitMQConfig
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		Auth: LoadAuthConfig(),
		HTTP: LoadHTTPConfig(),
		DB:   LoadDBConfig(),
		Redis : LoadRedisConfig(),
		RabbitMQ: LoadRabbitMQConfig(),
	}
}
