package config

type ServiceConfig struct {
	MessagesServicePort string `envconfig:"MESSAGES_SERVICE_PORT" default:"50051"`
	UsersServicePort    string `envconfig:"USERS_SERVICE_PORT" default:"50052"`
}
