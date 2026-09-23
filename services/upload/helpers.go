package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

func generatePresignedUrl(filename string, filesize int64) (presignedUrl string, formData map[string]string, err error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+fmt.Sprintf("%v", filename))

	policy := minio.NewPostPolicy()

	policy.SetBucket(bucketName)
	policy.SetKey(filename)
	policy.SetExpires(time.Now().UTC().Add(15 * time.Minute))

	policy.SetContentType("video/mp4")
	policy.SetContentLengthRange(0, filesize)
	policy.SetUserMetadata("custom", "user")

	url, formData, err := minioClient.PresignedPostPolicy(context.Background(), policy)
	if err != nil {
		log.Printf("Error while generating presigned POST url:%v", err)
		return "", nil, err
	}

	return url.Host + url.Path, formData, nil
}
