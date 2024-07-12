package s3helper

import (
	"amartha/loan-service/helpers/envhelper"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func CreateClientFromSession() (*s3.S3, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(envhelper.GetEnvVar("AWS_DEPLOYMENT_REGION")),
	})
	if err != nil {
		err = fmt.Errorf("error creating SES client: %s", err.Error())
		return nil, err
	}
	return s3.New(session), nil
}

const DefaultExpireDurationForPresignedURL = 24 * time.Hour

func UploadFile(
	s3Client *s3.S3, file io.Reader, bucketName, keyName string,
) (
	err error,
) {
	uploader := s3manager.NewUploaderWithClient(s3Client)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})
	if err != nil {
		err = fmt.Errorf("error uploading file %s/%s: %v", bucketName, keyName, err)
		return err
	}

	return nil
}

func GeneratePresignedURLForGetObject(s3Client *s3.S3, bucketName, keyName string) (
	presignedUrl string,
	err error,
) {
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	presignedUrl, err = req.Presign(DefaultExpireDurationForPresignedURL)
	if err != nil {
		err = fmt.Errorf("error presign GetObjectRequest for %s/%s: %v", bucketName, keyName, err)
		return presignedUrl, err
	}

	return presignedUrl, nil
}

func GeneratePresignedURLForPutObject(s3Client *s3.S3, bucketName, keyName string) (
	presignedUrl string,
	err error,
) {
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	presignedUrl, err = req.Presign(DefaultExpireDurationForPresignedURL)
	if err != nil {
		err = fmt.Errorf("error presign PutObjectRequest for %s/%s: %v", bucketName, keyName, err)
		return presignedUrl, err
	}

	return presignedUrl, nil
}
