package debug

import (
	log "github.com/juju/loggo"
)

func init() {
	Init("<root>=ERROR;mal=WARNING;mal.transport=WARNING;isis=WARNING;mal.api=WARNING")
}

type Logger interface {
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	Fatalf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Debugf(message string, args ...interface{})
}

func Init(args string) {
	log.ConfigureLoggers(args)
}

func GetLogger(name string) Logger {
	logger := log.GetLogger(name)
	return &LoggerX{logger: logger}
}

type LoggerX struct {
	logger log.Logger
}

func (l *LoggerX) IsDebugEnabled() bool {
	return l.logger.IsDebugEnabled()
}

func (l *LoggerX) IsInfoEnabled() bool {
	return l.logger.IsInfoEnabled()
}

func (l *LoggerX) IsWarnEnabled() bool {
	return l.logger.IsWarningEnabled()
}

func (l *LoggerX) Fatalf(message string, args ...interface{}) {
	l.logger.Logf(log.CRITICAL, message, args...)

}

func (l *LoggerX) Errorf(message string, args ...interface{}) {
	l.logger.Logf(log.ERROR, message, args...)
}

func (l *LoggerX) Warnf(message string, args ...interface{}) {
	l.logger.Logf(log.WARNING, message, args...)
}

func (l *LoggerX) Infof(message string, args ...interface{}) {
	l.logger.Logf(log.INFO, message, args...)
}

func (l *LoggerX) Debugf(message string, args ...interface{}) {
	l.logger.Logf(log.DEBUG, message, args...)
}
