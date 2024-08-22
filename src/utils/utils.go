package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PromptUserInput displays a question to the user and captures their input.
// The `question` parameter is the prompt displayed to the user.
// It returns the user's input as a string.
func PromptUserInput(question string) string {
	fmt.Print(question)
	var input string
	fmt.Scanln(&input)
	return input
}

// CreateLogFile creates a log file in the directory of the executable.
// The `filename` parameter is the name of the log file to be created.
// It returns the full path to the created log file and any error encountered.
func CreateLogFile(filename string) (string, error) {
	var filePath string

	if !filepath.IsAbs(filename) {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		filePath = filepath.Join(dir, filename)
	} else {
		filePath = filename
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// GetFormattedDateTime returns the current date and time in a formatted string.
// The format is "YYYY-MM-DD HH:MM:SS".
func GetFormattedDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// CheckFileExists checks if a file exists at the given filepath.
// The `filepath` parameter is the path to the file to check.
// It returns true if the file exists, or false otherwise.
func CheckFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// WriteCustomBinaryFile writes a custom formatted binary file with given content.
// It includes metadata such as file version and content length.
func WriteCustomBinaryFile(filename string, content string) error {
	data := []byte(content)

	var buf bytes.Buffer

	version := uint32(1)
	if err := binary.Write(&buf, binary.LittleEndian, version); err != nil {
		return fmt.Errorf("failed to write version: %w", err)
	}

	dataLength := uint32(len(data))
	if err := binary.Write(&buf, binary.LittleEndian, dataLength); err != nil {
		return fmt.Errorf("failed to write data length: %w", err)
	}

	if _, err := buf.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// ReadCustomBinaryFile reads a custom formatted binary file and returns its content as a string.
// It parses the metadata such as file version and data length before reading the actual content.
func ReadCustomBinaryFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	buf := bytes.NewReader(data)

	var version uint32
	if err := binary.Read(buf, binary.LittleEndian, &version); err != nil {
		return "", fmt.Errorf("failed to read version: %w", err)
	}

	var dataLength uint32
	if err := binary.Read(buf, binary.LittleEndian, &dataLength); err != nil {
		return "", fmt.Errorf("failed to read data length: %w", err)
	}

	readData := make([]byte, dataLength)
	if _, err := buf.Read(readData); err != nil {
		return "", fmt.Errorf("failed to read data: %w", err)
	}

	return string(readData), nil
}
