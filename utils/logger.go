package utils

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
}

func Log(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)
	logFile.WriteString(logMessage)
	fmt.Print(logMessage)
}

func LogError(err error) {
	if err != nil {
		Log("ERROR: " + err.Error())
	}
}

func LogInfo(message string) {
	Log("INFO: " + message)
}

func LogWarning(message string) {
	Log("WARNING: " + message)
}
