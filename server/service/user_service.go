package service

import (
	"context"
	"errors"
	userpb "go-grpc-user/proto"
	"io"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetUser(
	ctx context.Context,
	req *userpb.GetUserRequest,
) (*userpb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var user userpb.User

	err := s.db.QueryRow(
		ctx,
		`
			SELECT id, name, email
			FROM users
			WHERE id = $1
		`,
		req.Id,
	).Scan(&user.Id, &user.Name, &user.Email)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.Id)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users %v", err)
	}

	return &user, nil
}

func (s *UserService) GetUsers(
	req *userpb.GetUsersRequest,
	stream grpc.ServerStreamingServer[userpb.User],
) error {
	ctx, cancel := context.WithTimeout(stream.Context(), 2*time.Minute)
	defer cancel()

	rows, err := s.db.Query(
		ctx,
		`
			SELECT id, name, email
			FROM users
			ORDER BY id
		`,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return status.Errorf(codes.DeadlineExceeded, "get users is timeout")
		}

		return status.Errorf(codes.Internal, "failed to get users %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user userpb.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			return status.Errorf(codes.Internal, "failed to scan users %v", err)
		}
		if err := stream.Send(&user); err != nil {
			return status.Errorf(codes.Internal, "failed to stream users %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return status.Errorf(codes.Internal, "rows error %v", err)
	}

	return nil
}

func (s *UserService) CreateUsers(
	stream grpc.ClientStreamingServer[userpb.CreateUserRequest, userpb.CreateUserResponse],
) error {
	ctx := stream.Context()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to begin transaction %v", err)
	}
	defer tx.Rollback(ctx)

	createdCount := int32(0)

	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			if err := tx.Commit(ctx); err != nil {
				return status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
			}
			return stream.SendAndClose(&userpb.CreateUserResponse{
				CreatedCount: createdCount,
				Message:      "users created successfully",
			})
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive users %v", err)
		}
		if req.Name == "" || req.Email == "" {
			return status.Error(codes.Internal, "named & email are required")
		}

		var insertedId int32

		err = tx.QueryRow(
			ctx,
			`
				INSERT INTO users (name, email)
				VALUES ($1, $2)
				ON CONFLICT (email) DO NOTHING
				RETURNING ID
			`,
			req.Name,
			req.Email,
		).Scan(&insertedId)

		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to insert users: %v", err)
		}
		createdCount++
	}
}
