package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrus *logrus.Logger
}

func NewLogger() *Logger {
	log := logrus.New()
	log.Out = os.Stdout

	return &Logger{
		logrus: log,
	}
}

func (l *Logger) Error(err error) {
	l.logrus.Error(err)
}
