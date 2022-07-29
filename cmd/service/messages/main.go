package main

import (
	"context"
	"dialogv2/internal/config"
	"dialogv2/internal/database"
	"dialogv2/internal/services/message"
	"dialogv2/pb/messages"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
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

	go runGateway()
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runGateway() {
	log.Println("Setup gateway")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := messages.RegisterMessageServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to Register: %v", err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
