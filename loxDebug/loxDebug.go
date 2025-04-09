package loxDebug

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logger *log.Logger
var logFile *os.File
var logDir = "logs"

func InitializeLogger() {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("golox_lox_%s.txt", timestamp)

	var err error
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return
	}

	logPath := filepath.Join(logDir, filename)

	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}

	logger = log.New(logFile, "", log.LstdFlags)
	log.Println("Logger initialized, writing to", filename)
}

func LogDebug(format string, v ...interface{}) {
	if logger != nil {
		logger.Printf("[DEBUG] "+format, v...)
	}
}

func LogInfo(format string, v ...interface{}) {
	if logger != nil {
		logger.Printf("[INFO] "+format, v...)
	}
}

func LogError(format string, v ...interface{}) {
	if logger != nil {
		logger.Printf("[ERROR] "+format, v...)
	}
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
