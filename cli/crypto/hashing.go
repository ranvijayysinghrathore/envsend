package crypto

import (
	"encoding/hex"
	"fmt"

	"github.com/zeebo/blake3"
)

// HashBLAKE3 computes the BLAKE3 hash of data.
// Returns the hash as a hex-encoded string.
func HashBLAKE3(data []byte) string {
	hash := blake3.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// HashBLAKE3Bytes computes the BLAKE3 hash and returns raw bytes.
func HashBLAKE3Bytes(data []byte) []byte {
	hash := blake3.Sum256(data)
	return hash[:]
}

// VerifyBLAKE3 verifies that data matches the expected hash.
func VerifyBLAKE3(data []byte, expectedHash string) (bool, error) {
	actualHash := HashBLAKE3(data)
	return actualHash == expectedHash, nil
}

// HashFile computes BLAKE3 hash of file contents (for integrity verification).
func HashFile(content []byte) string {
	return HashBLAKE3(content)
}

// GenerateSecretID generates a unique identifier for a secret using BLAKE3.
// This can be used as a deterministic ID based on content + timestamp.
func GenerateSecretID(content []byte, timestamp string) string {
	combined := append(content, []byte(timestamp)...)
	return HashBLAKE3(combined)
}

// ComputeChecksum computes a checksum for encrypted data integrity verification.
func ComputeChecksum(ciphertext string) string {
	return HashBLAKE3([]byte(ciphertext))
}

// VerifyChecksum verifies the integrity of encrypted data.
func VerifyChecksum(ciphertext string, expectedChecksum string) error {
	actualChecksum := ComputeChecksum(ciphertext)
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}
	return nil
}
