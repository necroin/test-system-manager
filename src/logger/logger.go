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
	errorLevel   = 0
	infoLevel    = 1
	debugLevel   = 2
	verboseLevel = 3
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

func Info(message any) {
	if logLevel >= infoLevel {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			log.Printf("INFO: %s\n", message)
		}()
	}
}

func Error(message any) {
	if logLevel >= errorLevel {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			log.Printf("ERROR: %s\n", message)
		}()
	}
}

func Debug(message any) {
	if logLevel >= debugLevel {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			log.Printf("DEBUG: %s\n", message)
		}()
	}
}

func Verbose(message any) {
	if logLevel >= verboseLevel {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			log.Printf("VERBOSE: %s\n", message)
		}()
	}
}
