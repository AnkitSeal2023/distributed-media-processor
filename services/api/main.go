package main

import (
	"distributed-media-processing-platform/constants/endpoints"
	"distributed-media-processing-platform/services/api/handlers"
	"log"
	"net/http"
	"os"

	"flag"

	env "github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Fatal("Could not load API env")
	}
	port := "localhost:" + os.Getenv("API_PORT")
	addr := flag.String("addr", port, "the address to connect to")
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post(endpoints.UserEndpoints.Register, func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRegisterUser(conn, w, r)
	})
	log.Print("API Server is running on port 3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatalf("Failed to start api server %v", err)
	}

	defer conn.Close()
}
