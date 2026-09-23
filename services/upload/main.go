package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	pb "distributed-media-processing-platform/proto/generated/proto/upload/v1"
	"distributed-media-processing-platform/services/upload/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUploadServiceServer
}

var (
	pool                  *pgxpool.Pool
	UploadRepo            *repository.UploadRepository
	db_ctx                = context.Background()
	dsn                   = os.Getenv("UPLOAD_DATABASE_URL")
	port                  = os.Getenv("UPLOAD_PORT")
	bucketName            = os.Getenv("BUCKET_NAME")
	minio_accessKeyID     = os.Getenv("MINIO_ACCESSKEYID")
	minio_secretAccessKey = os.Getenv("MINIO_SECRETACCESSKEY")
	minio_endpoint        = fmt.Sprintf("http://localhost:%v", os.Getenv("MINIO_SERVER_PORT"))
	minioClient           *minio.Client
)

func main() {
	if port == "" {
		fmt.Println("UPLOAD_PORT environment variable is not set")
		os.Exit(1)
	}
	if bucketName == "" {
		fmt.Println("BUCKET_NAME environment variable is not set")
		os.Exit(1)
	}
	if minio_accessKeyID == "" {
		fmt.Println("MINIO_ACCESSKEYID environment variable is not set")
		os.Exit(1)
	}
	if minio_secretAccessKey == "" {
		fmt.Println("MINIO_SECRETACCESSKEY environment variable is not set")
		os.Exit(1)
	}

	//init minio client
	useSSL := false
	minioclient, err := minio.New(minio_endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minio_accessKeyID, minio_secretAccessKey, ""),
		Secure: useSSL,
	})
	minioClient = minioclient
	if err != nil {
		log.Fatalln(err)
	}

	pool, err := pgxpool.New(db_ctx, dsn)
	if err != nil {
		log.Fatalf("UPLOAD DB err:%v", err)
	}
	err = pool.Ping(db_ctx)
	if err != nil {
		log.Fatalf("UPLOAD DB is unreachable: %v", err)
	}
	UploadRepo = repository.NewUploadRepository(pool)

	// --GRPC SERVER START--
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUploadServiceServer(s, &server{})
	log.Printf("UPLOAD service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve AUTH: %v", err)
	}
}
