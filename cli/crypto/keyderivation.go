package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters (OWASP recommended for password hashing)
	Argon2Time      = 3      // Number of iterations
	Argon2Memory    = 64 * 1024 // Memory in KiB (64 MB)
	Argon2Threads   = 4      // Number of threads
	Argon2KeyLength = 32     // Output key length (256 bits)
	Argon2SaltSize  = 16     // Salt size in bytes
)

// DeriveKeyFromPassphrase derives a 32-byte encryption key from a passphrase using Argon2id.
// Returns the derived key and the salt used (salt must be stored for decryption).
func DeriveKeyFromPassphrase(passphrase string) (key []byte, salt []byte, err error) {
	// Generate random salt
	salt = make([]byte, Argon2SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key using Argon2id
	key = argon2.IDKey(
		[]byte(passphrase),
		salt,
		Argon2Time,
		Argon2Memory,
		Argon2Threads,
		Argon2KeyLength,
	)

	return key, salt, nil
}

// DeriveKeyFromPassphraseWithSalt derives a key using a known salt (for decryption).
func DeriveKeyFromPassphraseWithSalt(passphrase string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(passphrase),
		salt,
		Argon2Time,
		Argon2Memory,
		Argon2Threads,
		Argon2KeyLength,
	)
}

// EncodeSalt encodes salt to base64 for storage.
func EncodeSalt(salt []byte) string {
	return base64.StdEncoding.EncodeToString(salt)
}

// DecodeSalt decodes base64-encoded salt.
func DecodeSalt(encoded string) ([]byte, error) {
	salt, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}
	return salt, nil
}

// EncryptWithPassphrase encrypts plaintext using a passphrase-derived key.
// Returns ciphertext and metadata including the salt.
func EncryptWithPassphrase(plaintext []byte, passphrase string) (ciphertext string, metadata EncryptionMetadata, err error) {
	// Derive key from passphrase
	key, salt, err := DeriveKeyFromPassphrase(passphrase)
	if err != nil {
		return "", metadata, err
	}
	defer ZeroBytes(key) // Clear key from memory after use

	// Encrypt with derived key
	ciphertext, metadata, err = EncryptAES256GCM(plaintext, key)
	if err != nil {
		return "", metadata, err
	}

	// Update metadata to include salt and key derivation info
	metadata.KeyDerivation = "argon2id"
	metadata.Salt = EncodeSalt(salt)

	return ciphertext, metadata, nil
}

// DecryptWithPassphrase decrypts ciphertext using a passphrase.
func DecryptWithPassphrase(ciphertext string, passphrase string, metadata EncryptionMetadata) ([]byte, error) {
	// Validate key derivation method
	if metadata.KeyDerivation != "argon2id" {
		return nil, fmt.Errorf("unsupported key derivation: %s", metadata.KeyDerivation)
	}

	// Decode salt
	salt, err := DecodeSalt(metadata.Salt)
	if err != nil {
		return nil, err
	}

	// Derive key from passphrase
	key := DeriveKeyFromPassphraseWithSalt(passphrase, salt)
	defer ZeroBytes(key) // Clear key from memory after use

	// Decrypt
	plaintext, err := DecryptAES256GCM(ciphertext, key, metadata)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
