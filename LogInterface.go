package imLog

import (
	"bufio"
	"os"
	"runtime"
	"runtime/debug"
)

var log = GetLogger()

//异常处理
func recoverErr() {
	if err := recover(); err != nil {
		write := bufio.NewWriter(os.Stdout)
		buff := debug.Stack()
		if nil != buff {
			write.Write(debug.Stack())
		}
	}
}

//获取日志对象
func GetLogger() *Logger {
	_, path, _, _ := runtime.Caller(1)
	logger := lm.get(path)
	if nil != logger {
		return logger
	} else {
		tmp := newLogger(path)
		lm.put(path, tmp)
		return tmp
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

//是否运行warn级别日志
func (log *Logger) IsWarn() bool {
	return log.level <= WARN
}

//是否运行error级别日志
func (log *Logger) IsError() bool {
	return log.level <= ERROR
}

//写入debug日志
func (log *Logger) Debug(args ...interface{}) {
	if log.IsDebug() {
		data := log.buildLogRecord(DEBUG, args...)
		log.write(DEBUG, data)
	}
}

//写入info日志
func (log *Logger) Info(args ...interface{}) {
	if log.IsInfo() {
		data := log.buildLogRecord(INFO, args...)
		log.write(INFO, data)
	}
}

//写入warn日志
func (log *Logger) Warn(args ...interface{}) {
	if log.IsWarn() {
		data := log.buildLogRecord(WARN, args...)
		log.write(WARN, data)
	}
}

//写入error日志
func (log *Logger) Error(args ...interface{}) {
	if log.IsWarn() {
		data := log.buildLogRecord(ERROR, args...)
		log.write(ERROR, data)
	}
}

//日志管理器
var lm LoggerManager = newLoggerManager()
