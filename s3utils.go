package s3utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// CheckS3FileExists checks if a file exists in the S3 bucket
func CheckS3FileExists(sess *session.Session, bucket, key string) (bool, error) {
	svc := s3.New(sess)
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GenerateUniqueFileName generates a unique file name for S3
func GenerateUniqueFileName(sess *session.Session, bucket, folder, baseName string) (string, error) {
	// Get file extension
	ext := filepath.Ext(baseName)
	nameWithoutExt := baseName[:len(baseName)-len(ext)]

	// First, try the original filename
	key := filepath.Join(folder, baseName)
	exists, err := CheckS3FileExists(sess, bucket, key)
	if err != nil {
		return "", err
	}
	if !exists {
		return baseName, nil
	}

	// If the original filename exists, start appending numbers
	for i := 1; ; i++ {
		fileName := fmt.Sprintf("%s_%d%s", nameWithoutExt, i, ext)
		key = filepath.Join(folder, fileName)
		exists, err := CheckS3FileExists(sess, bucket, key)
		if err != nil {
			return "", err
		}
		if !exists {
			return fileName, nil
		}
	}
}

// UploadToS3 uploads a file to S3
func UploadToS3(region, profile, fileName, bucket, folder string) error {
	sess, err := NewAWSSession(region, profile)
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	key := filepath.Join(folder, filepath.Base(fileName))

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return err
	}
	return nil
}

// NewAWSSession creates a new AWS session
func NewAWSSession(region, profile string) (*session.Session, error) {
	// First try to use environment variables if they exist
	if hasEnvCredentials() {
		fmt.Println("Using AWS credentials from environment variables")
		return session.NewSession(&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
				SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
			}),
		})
	}

	// Fallback to profile-based credentials
	fmt.Printf("Using AWS credentials from profile: %s\n", profile)
	return session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
}

// hasEnvCredentials checks if all required AWS credentials are present in environment variables
func hasEnvCredentials() bool {
	requiredEnvVars := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			fmt.Println("AWS env credentials not fully set, falling back to profile")
			return false
		}
	}
	return true
}
