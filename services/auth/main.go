package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	pb "distributed-media-processing-platform/proto/generated/proto/auth/v1"

	"github.com/jackc/pgx/v5/pgxpool"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

var (
	pool         *pgxpool.Pool
	db_ctx       = context.Background()
	AuthRepo     *AuthRepository
	port         = os.Getenv("AUTH_PORT")
	envJwtSecret = os.Getenv("ACCESS_KEY_JWT_SECRET")
	dsn          = os.Getenv("AUTH_DATABASE_URL")
)

func main() {
	if envJwtSecret == "" {
		log.Fatal("ACCESS_KEY_JWT_SECRET is not set")
	}
	if dsn == "" {
		log.Fatal("AUTH_DATABASE_URL is not set")
	}
	if port == "" {
		log.Fatal("AUTH_PORT is not set")
	}

	pool, err := pgxpool.New(db_ctx, dsn)
	if err != nil {
		log.Fatalf("AUTH DB err:%v", err)
	}
	err = pool.Ping(db_ctx)
	if err != nil {
		log.Fatalf("AUTH DB is unreachable: %v", err)
	}
	AuthRepo = NewAuthRepository(pool)

	// --GRPC SERVER START--
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})
	log.Printf("AUTH service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve AUTH: %v", err)
	}
}
