package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10 MB default max
)

// ReadFile reads a file from the filesystem.
// Returns the file contents as bytes.
func ReadFile(path string) ([]byte, error) {
	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Check if it's a directory
	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file: %s", path)
	}

	// Check file size
	if info.Size() > MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max %d bytes)", info.Size(), MaxFileSize)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// ReadFromStdin reads data from stdin (for pipe support).
// Example: cat .env | envsend
func ReadFromStdin() ([]byte, error) {
	// Check if stdin is a pipe
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat stdin: %w", err)
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("no data piped to stdin")
	}

	// Read from stdin with size limit
	reader := io.LimitReader(os.Stdin, MaxFileSize)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no data received from stdin")
	}

	return data, nil
}

// ReadFileOrStdin reads from a file if path is provided, otherwise from stdin.
func ReadFileOrStdin(path string) ([]byte, error) {
	if path == "" || path == "-" {
		return ReadFromStdin()
	}
	return ReadFile(path)
}

// WriteFile writes data to a file.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// WriteToStdout writes data to stdout.
func WriteToStdout(data []byte) error {
	_, err := os.Stdout.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}
	return nil
}

// IsStdinPiped checks if data is being piped to stdin.
func IsStdinPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// PromptUser prompts the user for input and returns the response.
func PromptUser(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read user input: %w", err)
	}

	// Trim newline
	response = response[:len(response)-1]
	return response, nil
}

// PromptPassword prompts for a password without echoing.
func PromptPassword(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	
	// Note: This requires platform-specific implementation
	// For production, use golang.org/x/term
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	// Trim newline
	password = password[:len(password)-1]
	fmt.Fprintln(os.Stderr) // Print newline after password input
	return password, nil
}

// FileExists checks if a file exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetFileSize returns the size of a file in bytes.
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
