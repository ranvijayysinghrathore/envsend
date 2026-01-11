package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/shamir"
)

// ShamirShare represents a single share of a secret.
type ShamirShare struct {
	Index int    `json:"index"`
	Data  string `json:"data"` // base64-encoded
}

// SplitSecretShamir splits a secret into n shares, requiring threshold shares to reconstruct.
// This is useful for distributing secrets among multiple parties.
// Example: Split a 32-byte key into 5 shares, requiring any 3 to reconstruct.
func SplitSecretShamir(secret []byte, threshold, totalShares int) ([]ShamirShare, error) {
	if threshold < 2 {
		return nil, fmt.Errorf("threshold must be at least 2")
	}
	if totalShares < threshold {
		return nil, fmt.Errorf("total shares must be >= threshold")
	}
	if threshold > 255 || totalShares > 255 {
		return nil, fmt.Errorf("threshold and total shares must be <= 255")
	}

	// Split secret using Shamir's Secret Sharing
	shares, err := shamir.Split(secret, totalShares, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to split secret: %w", err)
	}

	// Convert to ShamirShare format
	result := make([]ShamirShare, len(shares))
	for i, share := range shares {
		result[i] = ShamirShare{
			Index: i + 1,
			Data:  base64.StdEncoding.EncodeToString(share),
		}
	}

	return result, nil
}

// CombineSecretShamir reconstructs a secret from shares.
func CombineSecretShamir(shares []ShamirShare) ([]byte, error) {
	if len(shares) < 2 {
		return nil, fmt.Errorf("at least 2 shares required")
	}

	// Convert from ShamirShare format to byte slices
	shareBytes := make([][]byte, len(shares))
	for i, share := range shares {
		data, err := base64.StdEncoding.DecodeString(share.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode share %d: %w", share.Index, err)
		}
		shareBytes[i] = data
	}

	// Combine shares
	secret, err := shamir.Combine(shareBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to combine shares: %w", err)
	}

	return secret, nil
}

// GenerateRandomSecret generates a random secret of the specified size.
// Useful for generating keys that will be split using Shamir.
func GenerateRandomSecret(size int) ([]byte, error) {
	secret := make([]byte, size)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %w", err)
	}
	return secret, nil
}

// EncryptWithShamir encrypts data and splits the encryption key using Shamir's Secret Sharing.
// Returns the ciphertext, metadata, and shares of the encryption key.
func EncryptWithShamir(plaintext []byte, threshold, totalShares int) (ciphertext string, metadata EncryptionMetadata, shares []ShamirShare, err error) {
	// Generate encryption key
	key, err := GenerateKey()
	if err != nil {
		return "", metadata, nil, err
	}
	defer ZeroBytes(key)

	// Encrypt data
	ciphertext, metadata, err = EncryptAES256GCM(plaintext, key)
	if err != nil {
		return "", metadata, nil, err
	}

	// Split key into shares
	shares, err = SplitSecretShamir(key, threshold, totalShares)
	if err != nil {
		return "", metadata, nil, err
	}

	return ciphertext, metadata, shares, nil
}

// DecryptWithShamir decrypts data using shares to reconstruct the encryption key.
func DecryptWithShamir(ciphertext string, metadata EncryptionMetadata, shares []ShamirShare) ([]byte, error) {
	// Reconstruct key from shares
	key, err := CombineSecretShamir(shares)
	if err != nil {
		return nil, err
	}
	defer ZeroBytes(key)

	// Decrypt data
	plaintext, err := DecryptAES256GCM(ciphertext, key, metadata)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
