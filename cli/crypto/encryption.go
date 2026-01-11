package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidKeySize     = errors.New("invalid key size: must be 32 bytes for AES-256")
	ErrInvalidCiphertext  = errors.New("invalid ciphertext: too short or malformed")
	ErrDecryptionFailed   = errors.New("decryption failed: authentication tag mismatch")
)

// EncryptionMetadata contains information needed for decryption
type EncryptionMetadata struct {
	Algorithm    string `json:"algorithm"`     // "AES-256-GCM"
	IV           string `json:"iv"`            // base64-encoded initialization vector
	KeyDerivation string `json:"keyDerivation"` // "argon2id" or "none"
	Salt         string `json:"salt,omitempty"` // base64-encoded salt (if key derivation used)
	Version      string `json:"version"`       // "1.0"
}

// EncryptAES256GCM encrypts plaintext using AES-256-GCM with the provided 32-byte key.
// Returns base64-encoded ciphertext and metadata needed for decryption.
// The key must be exactly 32 bytes (256 bits).
func EncryptAES256GCM(plaintext []byte, key []byte) (ciphertext string, metadata EncryptionMetadata, err error) {
	if len(key) != 32 {
		return "", metadata, ErrInvalidKeySize
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", metadata, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", metadata, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce/IV
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", metadata, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and authenticate
	ciphertextBytes := gcm.Seal(nil, nonce, plaintext, nil)

	// Encode to base64 for transport
	ciphertext = base64.StdEncoding.EncodeToString(ciphertextBytes)

	// Create metadata
	metadata = EncryptionMetadata{
		Algorithm:    "AES-256-GCM",
		IV:           base64.StdEncoding.EncodeToString(nonce),
		KeyDerivation: "none",
		Version:      "1.0",
	}

	return ciphertext, metadata, nil
}

// DecryptAES256GCM decrypts base64-encoded ciphertext using AES-256-GCM.
// Returns the plaintext bytes.
func DecryptAES256GCM(ciphertext string, key []byte, metadata EncryptionMetadata) ([]byte, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	// Validate algorithm
	if metadata.Algorithm != "AES-256-GCM" {
		return nil, fmt.Errorf("unsupported algorithm: %s", metadata.Algorithm)
	}

	// Decode base64 ciphertext
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	// Decode IV
	nonce, err := base64.StdEncoding.DecodeString(metadata.IV)
	if err != nil {
		return nil, fmt.Errorf("failed to decode IV: %w", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Validate nonce size
	if len(nonce) != gcm.NonceSize() {
		return nil, ErrInvalidCiphertext
	}

	// Decrypt and verify authentication tag
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// GenerateKey generates a cryptographically secure random 32-byte key for AES-256.
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}

// EncodeKey encodes a key to base64 for sharing (e.g., in URLs or QR codes).
func EncodeKey(key []byte) string {
	return base64.URLEncoding.EncodeToString(key)
}

// DecodeKey decodes a base64-encoded key.
func DecodeKey(encoded string) ([]byte, error) {
	key, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}
	return key, nil
}

// MetadataToJSON converts encryption metadata to JSON string.
func MetadataToJSON(metadata EncryptionMetadata) (string, error) {
	data, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}
	return string(data), nil
}

// MetadataFromJSON parses encryption metadata from JSON string.
func MetadataFromJSON(jsonStr string) (EncryptionMetadata, error) {
	var metadata EncryptionMetadata
	if err := json.Unmarshal([]byte(jsonStr), &metadata); err != nil {
		return metadata, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}
	return metadata, nil
}
