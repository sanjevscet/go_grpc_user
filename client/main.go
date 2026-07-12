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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("=========== Get Single User Unary Call=============")

	userId := int32(1)
	user, err := client.GetUser(ctx, &userpb.GetUserRequest{
		Id: userId,
	})
	if err != nil {
		log.Fatalf("Failed to get user with Id %d: err: %v", userId, err)
	}
	log.Println("Getting Single User")
	log.Printf("User Details %v", user)

	log.Println("=========== Get All User via Stream =============")

	stream, err := client.GetUsers(ctx, &userpb.GetUsersRequest{})
	if err != nil {
		log.Fatalf("Get Users failed %v", err)
	}

	log.Println("All Users:")

	for {
		user, err := stream.Recv()
		if err == io.EOF {
			log.Println("Fetching all users completed")
			break
		}
		if err != nil {
			log.Fatalf("Error received %v", err)
		}
		log.Printf("Id = %d, Name= %s, Email= %s", user.Id, user.Name, user.Email)
	}

}
