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
	port := os.Getenv("API_PORT")
	auth_port := os.Getenv("AUTH_PORT")
	addr := flag.String("addr", "localhost:"+auth_port, "AUTH service Port")
	grpc_conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// unprotected routes
	r.Group(func(r chi.Router) {
		r.Post(endpoints.UserEndpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			handlers.RegisterUser(grpc_conn, w, r)
		})
		r.Post(endpoints.UserEndpoints.SignIn, func(w http.ResponseWriter, r *http.Request) {
			handlers.SigninUser(grpc_conn, w, r)
		})
		r.Get(endpoints.Ping, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Pong"))
		})
	})

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
	})

	log.Printf("API Server is running on port %v", port)
	err = http.ListenAndServe("localhost:"+port, r)
	if err != nil {
		log.Fatalf("Failed to start API server %v", err)
	}

	defer grpc_conn.Close()
}
