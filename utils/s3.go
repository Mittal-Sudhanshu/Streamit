package utils

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client
var S3Bucket string

func InitS3(bucket string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	S3Client = s3.NewFromConfig(cfg)
	S3Bucket = bucket
	return nil
}

func UploadFileToS3(localPath string, s3Key string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(s3Key),
		Body:   file,
	})

	return err
}
