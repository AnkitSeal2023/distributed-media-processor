package main

import (
	"context"
	"fmt"

	pb "distributed-media-processing-platform/proto/generated/proto/upload/v1"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *server) UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	filename := req.GetFileName()
	filesize := req.GetFileSize()
	filetype := req.GetFileType()
	fmt.Printf("%v %v %v", filename, filesize, filetype)
	if len(filename) >= 50 {
		return nil, status.Error(codes.InvalidArgument, "File name should be less than 50 characters")
	}

	if filesize > 1000000000 {
		return nil, status.Error(codes.InvalidArgument, "File size should be less than 1GB")
	}

	return &pb.UploadVideoResponse{
		PresignedUrl: "hello from upload service",
	}, nil

}
