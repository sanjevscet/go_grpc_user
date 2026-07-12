package auth

import (
	"context"
	"errors"
	authpb "go-grpc-user/proto/auth"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "usrename & password are required")
	}

	var id int32
	var usrename string
	var role string

	err := s.db.QueryRow(
		ctx,
		`
			SELECT id, username, role
			FROM users
			WHERE  username = $1 and password = $2	
		`,
		req.Username,
		req.Password,
	).Scan(&id, &usrename, &role)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.Unauthenticated, "invalid user or password")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login faield %v", err)
	}

	token, err := GenerateJWT(id, usrename, role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate Token %v", err)
	}

	return &authpb.LoginResponse{
		AccessToken: token,
	}, nil
}
