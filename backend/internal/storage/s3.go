package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ranvijayysinghrathore/envsend/backend/internal/config"
)

// S3Storage handles encrypted blob storage in MinIO/S3.
type S3Storage struct {
	client *minio.Client
	bucket string
}

// NewS3Storage creates a new S3 storage client.
func NewS3Storage(cfg config.S3Config) (*S3Storage, error) {
	// Initialize MinIO client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{Region: cfg.Region})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &S3Storage{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

// UploadEncryptedBlob uploads an encrypted blob to S3/MinIO.
// Returns the object URL.
func (s *S3Storage) UploadEncryptedBlob(ctx context.Context, data []byte) (string, error) {
	// Generate unique object key
	objectKey := fmt.Sprintf("secrets/%s.enc", uuid.New().String())

	// Upload to MinIO
	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(ctx, s.bucket, objectKey, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	// Return object URL (internal reference)
	return fmt.Sprintf("s3://%s/%s", s.bucket, objectKey), nil
}

// DownloadEncryptedBlob downloads an encrypted blob from S3/MinIO.
func (s *S3Storage) DownloadEncryptedBlob(ctx context.Context, objectURL string) ([]byte, error) {
	// Parse object URL to extract key
	// Format: s3://bucket/key
	objectKey := extractObjectKey(objectURL)

	// Download from MinIO
	object, err := s.client.GetObject(ctx, s.bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	// Read all data
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return data, nil
}

// DeleteEncryptedBlob deletes an encrypted blob from S3/MinIO.
func (s *S3Storage) DeleteEncryptedBlob(ctx context.Context, objectURL string) error {
	objectKey := extractObjectKey(objectURL)

	err := s.client.RemoveObject(ctx, s.bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// extractObjectKey extracts the object key from an S3 URL.
// Format: s3://bucket/key -> key
func extractObjectKey(objectURL string) string {
	// Simple parsing: remove "s3://bucket/" prefix
	const prefix = "s3://"
	if len(objectURL) > len(prefix) {
		// Find first "/" after bucket name
		start := len(prefix)
		for i := start; i < len(objectURL); i++ {
			if objectURL[i] == '/' {
				return objectURL[i+1:]
			}
		}
	}
	return objectURL
}
