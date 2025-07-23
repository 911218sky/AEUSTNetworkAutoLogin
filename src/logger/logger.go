package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"AEUSTNetworkAutoLogin/src/utils"
)

const (
	// FilePermission defines default file permissions
	FilePermission = 0644
)

var (
	// logMutex ensures thread-safe log writing
	logMutex sync.Mutex
)

// LogSuccess logs successful operations
// username: the username
// message: the log message
// loginLogPath: path to the login log file
func LogSuccess(username, message, loginLogPath string) error {
	logEntry := fmt.Sprintf("[%s] SUCCESS: %s %s\n", utils.GetFormattedDateTime(), username, message)
	fmt.Print(logEntry)
	return appendToLogFile(loginLogPath, logEntry)
}

// LogError logs error messages
// err: the error message
// errorLogPath: path to the error log file
func LogError(err error, errorLogPath string) error {
	logEntry := fmt.Sprintf("[%s] ERROR: %v\n", utils.GetFormattedDateTime(), err)
	fmt.Print(logEntry)
	return appendToLogFile(errorLogPath, logEntry)
}

// appendToLogFile appends a log message to the specified log file
// Uses mutex to ensure thread safety when multiple goroutines write simultaneously
func appendToLogFile(logFilePath, logEntry string) error {
	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, FilePermission)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}

	return nil
}

// ClearLogFile clears the specified log file
func ClearLogFile(logFilePath string) error {
	logMutex.Lock()
	defer logMutex.Unlock()

	if err := os.WriteFile(logFilePath, []byte{}, FilePermission); err != nil {
		return fmt.Errorf("failed to clear log file: %w", err)
	}
	return nil
}

// RotateLogFile rotates the log file
// When the log file exceeds the specified size, it will be renamed as a backup file
func RotateLogFile(logFilePath string, maxSizeMB int64) error {
	logMutex.Lock()
	defer logMutex.Unlock()

	info, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get log file info: %w", err)
	}

	if info.Size() < maxSizeMB*1024*1024 {
		return nil
	}

	backupPath := fmt.Sprintf("%s.%s.backup", logFilePath, utils.GetFormattedDateTime())
	if err := os.Rename(logFilePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	return nil
}
