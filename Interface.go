package log4g

import (
	"runtime"
)

//获取日志对象
func GetLogger() *Logger {
	_, file, _, _ := runtime.Caller(1)
	logger := loggerManager.get(file)
	if nil != logger {
		return logger
	} else {
		tmp := Logger{codePath: file}
		initLogger(&tmp)
		loggerManager.put(file, &tmp)

		return &tmp
	}
}

//是否为debug
func (log *Logger) IsDebug() bool {
	return log.level <= DEBUG
}

//是否运行info级别日志
func (log *Logger) IsInfo() bool {
	return log.level <= INFO
}
func (log *Logger) IsWarn() bool {
	return log.level <= WARN
}
func (log *Logger) IsError() bool {
	return log.level <= ERROR
}

func (log *Logger) Debug(args ...interface{}) {
	if log.IsDebug() {
		data := log.buildLogRecord(DEBUG, args...)
		log.write(DEBUG, data)
	}
}

func (log *Logger) Info(args ...interface{}) {
	if log.IsInfo() {
		data := log.buildLogRecord(INFO, args...)
		log.write(INFO, data)
	}
}
func (log *Logger) Warn(args ...interface{}) {
	if log.IsWarn() {
		data := log.buildLogRecord(WARN, args...)
		log.write(WARN, data)
	}
}
func (log *Logger) Error(args ...interface{}) {
	if log.IsWarn() {
		data := log.buildLogRecord(ERROR, args...)
		log.write(ERROR, data)
	}
}
