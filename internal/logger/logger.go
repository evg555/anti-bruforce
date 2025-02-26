package logger

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	timestampFormat = "2006-01-02 15:04:05"

	jsonFormat = "json"
	textFormat = "text"
)

var ErrFormatNotExist = errors.New("format not exist")

type Logger struct {
	logger *log.Logger
}

func New(level, format string) Logger {
	logger := Logger{
		logger: log.New(),
	}

	err := logger.setFormat(format)
	if err != nil {
		panic(fmt.Sprintf("init logger error: %v: %s", err, format))
	}

	err = logger.setLevel(level)
	if err != nil {
		panic(fmt.Sprintf("init logger error: %v", err))
	}

	return logger
}

func (l *Logger) setLevel(level string) error {
	loglevel, err := log.ParseLevel(level)
	if err != nil {
		return err
	}

	l.logger.SetLevel(loglevel)
	return nil
}

func (l *Logger) setFormat(format string) error {
	switch format {
	case jsonFormat:
		l.logger.SetFormatter(&log.JSONFormatter{
			TimestampFormat: timestampFormat,
		})
	case textFormat:
		l.logger.SetFormatter(&log.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: timestampFormat,
		})
	default:
		return ErrFormatNotExist
	}

	return nil
}

func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}
