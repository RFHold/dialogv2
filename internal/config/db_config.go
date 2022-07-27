package config

import "fmt"

type DBConfig struct {
	Host     string `envconfig:"POSTGRES_HOST"`
	Port     string `envconfig:"POSTGRES_PORT" default:"5432"`
	Database string `envconfig:"POSTGRES_DATABASE"`
	Username string `envconfig:"POSTGRES_USERNAME"`
	Password string `envconfig:"POSTGRES_PASSWORD"`
	SSL      string `envconfig:"POSTGRES_SSL"`
}

func (c DBConfig) DBString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=%s", c.Username, c.Password, c.Database, c.Port, c.Host, c.SSL)
}
