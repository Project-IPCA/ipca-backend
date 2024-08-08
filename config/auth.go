package config

import (
	"os"
)

type AuthConfig struct {
	AccessSecretUserStudent  string
	RefreshSecretUserStudent string
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		AccessSecretUserStudent:  os.Getenv("ACCESS_SECRET_USER_STUDENT"),
		RefreshSecretUserStudent: os.Getenv("REFRESH_SECRET_USER_STUDENT"),
	}
}
