package main

import (
	"context"
	"log"

	"distributed-media-processing-platform/constants/error_msgs"
	pb "distributed-media-processing-platform/proto/generated/proto/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	hashPassword, err := hashPassword(in.GetPassword())
	userID, err := AuthRepo.CreateUser(ctx, in.GetEmail(), hashPassword)
	if err != nil {
		if err == error_msgs.ErrUserAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	access_token, err := generateAccessToken(userID)
	if err != nil {
		log.Printf("Failed to generate access token for user: %v", userID)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RegisterUserResponse{Message: "ok", AccessToken: access_token}, nil
}
