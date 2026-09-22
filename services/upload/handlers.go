package main

import (
	"context"

	pb "distributed-media-processing-platform/proto/generated/proto/upload/v1"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *server) UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	filename := req.GetFileName()
	filesize := req.GetFileSize()
	filetype := req.GetFileType()

	if len(filename) >= 50 {
		return nil, status.Error(codes.InvalidArgument, "File name should be less than 50 characters")
	}
	if filesize > 1000000000 {
		return nil, status.Error(codes.InvalidArgument, "File size should be less than 1GB")
	}
	if !(filetype == "video/mp4" || filetype == "video/mkv" || filetype == "video/avi") {
		return nil, status.Error(codes.InvalidArgument, "File type should be video/mp4, video/mkv or video/avi")
	}

	presigned_url, formData, err := generatePresignedUrl(filename)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error generating presigned URL")
	}

	return &pb.UploadVideoResponse{
		PresignedUrl: presigned_url,
		FormData:     formData,
	}, nil

}
