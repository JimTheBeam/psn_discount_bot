package logger

import "github.com/sirupsen/logrus"

type Entry struct {
	*logrus.Entry
	logger *Logger
}

func (e *Entry) Logger() *Logger {
	return e.logger
}
