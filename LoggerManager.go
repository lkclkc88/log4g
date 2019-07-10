package log4g

import (
	"sync"
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
	lm.config = config
	for _, v := range lm.loggerManager {
		//		initLogger(v)
		v.initLogger()
	}
}
