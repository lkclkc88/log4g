package log4g

import (
	"runtime"
	"strings"
)

var log = GetLogger()

//日志管理器
var lm LoggerManager = newLoggerManager()

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

//根据文件路径，获取对应的日志配置信息
func getLoggerConfigByPath(path string) *LoggerConfig {
	if nil == lm.config {
		return nil
	}
	configMap := make(map[string]LoggerConfig)
	for _, v := range lm.config.Loggers {
		name := v.Name
		if "" != name {
			if v.Name == "root" || strings.Contains(path, name) {
				if nil != v.Appenders {
					configMap[v.Name] = v
				}
			}
		}
	}
	size := len(configMap)
	if size > 0 {

		if size > 1 {
			delete(configMap, "root") //当存在多个配置的时候,优先使用局部被指,移除全局配置,只保留路径最长的配置
			key := ""
			for k, _ := range configMap {
				if len(k) > len(key) {
					key = k
				}
			}
			config := configMap[key]
			return &config
		} else {
			//只有一个配置的时候,直接返回
			for _, v := range configMap {
				return &v
			}
		}
	}
	return nil
}
