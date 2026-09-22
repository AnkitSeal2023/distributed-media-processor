package main

import (
	"context"
	"encoding/json"
	_ "fmt"
	"log"
	"net/http"
	"time"

	auth_pb "distributed-media-processing-platform/proto/generated/proto/auth/v1"
	upload_pb "distributed-media-processing-platform/proto/generated/proto/upload/v1"
)

type UserAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type UserSigninResponse struct {
	AccessToken string `json:"access_token"`
}

type UploadVideoResponse struct {
	PresignedUrl string            `json:"presigned_url"`
	FormData     map[string]string `json:"form_data"`
}

type UploadVideoRequest struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}

func HandleRegisterUser(c auth_pb.AuthServiceClient, ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ur UserAuthRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		log.Printf("could not decode JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rpb, err := c.RegisterUser(ctx, &auth_pb.RegisterUserRequest{Email: ur.Email, Password: ur.Password})
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

func HandleSigninUser(c auth_pb.AuthServiceClient, ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ur UserAuthRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		log.Printf("could not decode JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rpb, err := c.SigninUser(ctx, &auth_pb.SigninUserRequest{Email: ur.Email, Password: ur.Password})
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

func HandleUploadVideo(c upload_pb.UploadServiceClient, ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ur UploadVideoRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		log.Printf("could not decode JSON: %v", err)
		EncodeResponse(w, nil, "internal server error", nil, http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rpb, err := c.UploadVideo(ctx, &upload_pb.UploadVideoRequest{FileName: ur.FileName, FileType: ur.FileType, FileSize: ur.FileSize})
	if err != nil {
		HandleGrpcError(w, err)
		return
	}

	response := UploadVideoResponse{
		PresignedUrl: rpb.GetPresignedUrl(),
		FormData:     rpb.GetFormData(),
	}

	EncodeResponse(w, response, "ok", nil, http.StatusOK)
}
