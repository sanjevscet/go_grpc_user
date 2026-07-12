package config

import "os"

type Config struct {
	DatabaseUrl string
	GRPCPort    string
}

func Load() Config {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "postgres://postgres:password@localhost:14432/grpc_db?sslmode=disable"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50501"
	}
	return Config{
		DatabaseUrl: databaseUrl,
		GRPCPort:    grpcPort,
	}
}
