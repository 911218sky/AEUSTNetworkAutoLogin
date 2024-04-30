package utils

import (
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
