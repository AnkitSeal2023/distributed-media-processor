package main

import (
	"context"
	"encoding/json"
	_ "fmt"
	"log"
	"net/http"
	"time"

	pb "distributed-media-processing-platform/proto/generated/proto/auth/v1"

	"google.golang.org/grpc"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type UserSigninResponse struct {
	AccessToken string `json:"access_token"`
}

func RegisterUser(conn *grpc.ClientConn, w http.ResponseWriter, r *http.Request) {
	var ur UserRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		log.Printf("could not decode JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rpb, err := c.RegisterUser(ctx, &pb.RegisterUserRequest{Email: ur.Email, Password: ur.Password})
	if err != nil {
		HandleGrpcError(w, err)
		return
	}

	access_token := rpb.AccessToken
	response := UserRegisterResponse{
		AccessToken: access_token,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("could not encode JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SigninUser(conn *grpc.ClientConn, w http.ResponseWriter, r *http.Request) {
	var ur UserRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		log.Printf("could not decode JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rpb, err := c.SigninUser(ctx, &pb.SigninUserRequest{Email: ur.Email, Password: ur.Password})
	if err != nil {
		HandleGrpcError(w, err)
		return
	}

	access_token := rpb.AccessToken
	response := UserSigninResponse{
		AccessToken: access_token,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("could not encode JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
