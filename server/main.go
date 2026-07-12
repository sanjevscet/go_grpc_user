package main

import (
	"context"
	"fmt"
	userpb "go-grpc-user/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type userServer struct {
	userpb.UnimplementedUserServiceServer
}

var users = []*userpb.User{
	{
		Id:    1,
		Name:  "Rahul",
		Email: "rahul@example.com",
	},
	{
		Id:    2,
		Name:  "Aman",
		Email: "aman@example.com",
	},
	{
		Id:    3,
		Name:  "Deepak",
		Email: "deepak@example.com",
	},
	{
		Id:    4,
		Name:  "Sanjeev",
		Email: "sanjeev@example.com",
	},
	{
		Id:    5,
		Name:  "Rohit",
		Email: "rohit@example.com",
	},
}

func (s *userServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	for _, user := range users {
		if user.Id == req.Id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with id %d not found", req.Id)
}

func (s *userServer) GetUsers(
	req *userpb.GetUsersRequest,
	stream grpc.ServerStreamingServer[userpb.User],
) error {
	for _, user := range users {
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":50501")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, &userServer{})

	log.Println("grcp server is running on port tcp:50501")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}
