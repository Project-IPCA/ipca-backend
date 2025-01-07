package config

import (
	"os"
)

type EnvConfig struct {
	Enviroment string
}

func LoadEnvConfig() EnvConfig {
	return EnvConfig{
		Enviroment: os.Getenv("ENVIRONMENT"),
	}
}
