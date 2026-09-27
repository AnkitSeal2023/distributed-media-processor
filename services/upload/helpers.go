package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

func generatePresignedUrl(filename string, filesize int64, filetype string) (presignedUrl string, formData map[string]string, err error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+fmt.Sprintf("%v", filename))
	var extension string
	if strings.Contains(filetype, "/") {
		fileExtension := strings.Split(filetype, "/")
		if len(fileExtension) != 2 {
			return "", nil, errors.New("wrong file type")
		}
		extension = fileExtension[1]
	} else {
		return "", nil, errors.New("wrong file type")
	}

	policy := minio.NewPostPolicy()

	policy.SetBucket(bucketName)
	policy.SetKey(filename + "." + extension)
	policy.SetExpires(time.Now().UTC().Add(15 * time.Minute))

	policy.SetContentType(filetype)
	policy.SetContentLengthRange(0, filesize)
	policy.SetUserMetadata("custom", "user")

	url, formData, err := minioClient.PresignedPostPolicy(context.Background(), policy)
	if err != nil {
		log.Printf("Error while generating presigned POST url:%v", err)
		return "", nil, err
	}

	return url.Host + url.Path, formData, nil
}
