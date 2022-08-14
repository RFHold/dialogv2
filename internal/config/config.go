package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func LoadConfig[E interface{}](cfg E) (err error) {
	_ = godotenv.Overload(".env", ".env.local") // For .env file in local development

	err = envconfig.Process("", &cfg)

	return
}
