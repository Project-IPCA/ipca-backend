package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      EnvConfig
	HTTP     HTTPConfig
	DB       DBConfig
	Auth     AuthConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	Minio    MinioConfig
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		Env:      LoadEnvConfig(),
		Auth:     LoadAuthConfig(),
		HTTP:     LoadHTTPConfig(),
		DB:       LoadDBConfig(),
		Redis:    LoadRedisConfig(),
		RabbitMQ: LoadRabbitMQConfig(),
		Minio:    LoadMinioConfig(),
	}
}
