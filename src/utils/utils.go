package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// FilePermission defines default file permissions
	FilePermission = 0644
	// BinaryFileVersion defines the binary file format version
	BinaryFileVersion = uint32(1)
)

// PromptUserInput displays a question to the user and captures their input
// Supports multi-line input until the user presses Enter
func PromptUserInput(question string) string {
	fmt.Print(question)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// CreateLogFile creates a log file at the specified path
// If the directory doesn't exist, it will be created automatically
func CreateLogFile(filename string) (string, error) {
	filePath := filename
	if !filepath.IsAbs(filename) {
		dir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
		filePath = filepath.Join(dir, filename)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, FilePermission)
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}
	defer file.Close()

	return filePath, nil
}

// GetFormattedDateTime returns the current date and time in formatted string
// Format: "YYYY-MM-DD HH:MM:SS"
func GetFormattedDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// CheckFileExists checks if a file exists at the specified path
func CheckFileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// WriteBinaryFile writes content to a custom format binary file
// Includes metadata such as version and content length
func WriteBinaryFile(filename string, content string) error {
	var buf bytes.Buffer

	// Write version information
	if err := binary.Write(&buf, binary.LittleEndian, BinaryFileVersion); err != nil {
		return fmt.Errorf("failed to write version info: %w", err)
	}

	// Write data length
	data := []byte(content)
	dataLength := uint32(len(data))
	if err := binary.Write(&buf, binary.LittleEndian, dataLength); err != nil {
		return fmt.Errorf("failed to write data length: %w", err)
	}

	// Write actual data
	if _, err := buf.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, buf.Bytes(), FilePermission); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ReadBinaryFile reads a custom format binary file and returns its content
// Parses metadata such as version and data length
func ReadBinaryFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	buf := bytes.NewReader(data)

	// Read version information
	var version uint32
	if err := binary.Read(buf, binary.LittleEndian, &version); err != nil {
		return "", fmt.Errorf("failed to read version info: %w", err)
	}

	// Validate version
	if version != BinaryFileVersion {
		return "", fmt.Errorf("unsupported file version: %d", version)
	}

	// Read data length
	var dataLength uint32
	if err := binary.Read(buf, binary.LittleEndian, &dataLength); err != nil {
		return "", fmt.Errorf("failed to read data length: %w", err)
	}

	// Read actual data
	readData := make([]byte, dataLength)
	if _, err := buf.Read(readData); err != nil {
		return "", fmt.Errorf("failed to read data: %w", err)
	}

	return string(readData), nil
}

// EnsureDirectoryExists ensures that the directory at the specified path exists
func EnsureDirectoryExists(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}
