package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/lxmwaniky/file-storage-service/internal/domain"
)

type S3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
}

func NewS3Storage(client *s3.Client, bucketName string) domain.StorageService {
	return &S3Storage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    bucketName,
	}
}

// Upload to s3
func (s *S3Storage) Save(ctx context.Context, key string, src io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   src,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to s3: %w", err)
	}
	return nil
}

// Get retrieves an object from s3
func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return nil, domain.ErrObjectNotFound
		}
		return nil, fmt.Errorf("failed to get object from s3: %w", err)
	}
	return output.Body, nil
}

// Delete from s3
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from s3: %w", err)
	}
	return nil
}

func (s *S3Storage) GeneratePresignedUploadURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	request, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	return request.URL, nil
}
