package imLog

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

//日志管理，存储当前应用使用的日志信息，混粗日志构建工具
type LoggerManager struct {
	loggerManager map[string]*Logger
	config        *LogConfig
	lock          *sync.RWMutex
}

//构建日志管理
func newLoggerManager() LoggerManager {
	tmp := LoggerManager{loggerManager: make(map[string]*Logger), lock: &sync.RWMutex{}}
	return tmp
}

// 根据key存放日志对象,如果已经存在，返回存在的数据
func (loggerManager *LoggerManager) put(key string, v *Logger) *Logger {
	loggerManager.lock.Lock()
	defer loggerManager.lock.Unlock()
	old := loggerManager.loggerManager[key]
	if nil == old {
		loggerManager.loggerManager[key] = v
		return v
	} else {
		return old
	}

}

//获取日志对象
func (loggerManager *LoggerManager) get(key string) *Logger {
	loggerManager.lock.RLock()
	defer loggerManager.lock.RUnlock()
	return loggerManager.loggerManager[key]
}

//初始配置。获取日志对象时，可能没有初始配置，加入配置文件之后，需要重新初始化配置
func (lm *LoggerManager) initLoggerManager(config *LogConfig) {
	lm.lock.Lock()
	defer lm.lock.Unlock()
	lm.config=config
	for _, v := range lm.loggerManager {
		//		initLogger(v)
		v.initLogger()
	}
}

//初始日志，给日志配置输出对象，日志级别，第一次获取日志工具结构体时，可能没有初始化完成。
func (tmp *Logger) initLogger() {
	logConfig := getLoggerConfigByPath(tmp.codePath)
	if nil != logConfig {
		var level Level = OFF
		appenders := make(map[string]*Appender, 0)
		if nil != logConfig {
			level = stringToLevel(logConfig.Level)
			for _, v := range logConfig.Appenders {
				appenders[v] = lm.config.appenders[v]
			}
		}
		tmp.lock.Lock()
		defer tmp.lock.Unlock()
		tmp.level = level
		tmp.appenders = appenders
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

//日志工具结构体
type Logger struct {
	codePath  string               //代码路径
	level     Level                //日志级别
	appenders map[string]*Appender //日志记录工具集合
	lock      *sync.RWMutex        //加锁
}

//新建日志工具结构体
func newLogger(path string) *Logger {
	tmp := &Logger{codePath: path, lock: &sync.RWMutex{}}
	tmp.initLogger()
	return tmp
}

//日志记录
type LogRecord struct {
	date     time.Time //时间
	content  string    //内容
	codePath string    //代码路径
	method   string    //代码方法
	line     int       //行数
	level    Level     //日志级别
}

//日志记录转字符串
func (record *LogRecord) toString() string {
	format := "[%s] [%s]  %s(%d) %s "
	return fmt.Sprintf(format, timeToString(record.date), levelToString(record.level), record.method, record.line, record.content)
}

// 时间转字符串
func timeToString(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
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
	log.lock.RLock()
	defer log.lock.RUnlock()
	if nil != log.appenders {
		for _, v := range log.appenders {
			t := *v
			if t.getLevel() <= level {
				t.write(record)
			}
		}
	}
}


//输出接口，对应工作台，文件等输出
type Appender interface {
	/**
	写日志
	*/
	write(log *LogRecord)
	/**
	  获取日志级别
	*/
	getLevel() Level
	/**
	  初始化配置
	*/
	initConfig(config LoggerAppenderConfig)
}

