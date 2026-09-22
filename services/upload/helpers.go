package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

func generatePresignedUrl(filename string) (presignedUrl string, formData map[string]string, err error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+fmt.Sprintf("%v", filename))
	expiry := time.Second * 5 * 60 // 5mins.
	presignedURL, err := minioClient.PresignedPutObject(context.Background(), bucketName, filename, expiry)
	if err != nil {
		fmt.Println("Error generating presigned URL", err)
		return "", nil, err
	}
	// return presignedURL.RawPath, nil
	presignedURLstr := presignedURL.Host + presignedURL.Path

	// Initialize policy condition config.
	policy := minio.NewPostPolicy()

	// Apply upload policy restrictions:
	policy.SetBucket(bucketName)
	policy.SetKey(filename)
	policy.SetExpires(time.Now().UTC().AddDate(0, 0, 10)) // expires in 10 days

	// Only allow 'png' images.
	policy.SetContentType("video/mp4")

	// Only allow content size in range 1KB to 1MB.
	policy.SetContentLengthRange(1024, 1024*1024*1024)

	// Add a user metadata using the key "custom" and value "user"
	policy.SetUserMetadata("custom", "user")

	// Get the POST form key/value object:
	url, formData, err := minioClient.PresignedPostPolicy(context.Background(), policy)
	if err != nil {
		fmt.Println(err)
		return "", nil, err
	}
	fmt.Println("Presigned POST URL:", url)
	fmt.Println("\nPresigned POST form data:\n", formData)

	return presignedURLstr, formData, nil
}
