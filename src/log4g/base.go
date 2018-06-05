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

//全局配置信息
var GlobalConfig *LogConfig

//日志管理
type LoggerManager struct {
	loggerManager map[string]*Logger
	lock          sync.RWMutex
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

//日志记录
type LogRecord struct {
	date     time.Time //时间
	content  string    //内容
	codePath string    //代码路径
	method   string    //代码方法
	line     int       //行数
	level    Level     //日志级别
}

//输出工具接口
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

//构建日志管理
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

//level转字符串
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

//　字符串转level
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

// 时间转字符串
func timeToString(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

//日志记录转字符串
func (record *LogRecord) toString() string {
	format := "[%s] [%s]  %s(%d) %s "
	return fmt.Sprintf(format, timeToString(record.date), levelToString(record.level), record.method, record.line, record.content)
}

//根据文件路径，获取对应的日志配置信息
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

//初始日志，给日志配置输出对象，日志级别
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

//构建内容
func buildContent(args ...interface{}) string {
	tmp := fmt.Sprintln(args...)
	return tmp
}

//构建日志记录
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
