package urlify

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/docker/docker/pkg/namesgenerator"
)

func Urlify(bucket string, prefix string, region string, buffer []byte) (string, error) {

	// get a session object so we can interact with AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Printf("Could not load session in %s: %v\n", region, err)
		return "", errors.New(fmt.Sprintf("Could not load session in %s: %v", region, err))
	}

	// upload the buffer to the given S3 bucket at the given prefix
	// use a random name
	uploader := s3manager.NewUploader(sess)
	fileName := namesgenerator.GetRandomName(3)
	extension := filepath.Ext(fileName)
	key := fmt.Sprintf("%s/%s%s", prefix, fileName, extension)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buffer),
	})

	if err != nil {
		log.Printf("Could not load session in %s: %v\n", region, err)
		return "", errors.New(fmt.Sprintf("Could not upload file to %s: %v", key, err))
	}

	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(5 * time.Minute)

	if err != nil {
		log.Printf("Failed to presign request: %v\n", err)
		return "", errors.New(fmt.Sprintf("Failed to presign request: %v\n", err))
	}
	return urlStr, nil
}
