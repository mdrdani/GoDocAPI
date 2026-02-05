package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"godocapi/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type RustFSStore struct {
	client *s3.Client
	bucket string
}

func NewRustFSStore(cfg *config.Config) (*RustFSStore, error) {
	// Load AWS config with custom endpoint for RustFS
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.RustFSRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.RustFSAccessKey, cfg.RustFSSecretKey, "")),
		awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               cfg.RustFSEndpoint,
					SigningRegion:     cfg.RustFSRegion,
					HostnameImmutable: true, // Needed for path-style buckets if RustFS requires it, usually safer with custom endpoints
				}, nil
			},
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // Use path-style addressing for S3 compatible implementations
	})

	return &RustFSStore{
		client: client,
		bucket: cfg.RustFSBucket,
	}, nil
}

func (s *RustFSStore) Upload(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (string, error) {
	// Create a unique key for the file
	key := fmt.Sprintf("%d/%s", time.Now().UnixNano(), filename)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to RustFS: %w", err)
	}

	return key, nil
}

func (s *RustFSStore) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from RustFS: %w", err)
	}

	return output.Body, nil
}

func (s *RustFSStore) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from RustFS: %w", err)
	}

	return nil
}
