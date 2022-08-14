package main

import (
	"dialogv2/internal/config"
	"dialogv2/internal/database"
	"dialogv2/internal/services/users"
	pb "dialogv2/pb/users"
	"google.golang.org/grpc"
	"log"
	"net"
)

type environmentConfig struct {
	Environment   string `envconfig:"ENVIRONMENT"`
	Debug         bool   `envconfig:"DEBUG"`
	DBConfig      config.DBConfig
	ServiceConfig config.ServiceConfig
}

func main() {
	var cfg environmentConfig
	err := config.LoadConfig(cfg)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+cfg.ServiceConfig.UsersServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	db, err := database.Connect(&cfg.DBConfig)
	if err != nil {
		log.Fatalf("failed to connect to Database: %v", err)
	}

	log.Println("Connected to Database")

	repository := users.Repository{
		DB: db,
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterServiceServer(grpcServer, &users.ServiceServer{
		Repository: &repository,
	})

	log.Println("Service registered")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
