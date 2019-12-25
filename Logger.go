package log4g

import (
	"sync"
	"time"
	"runtime"
)

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

//写数据
func (log *Logger) syncWrite(level Level, record *LogRecord) {
	log.lock.RLock()
	defer log.lock.RUnlock()
	if nil != log.appenders {
		for _, v := range log.appenders {
			t := *v
			if t.getLevel() <= level {
				t.syncWrite(record)
			}
		}
	}
}
