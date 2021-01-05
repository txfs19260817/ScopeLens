package storage

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var S3Client *AmazonS3

// AmazonS3 is an Amazon S3 file storage.
type AmazonS3 struct {
	bucket string
	svc    *s3.S3
}

// NewAmazonS3 returns a new Amazon S3 file storage.
func NewAmazonS3(accessKey, secretKey, region, bucket string) (*AmazonS3, error) {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	cfg := aws.NewConfig().WithCredentials(creds).WithRegion(region)

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	return &AmazonS3{
		bucket: bucket,
		svc:    s3.New(sess),
	}, nil
}

// Save saves data from r to file with the given path.
func (s *AmazonS3) Save(path string, r io.Reader) (string, error) {
	contentType := aws.String("binary/octet-stream")

	switch ext := filepath.Ext(path); ext {
	case ".jpg", ".jpeg":
		contentType = aws.String("image/jpeg")
	case ".png":
		contentType = aws.String("image/png")
	default:
		contentType = aws.String("binary/octet-stream")
	}

	res, err := s3manager.NewUploaderWithClient(s.svc).Upload(
		&s3manager.UploadInput{
			Bucket:      aws.String(s.bucket),
			Key:         aws.String(path),
			ACL:         aws.String(s3.ObjectCannedACLPublicRead),
			Body:        r,
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload object to S3: %w", err)
	}
	return res.Location, nil
}

// Remove removes the file with the given path.
func (s *AmazonS3) Remove(path string) error {
	_, err := s.svc.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(path),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}
	return nil
}

// URL returns an URL of the file with the given path.
func (s *AmazonS3) URL(path, region string) string {
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", s.bucket, region, path)
}
