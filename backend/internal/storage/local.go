package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// LocalStorage handles encrypted blob storage on local disk.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local disk storage.
func NewLocalStorage() (*LocalStorage, error) {
	// Use /tmp for Render, or current directory for local dev
	basePath := os.Getenv("LOCAL_STORAGE_PATH")
	if basePath == "" {
		basePath = filepath.Join(os.TempDir(), "envsend-secrets")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
	}, nil
}

// UploadEncryptedBlob saves an encrypted blob to local disk.
func (l *LocalStorage) UploadEncryptedBlob(ctx context.Context, data []byte) (string, error) {
	// Generate unique filename
	filename := fmt.Sprintf("%s.enc", uuid.New().String())
	filepath := filepath.Join(l.basePath, filename)

	// Write file
	if err := ioutil.WriteFile(filepath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return local:// URL
	return fmt.Sprintf("local://%s", filename), nil
}

// DownloadEncryptedBlob reads an encrypted blob from local disk.
func (l *LocalStorage) DownloadEncryptedBlob(ctx context.Context, objectURL string) ([]byte, error) {
	// Extract filename from URL (format: local://filename.enc)
	filename := extractLocalFilename(objectURL)
	filepath := filepath.Join(l.basePath, filename)

	// Read file
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("blob not found")
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// DeleteEncryptedBlob removes an encrypted blob from local disk.
func (l *LocalStorage) DeleteEncryptedBlob(ctx context.Context, objectURL string) error {
	filename := extractLocalFilename(objectURL)
	filepath := filepath.Join(l.basePath, filename)

	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// extractLocalFilename extracts filename from local:// URL
func extractLocalFilename(url string) string {
	// Remove "local://" prefix
	if len(url) > 8 && url[:8] == "local://" {
		return url[8:]
	}
	return url
}
