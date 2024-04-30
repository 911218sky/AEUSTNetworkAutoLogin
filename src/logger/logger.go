package logger

import (
	"fmt"
	"os"

	"AEUSTNetworkAutoLogin/src/utils"
)

func LogSuccess(username, message, loginLogPath string) {
	logEntry := fmt.Sprintf("[%s] : %s %s\n", utils.GetFormattedDateTime(), username, message)
	fmt.Print(logEntry)
	appendToLogFile(loginLogPath, logEntry)
}

func LogError(err error, errorLogPath string) {
	logEntry := fmt.Sprintf("[%s] ERROR: %v\n", utils.GetFormattedDateTime(), err)
	fmt.Print(logEntry)
	appendToLogFile(errorLogPath, logEntry)
}

func appendToLogFile(logFilePath, logEntry string) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(logEntry)
	if err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}
