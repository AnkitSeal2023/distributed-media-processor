package handlers

import (
	"context"
	"encoding/json"
	"flag"
	_ "fmt"
	"log"
	"net/http"
	"time"

	"distributed-media-processing-platform/constants/error_msgs"
	pb "distributed-media-processing-platform/proto/generated/proto/auth/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	AccessToken string `json:"access_token"`
}

func HandleRegisterUser(conn *grpc.ClientConn, w http.ResponseWriter, r *http.Request) {
	flag.Parse()
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
		handleGrpcError(w, err)
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

func handleGrpcError(w http.ResponseWriter, err error) {
	switch status.Code(err) {
	case codes.AlreadyExists:
		writeError(w, http.StatusConflict, error_msgs.ErrUserAlreadyExists.Error())
	default:
		writeError(w, http.StatusInternalServerError, error_msgs.ErrInternalServer.Error())
	}

}

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Printf("Failed to write message")
	}
}
