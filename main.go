package main

import (
	"context"
	"database-example/handlers"
	stakeholderspb "database-example/proto/stakeholders"
	"database-example/repo"
	"database-example/service"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := log.New(os.Stdout, "[stakeholders] ", log.LstdFlags)

	userRepo, err := repo.NewUserRepository(logger)
	if err != nil {
		logger.Fatal("Failed to create UserRepository:", err)
	}
	defer userRepo.DriverClose(context.Background())

	userService := &service.UserService{UserRepo: userRepo}
	userHandler := handlers.NewUserHandler(userService)

	addr := os.Getenv("STAKEHOLDERS_SERVICE_ADDRESS")
	if addr == "" {
		addr = ":8000"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	stakeholderspb.RegisterStakeholdersServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	go func() {
		logger.Println("Starting gRPC server on", addr)
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("gRPC server error:", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh

	logger.Println("Shutting down gRPC server...")
	grpcServer.Stop()
}
