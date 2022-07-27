package config

type RedisConfig struct {
	Port     string `envconfig:"REDIS_PORT"`
	Host     string `envconfig:"REDIS_HOST"`
	Password string `envconfig:"REDIS_PASSWORD"`
	User     string `envconfig:"REDIS_USER"`
}
