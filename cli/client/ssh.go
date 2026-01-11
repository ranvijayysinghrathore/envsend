package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GitHubPublicKey represents a GitHub user's SSH public key.
type GitHubPublicKey struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

// FetchGitHubSSHKey fetches a user's SSH public key from GitHub.
// GitHub username format: "github:username"
func FetchGitHubSSHKey(username string) (string, error) {
	// Remove "github:" prefix if present
	username = strings.TrimPrefix(username, "github:")

	// GitHub API endpoint
	url := fmt.Sprintf("https://api.github.com/users/%s/keys", username)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", UserAgent)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch GitHub keys: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("GitHub user not found: %s", username)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var keys []GitHubPublicKey
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if user has any keys
	if len(keys) == 0 {
		return "", fmt.Errorf("user %s has no SSH keys on GitHub", username)
	}

	// Iterate through keys to find a supported one (ssh-ed25519)
	for _, k := range keys {
		if strings.HasPrefix(k.Key, "ssh-ed25519") {
			return k.Key, nil
		}
	}

	// If we get here, we found keys but none were Ed25519
	var foundTypes []string
	for _, k := range keys {
		parts := strings.Fields(k.Key)
		if len(parts) > 0 {
			foundTypes = append(foundTypes, parts[0])
		}
	}

	return "", fmt.Errorf("user %s has %d key(s) (%v), but none are 'ssh-ed25519'. EnvSend requires Ed25519 keys for security", username, len(keys), foundTypes)
}

// GitLabPublicKey represents a GitLab user's SSH public key.
type GitLabPublicKey struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

// FetchGitLabSSHKey fetches a user's SSH public key from GitLab.
// GitLab username format: "gitlab:username"
func FetchGitLabSSHKey(username string) (string, error) {
	// Remove "gitlab:" prefix if present
	username = strings.TrimPrefix(username, "gitlab:")

	// GitLab API endpoint (gitlab.com)
	url := fmt.Sprintf("https://gitlab.com/api/v4/users/%s/keys", username)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch GitLab keys: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("GitLab user not found: %s", username)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitLab API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var keys []GitLabPublicKey
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if user has any keys
	if len(keys) == 0 {
		return "", fmt.Errorf("user %s has no SSH keys on GitLab", username)
	}

	// Iterate through keys to find a supported one (ssh-ed25519)
	for _, k := range keys {
		if strings.HasPrefix(k.Key, "ssh-ed25519") {
			return k.Key, nil
		}
	}

	// If we get here, we found keys but none were Ed25519
	var foundTypes []string
	for _, k := range keys {
		parts := strings.Fields(k.Key)
		if len(parts) > 0 {
			foundTypes = append(foundTypes, parts[0])
		}
	}

	return "", fmt.Errorf("user %s has %d key(s) (%v), but none are 'ssh-ed25519'. EnvSend requires Ed25519 keys for security", username, len(keys), foundTypes)
}

// FetchSSHKey fetches an SSH public key based on the recipient identifier.
// Supports formats: "github:username", "gitlab:username"
func FetchSSHKey(recipientID string) (string, error) {
	if strings.HasPrefix(recipientID, "github:") {
		return FetchGitHubSSHKey(recipientID)
	}
	if strings.HasPrefix(recipientID, "gitlab:") {
		return FetchGitLabSSHKey(recipientID)
	}
	return "", fmt.Errorf("unsupported recipient format: %s (use github:username or gitlab:username)", recipientID)
}
