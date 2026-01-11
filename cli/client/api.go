package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
	UserAgent      = "envsend-cli/1.0"
)

// APIClient handles communication with the EnvSend backend.
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAPIClient creates a new API client.
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// UploadSecretRequest represents the request to upload an encrypted secret.
type UploadSecretRequest struct {
	EncryptedBlob      string                 `json:"encrypted_blob"`
	EncryptionMetadata map[string]interface{} `json:"encryption_metadata"`
	ExpiresIn          string                 `json:"expires_in"`          // e.g., "10m", "1h"
	MaxViews           int                    `json:"max_views"`
	IPLock             string                 `json:"ip_lock,omitempty"`
	RecipientID        string                 `json:"recipient_id,omitempty"`
}

// UploadSecretResponse represents the response from uploading a secret.
type UploadSecretResponse struct {
	SecretID  string    `json:"secret_id"`
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
	MaxViews  int       `json:"max_views"`
}

// DownloadSecretResponse represents the response from downloading a secret.
type DownloadSecretResponse struct {
	EncryptedBlob      string                 `json:"encrypted_blob"`
	EncryptionMetadata map[string]interface{} `json:"encryption_metadata"`
	ViewsRemaining     int                    `json:"views_remaining"`
	ExpiresAt          time.Time              `json:"expires_at"`
}

// UploadSecret uploads an encrypted secret to the server.
func (c *APIClient) UploadSecret(req UploadSecretRequest) (*UploadSecretResponse, error) {
	// Marshal request
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/api/v1/secrets", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", UserAgent)

	// Send request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var uploadResp UploadSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &uploadResp, nil
}

// DownloadSecret downloads an encrypted secret from the server.
func (c *APIClient) DownloadSecret(secretID string) (*DownloadSecretResponse, error) {
	// Create HTTP request
	httpReq, err := http.NewRequest("GET", c.BaseURL+"/api/v1/secrets/"+secretID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("User-Agent", UserAgent)

	// Send request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("secret not found or has expired")
	}
	if resp.StatusCode == http.StatusGone {
		return nil, fmt.Errorf("secret has been destroyed (max views reached)")
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var downloadResp DownloadSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&downloadResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &downloadResp, nil
}

// HealthCheck checks if the server is healthy.
func (c *APIClient) HealthCheck() error {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/health")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server unhealthy: status %d", resp.StatusCode)
	}

	return nil
}
