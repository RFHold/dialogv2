package main

import (
	"dialogv2/internal/config"
	"dialogv2/internal/database"
	"dialogv2/internal/services/message"
	"dialogv2/pb/messages"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"log"
	"net"
)

type environmentConfig struct {
	Environment   string `envconfig:"ENVIRONMENT"`
	Debug         bool   `envconfig:"DEBUG"`
	DBConfig      config.DBConfig
	RedisConfig   config.RedisConfig
	ServiceConfig config.ServiceConfig
}

func getConfig() (environmentConfig, error) {
	_ = godotenv.Overload(".env", ".env.local") // For .env file in local development

	var cfg environmentConfig
	err := envconfig.Process("", &cfg)

	if err != nil {
		return environmentConfig{}, err
	}

	return cfg, nil
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+cfg.ServiceConfig.MessagesServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	redisClient := message.NewRedisClient(&cfg.RedisConfig)

	log.Println("Connected to Redis")

	db, err := database.Connect(&cfg.DBConfig)
	if err != nil {
		log.Fatalf("failed to connect to Database: %v", err)
	}

	log.Println("Connected to Database")

	repository := message.Repository{
		DB:          db,
		RedisClient: redisClient,
	}

	grpcServer := grpc.NewServer(opts...)

	messages.RegisterMessageServiceServer(grpcServer, &message.ServiceServer{
		Repository:  &repository,
		RedisClient: redisClient,
	})

	log.Println("Service registered")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
