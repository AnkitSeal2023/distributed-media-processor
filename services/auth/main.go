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

	env "github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

var (
	port     = flag.Int("port", 50051, "The server port")
	pool     *pgxpool.Pool
	ctx      = context.Background()
	AuthRepo *AuthRepository
)

func main() {
	err := env.Load()
	if err != nil {
		log.Fatalf("error occured during loading env variables: %v", err)
	}

	pool, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("AUTH DB err:%v", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("AUTH DB is unreachable: %v", err)
	}
	AuthRepo = NewAuthRepository(pool)

	// --GRPC SERVER START--
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
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
