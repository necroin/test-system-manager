package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"tsm/src/config"
)

const (
	errorLevel = iota
	infoLevel
	verboseLevel
	debugLevel
)

var (
	logLevel    = infoLevel
	mutex       sync.Mutex
	logLevelMap = map[string]int{
		"error":   errorLevel,
		"info":    infoLevel,
		"debug":   debugLevel,
		"verbose": verboseLevel,
	}
)

func Configure(config *config.Config) error {
	log.SetOutput(os.Stdout)
	if config.LogPath != "" {
		if err := os.MkdirAll(path.Dir(config.LogPath), os.ModePerm); err != nil {
			return fmt.Errorf("[Logger] [Error] failed create logs directory: %s", err)
		}

		var logsFile *os.File
		logsFile, err := os.OpenFile(config.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			logsFile, err = os.Create(config.LogPath)
			if err != nil {
				return fmt.Errorf("[Logger] [Error] failed create logs file: %s", err)
			}
		}
		log.SetOutput(logsFile)
	}

	configLogLevel, ok := logLevelMap[config.LogLevel]
	if !ok {
		configLogLevel = infoLevel
	}

	logLevel = configLogLevel

	return nil
}

func print(message string) {
	mutex.Lock()
	defer mutex.Unlock()
	log.Println(message)
}

func Error(message string, args ...any) {
	if logLevel >= errorLevel {
		print("ERROR: " + fmt.Sprintf(message, args...))
	}
}

func Info(message string, args ...any) {
	if logLevel >= infoLevel {
		print("INFO: " + fmt.Sprintf(message, args...))
	}
}

func Verbose(message string, args ...any) {
	if logLevel >= verboseLevel {
		print("VERBOSE: " + fmt.Sprintf(message, args...))
	}
}

func Debug(message string, args ...any) {
	if logLevel >= debugLevel {
		print("DEBUG: " + fmt.Sprintf(message, args...))
	}
}
