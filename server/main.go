package main

import (
	"context"
	authpb "go-grpc-user/proto/auth"
	userpb "go-grpc-user/proto/user"
	"go-grpc-user/server/auth"
	"go-grpc-user/server/config"
	"go-grpc-user/server/db"
	"go-grpc-user/server/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	ctx := context.Background()
	config := config.Load()

	dbPool, err := db.NewPostgresPool(ctx, config.DatabaseUrl)
	if err != nil {
		log.Fatalf("failed to connect to database %v", err)
	}
	defer dbPool.Close()

	listener, err := net.Listen("tcp", ":"+config.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryInterceptot()),
		grpc.StreamInterceptor(auth.StreamInterceptor()),
	)

	authService := auth.NewAuthService(dbPool)
	authpb.RegisterAuthServiceServer(grpcServer, authService)

	userService := service.NewUserService(dbPool)
	userpb.RegisterUserServiceServer(grpcServer, userService)

	log.Printf("grcp server is running on port %s", config.GRPCPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}
