package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryInterceptot() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if info.FullMethod == "/auth.AuthService/Login" {
			return handler(ctx, req)
		}
		if err := authenticate(ctx); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func StreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if err := authenticate(ss.Context()); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metaDat missing")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token missing")
	}

	authHeader := values[0]

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return status.Error(codes.Unauthenticated, "invalid authorization format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if _, err := ValidateJwt(tokenString); err != nil {
		return status.Error(codes.Unauthenticated, "invalid or Expired token")
	}
	return nil
}
