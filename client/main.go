package main

import (
	"context"
	authpb "go-grpc-user/proto/auth"
	userpb "go-grpc-user/proto/user"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.NewClient(
		"passthrough:///localhost:50501",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	userClient := userpb.NewUserServiceClient(conn)
	authClient := authpb.NewAuthServiceClient(conn)

	loginRes, err := authClient.Login(context.Background(), &authpb.LoginRequest{
		Username: "sanjeev",
		Password: "sanjeev123",
	})

	if err != nil {
		log.Fatalf("login failed %v", err)
	}

	token := loginRes.AccessToken
	log.Printf("Login successful %s\n", token)
	md := metadata.Pairs("authorization", "Bearer "+token)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, md)
	getUser(userClient, ctx, 1)
	getUsers(userClient, ctx)
	createUsers(userClient, ctx)
}

func getUser(client userpb.UserServiceClient, ctx context.Context, id int32) {
	user, err := client.GetUser(ctx, &userpb.GetUserRequest{Id: id})
	if err != nil {
		log.Fatalf("Get User failed %v", err)
	}

	log.Printf("Single User: ID = %d, Name = %s, Email = %s", user.Id, user.Name, user.Email)
}

func getUsers(client userpb.UserServiceClient, ctx context.Context) {
	stream, err := client.GetUsers(ctx, &userpb.GetUsersRequest{})
	if err != nil {
		log.Fatalf("Get Users Failed %v", err)
	}
	log.Println("Users from stream server")

	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream receive failed %v", err)
		}

		log.Printf("ID=%d, Name =%s, Email =%s", user.Id, user.Name, user.Email)
	}
}

func createUsers(client userpb.UserServiceClient, ctx context.Context) {
	stream, err := client.CreateUsers(ctx)
	if err != nil {
		log.Fatalf("CreateUsers failed %v", err)
	}

	users := []*userpb.CreateUserRequest{
		{Name: "Tannu", Email: "tannu@example.com", Password: "tanu123", Role: "user", Username: "tanu"},
		{Name: "Deepu", Email: "deepu@example.com", Password: "deepu123", Role: "user", Username: "deepu"},
		{Name: "Guddu", Email: "gudu@example.com", Password: "guddu123", Role: "user", Username: "guddu"},
		{Name: "John", Email: "john@example.com", Password: "john123", Role: "user", Username: "john"},
	}

	for _, user := range users {
		if err := stream.Send(user); err != nil {
			log.Fatalf("failed to send users %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to reveceive user response %v", err)
	}

	log.Printf("Created count %d", res.CreatedCount)
	log.Printf("Message %s", res.Message)
}
