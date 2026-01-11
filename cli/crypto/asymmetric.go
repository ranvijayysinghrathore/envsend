package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"
)

// X25519KeyPair represents an X25519 key pair for asymmetric encryption.
type X25519KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

// GenerateX25519KeyPair generates a new X25519 key pair.
func GenerateX25519KeyPair() (*X25519KeyPair, error) {
	// Generate private key (32 random bytes)
	privateKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, privateKey); err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Derive public key
	publicKey, err := curve25519.X25519(privateKey, curve25519.Basepoint)
	if err != nil {
		return nil, fmt.Errorf("failed to derive public key: %w", err)
	}

	return &X25519KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// ComputeSharedSecret computes the shared secret using X25519 key exchange.
// This is used for encrypting the symmetric key with the recipient's public key.
func ComputeSharedSecret(privateKey, publicKey []byte) ([]byte, error) {
	if len(privateKey) != 32 {
		return nil, fmt.Errorf("invalid private key size: expected 32 bytes, got %d", len(privateKey))
	}
	if len(publicKey) != 32 {
		return nil, fmt.Errorf("invalid public key size: expected 32 bytes, got %d", len(publicKey))
	}

	sharedSecret, err := curve25519.X25519(privateKey, publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute shared secret: %w", err)
	}

	return sharedSecret, nil
}

// EncryptSymmetricKeyForRecipient encrypts a symmetric key using the recipient's public key.
// Returns the encrypted key and the ephemeral public key (needed for decryption).
func EncryptSymmetricKeyForRecipient(symmetricKey, recipientPublicKey []byte) (encryptedKey, ephemeralPublicKey []byte, err error) {
	// Generate ephemeral key pair
	ephemeralKeyPair, err := GenerateX25519KeyPair()
	if err != nil {
		return nil, nil, err
	}
	defer ZeroBytes(ephemeralKeyPair.PrivateKey)

	// Compute shared secret
	sharedSecret, err := ComputeSharedSecret(ephemeralKeyPair.PrivateKey, recipientPublicKey)
	if err != nil {
		return nil, nil, err
	}
	defer ZeroBytes(sharedSecret)

	// Encrypt symmetric key with shared secret
	encryptedKeyStr, _, err := EncryptAES256GCM(symmetricKey, sharedSecret)
	if err != nil {
		return nil, nil, err
	}

	// Decode encrypted key from base64
	encryptedKey, err = base64.StdEncoding.DecodeString(encryptedKeyStr)
	if err != nil {
		return nil, nil, err
	}

	return encryptedKey, ephemeralKeyPair.PublicKey, nil
}

// DecryptSymmetricKeyFromSender decrypts a symmetric key using the recipient's private key.
func DecryptSymmetricKeyFromSender(encryptedKey, ephemeralPublicKey, recipientPrivateKey []byte, metadata EncryptionMetadata) ([]byte, error) {
	// Compute shared secret
	sharedSecret, err := ComputeSharedSecret(recipientPrivateKey, ephemeralPublicKey)
	if err != nil {
		return nil, err
	}
	defer ZeroBytes(sharedSecret)

	// Decrypt symmetric key
	encryptedKeyStr := base64.StdEncoding.EncodeToString(encryptedKey)
	symmetricKey, err := DecryptAES256GCM(encryptedKeyStr, sharedSecret, metadata)
	if err != nil {
		return nil, err
	}

	return symmetricKey, nil
}

// EncodePublicKey encodes a public key to base64 for sharing.
func EncodePublicKey(publicKey []byte) string {
	return base64.StdEncoding.EncodeToString(publicKey)
}

// DecodePublicKey decodes a base64-encoded public key.
func DecodePublicKey(encoded string) ([]byte, error) {
	publicKey, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}
	if len(publicKey) != 32 {
		return nil, fmt.Errorf("invalid public key size: expected 32 bytes, got %d", len(publicKey))
	}
	return publicKey, nil
}

// ParseSSHPublicKey extracts the X25519 public key from an SSH public key.
// For Ed25519 SSH keys, we convert them to X25519 format.
// Note: This is a simplified version. Production code should use crypto/ssh package.
func ParseSSHPublicKey(sshPublicKey string) ([]byte, error) {
	// TODO: Implement proper SSH public key parsing
	// For now, this is a placeholder that expects base64-encoded X25519 keys
	return DecodePublicKey(sshPublicKey)
}
