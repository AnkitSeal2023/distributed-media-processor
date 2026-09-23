package main

import (
	"distributed-media-processing-platform/constants/endpoints"
	auth_pb "distributed-media-processing-platform/proto/generated/proto/auth/v1"
	upload_pb "distributed-media-processing-platform/proto/generated/proto/upload/v1"
	"log"
	"net/http"
	"os"

	"flag"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	accessKeyJwtSecret = os.Getenv("ACCESS_KEY_JWT_SECRET")
	api_port           = os.Getenv("API_PORT")
	auth_port          = os.Getenv("AUTH_PORT")
	upload_port        = os.Getenv("UPLOAD_PORT")
)

func main() {
	if accessKeyJwtSecret == "" {
		log.Fatal("ACCESS_KEY_JWT_SECRET is not set")
	}
	if api_port == "" {
		log.Fatal("API_PORT is not set")
	}
	if auth_port == "" {
		log.Fatal("AUTH_PORT is not set")
	}
	if upload_port == "" {
		log.Fatal("UPLOAD_PORT is not set")
	}

	// auth service client init
	auth_addr := flag.String("auth_addr", "localhost:"+auth_port, "AUTH service Port")
	auth_grpc_conn, err := grpc.NewClient(*auth_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	auth_grpc_client := auth_pb.NewAuthServiceClient(auth_grpc_conn)
	if err != nil {
		log.Fatalf("API could not connect to AUTH grpc service: %v", err)
	}

	// upload service client init
	upload_adddr := flag.String("upload_addr", "localhost:"+upload_port, "UPLOAD service Port")
	upload_grpc_conn, err := grpc.NewClient(*upload_adddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	upload_grpc_client := upload_pb.NewUploadServiceClient(upload_grpc_conn)
	if err != nil {
		log.Fatalf("API could not connect to UPLOAD grpc service: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// unprotected routes
	router.Group(func(r chi.Router) {
		r.Post(endpoints.AuthEndpoints.SignUp, func(w http.ResponseWriter, r *http.Request) {
			HandleRegisterUser(auth_grpc_client, r.Context(), w, r)
		})
		r.Post(endpoints.AuthEndpoints.SignIn, func(w http.ResponseWriter, r *http.Request) {
			HandleSigninUser(auth_grpc_client, r.Context(), w, r)
		})
		r.Get(endpoints.Ping, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Pong"))
		})
	})

	// protected routes
	router.Group(func(chi_router chi.Router) {
		chi_router.Use(AuthMiddleware)
		chi_router.Post(endpoints.UserEndpoints.UploadVideo, func(w http.ResponseWriter, r *http.Request) {
			HandleUploadVideo(upload_grpc_client, r.Context(), w, r)
		})
	})

	log.Printf("API Server is running on port %v", api_port)
	err = http.ListenAndServe("localhost:"+api_port, router)
	if err != nil {
		log.Fatalf("Failed to start API server %v", err)
	}

	defer auth_grpc_conn.Close()
}
