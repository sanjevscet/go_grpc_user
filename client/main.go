package main

import (
	"context"
	userpb "go-grpc-user/proto"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	client := userpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	getUser(client, ctx, 1)
	// getUsers(client, ctx)
	createUsers(client, ctx)

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
		{Name: "Tannu", Email: "tannu@example.com"},
		{Name: "Deepu", Email: "deepu@example.com"},
		{Name: "Guddu", Email: "gudu@example.com"},
		{Name: "ac", Email: "gudu1@example.com"},
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
