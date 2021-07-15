package cipherPayload

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func isExist(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func recoveryCatch() {
	if r := recover(); r != nil {
		log.Debug().Msgf("Recovered in f: %v", r)
	}
}

// Logger
type logger interface {
	printf(logType string, anything ...interface{})
}

type loggerConfig struct {
	isShow bool
}

func newLogger(isShow bool) logger {
	return &loggerConfig{
		isShow: isShow,
	}
}

func (l *loggerConfig) printf(logType string, anything ...interface{}) {
	if !l.isShow {
		return
	}

	var logList []string
	for _, val := range anything {
		valStr := fmt.Sprintf("%v", val)
		logList = append(logList, valStr)
	}

	logHeader := "[CipherPayload]"
	logString := strings.Join(logList, " ")
	switch logType {
	case "error":
		log.Error().Msgf("%v %v", logHeader, logString)
	case "info":
		log.Info().Msgf("%v %v", logHeader, logString)
	default:
		log.Debug().Msgf("%v %v", logHeader, logString)
	}
}
