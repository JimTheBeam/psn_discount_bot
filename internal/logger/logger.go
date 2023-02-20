package logger

import (
	"github.com/sirupsen/logrus"
)

const DefaultLogLevel = logrus.InfoLevel

type Logger struct {
	*logrus.Logger
	config  Config
	service string
	version string
}

func New() *Logger {
	logger := Logger{
		Logger: logrus.New(),
	}
	logger.Formatter = new(logrus.JSONFormatter)

	return &logger
}

func (log *Logger) SetLevel(level string) *Logger {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrusLevel = DefaultLogLevel
	}

	log.Logger.SetLevel(logrusLevel)

	return log
}

func (log *Logger) SetService(service string) *Logger {
	log.service = service

	return log
}

func (log *Logger) SetVersion(version string) {
	log.version = version
}

func (log *Logger) NewEntry() *Entry {
	entry := logrus.NewEntry(log.Logger)

	if log.service != "" {
		entry = entry.WithField("service", log.service)
	}

	if log.version != "" {
		entry = entry.WithField("version", log.version)
	}

	return &Entry{Entry: entry, logger: log}
}
