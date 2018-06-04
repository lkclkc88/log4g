package log4g

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

//日志管理器
//var loggerManager map[string]*Logger = make(map[string]*Logger)
var loggerManager LoggerManager = newLoggerManager()
var GlobalConfig *LogConfig //全局配置信息

type LoggerManager struct {
	loggerManager map[string]*Logger
	lock          sync.RWMutex
}

func newLoggerManager() LoggerManager {
	tmp := LoggerManager{loggerManager: make(map[string]*Logger)}
	return tmp
}

func (loggerManager *LoggerManager) put(key string, v *Logger) {
	loggerManager.lock.Lock()
	defer loggerManager.lock.Unlock()
	loggerManager.loggerManager[key] = v
}

func (loggerManager *LoggerManager) get(key string) *Logger {
	loggerManager.lock.RLock()
	defer loggerManager.lock.RUnlock()
	return loggerManager.loggerManager[key]
}

func (loggerManager *LoggerManager) initLoggerManager() {
	loggerManager.lock.Lock()
	defer loggerManager.lock.Unlock()
	for _, v := range loggerManager.loggerManager {
		initLogger(v)
	}
}

type Level uint8 //日志级别

//日志级别
const (
	_ Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

func levelToString(level Level) string {

	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}

	return "Default"
}

func stringToLevel(level string) Level {
	if "" != level {
		switch strings.ToLower(level) {
		case "debug":
			return DEBUG
		case "info":
			return INFO

		case "warn":
			return WARN
		case "error":
			return ERROR
		case "FATAL":
			return FATAL

		}
	}
	return FATAL
}

//日志实体
type LogRecord struct {
	date     time.Time //时间
	content  string    //内容
	codePath string    //代码路径
	method   string    //代码方法
	line     int       //行数
	level    Level     //日志级别
}

func timeToString(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func (record *LogRecord) toString() string {
	format := "[%s] [%s]  %s(%d) %s "
	return fmt.Sprintf(format, timeToString(record.date), levelToString(record.level), record.method, record.line, record.content)
}

type Appender interface {
	write(log *LogRecord) //写日志

	getLevel() Level //获取日志级别

	initConfig(config LoggerAppenderConfig) //初始化配置

}

//日志工具结构体
type Logger struct {
	codePath  string               //代码路径
	level     Level                //日志级别
	appenders map[string]*Appender //日志记录工具集合
}

func getLoggerConfigByPath(path string) *LoggerConfig {
	if nil == GlobalConfig {
		return nil
	}
	configMap := make(map[string]LoggerConfig)
	for _, v := range GlobalConfig.Loggers {
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

func initLogger(tmp *Logger) {
	logConfig := getLoggerConfigByPath(tmp.codePath)
	if nil != logConfig {
		var level Level = FATAL
		appenders := make(map[string]*Appender, 0)
		if nil != logConfig {
			level = stringToLevel(logConfig.Level)
			for _, v := range logConfig.Appenders {
				appenders[v] = GlobalConfig.appenders[v]
			}
		}
		tmp.level = level
		tmp.appenders = appenders
	}
}

func GetLogger() *Logger {
	_, file, _, _ := runtime.Caller(1)
	logger := loggerManager.get(file)
	if nil != logger {
		return logger
	} else {
		tmp := Logger{codePath: file}
		initLogger(&tmp)
		//		loggerManager[file] = &tmp
		loggerManager.put(file, &tmp)
		//		logConfig := getLoggerConfigByPath(file)
		//		if nil != logConfig {
		//			var level Level = FATAL
		//			appenders := make(map[string]*Appender, 0)
		//			if nil != logConfig {
		//				level = stringToLevel(logConfig.Level)
		//				for _, v := range logConfig.Appenders {
		//					appenders[v] = GlobalConfig.appenders[v]
		//				}
		//			}
		//			tmp.level = level
		//			tmp.appenders = appenders
		//			loggerManager[file] = &tmp
		//		}
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

func buildStringContent(format string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	} else {
		return format
	}
}

//构建内容
func buildContent(args ...interface{}) string {
	tmp := fmt.Sprintln(args...)
	return tmp
}

func (log *Logger) buildLogRecord(level Level, args ...interface{}) *LogRecord {
	tmp := LogRecord{date: time.Now(), level: level}
	pc, _, lineno, ok := runtime.Caller(2)
	tmp.codePath = log.codePath
	if ok {
		tmp.line = lineno
		tmp.method = runtime.FuncForPC(pc).Name()
	}
	tmp.content = buildContent(args...)
	return &tmp
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
func (log *Logger) Warn(arg interface{}, args ...interface{}) {
	if log.IsWarn() {
		data := log.buildLogRecord(DEBUG, args...)
		log.write(WARN, data)
	}
}
func (log *Logger) Error(arg interface{}, args ...interface{}) {
	if log.IsWarn() {
		data := log.buildLogRecord(DEBUG, args...)
		log.write(ERROR, data)
	}
}

//写数据
func (log *Logger) write(level Level, record *LogRecord) {
	if nil != log.appenders {
		for _, v := range log.appenders {
			t := *v
			if t.getLevel() <= level {
				t.write(record)
			}
		}
	}
}
