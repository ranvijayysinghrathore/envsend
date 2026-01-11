package crypto

import (
	"runtime"
	"unsafe"
)

// ZeroBytes securely zeros out a byte slice to prevent sensitive data from remaining in memory.
// This is critical for clearing encryption keys, passphrases, and decrypted secrets.
func ZeroBytes(data []byte) {
	if data == nil || len(data) == 0 {
		return
	}

	// Overwrite with zeros
	for i := range data {
		data[i] = 0
	}

	// Force garbage collector to not optimize away the zeroing
	runtime.KeepAlive(data)
}

// ZeroString securely zeros out a string by converting to bytes and zeroing.
// Note: Strings in Go are immutable, so this creates a mutable copy first.
func ZeroString(s *string) {
	if s == nil || *s == "" {
		return
	}

	// Convert string to byte slice (creates a copy)
	bytes := []byte(*s)
	ZeroBytes(bytes)

	// Clear the string pointer
	*s = ""
}

// SecureBuffer is a wrapper around a byte slice that automatically zeros on garbage collection.
type SecureBuffer struct {
	data []byte
}

// NewSecureBuffer creates a new secure buffer of the specified size.
func NewSecureBuffer(size int) *SecureBuffer {
	buf := &SecureBuffer{
		data: make([]byte, size),
	}

	// Set finalizer to zero memory when garbage collected
	runtime.SetFinalizer(buf, func(b *SecureBuffer) {
		b.Destroy()
	})

	return buf
}

// Bytes returns the underlying byte slice.
func (sb *SecureBuffer) Bytes() []byte {
	return sb.data
}

// Destroy explicitly zeros and releases the buffer.
func (sb *SecureBuffer) Destroy() {
	if sb.data != nil {
		ZeroBytes(sb.data)
		sb.data = nil
	}
}

// Size returns the size of the buffer.
func (sb *SecureBuffer) Size() int {
	if sb.data == nil {
		return 0
	}
	return len(sb.data)
}

// Copy copies data into the secure buffer.
func (sb *SecureBuffer) Copy(data []byte) error {
	if len(data) > len(sb.data) {
		return ErrInvalidKeySize
	}
	copy(sb.data, data)
	return nil
}

// WipeMemory attempts to overwrite a memory region (advanced usage).
// This is a best-effort function and may not work on all platforms.
func WipeMemory(ptr unsafe.Pointer, length int) {
	if ptr == nil || length <= 0 {
		return
	}

	// Create a slice from the pointer
	slice := (*[1 << 30]byte)(ptr)[:length:length]

	// Zero the memory
	for i := range slice {
		slice[i] = 0
	}

	runtime.KeepAlive(ptr)
}

// ClearSensitiveData is a helper to clear multiple sensitive data items at once.
func ClearSensitiveData(items ...interface{}) {
	for _, item := range items {
		switch v := item.(type) {
		case []byte:
			ZeroBytes(v)
		case *string:
			ZeroString(v)
		case *SecureBuffer:
			v.Destroy()
		}
	}
}
