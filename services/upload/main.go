package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	pb "distributed-media-processing-platform/proto/generated/proto/auth/v1"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

func main() {
	port := os.Getenv("UPLOAD_PORT")
	// --GRPC SERVER START--
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})
	log.Printf("UPLOAD service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve AUTH: %v", err)
	}
}
