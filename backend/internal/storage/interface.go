package storage

import "context"

// BlobStorage is an interface for encrypted blob storage.
// Implementations: S3Storage, LocalStorage
type BlobStorage interface {
	UploadEncryptedBlob(ctx context.Context, data []byte) (string, error)
	DownloadEncryptedBlob(ctx context.Context, objectURL string) ([]byte, error)
	DeleteEncryptedBlob(ctx context.Context, objectURL string) error
}
