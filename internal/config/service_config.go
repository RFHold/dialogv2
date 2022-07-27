package config

type ServiceConfig struct {
	MessagesServicePort string `envconfig:"MESSAGES_SERVICE_PORT" default:"50051"`
}
