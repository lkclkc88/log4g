package log4g

import ()

//输出接口，对应工作台，文件等输出
type Appender interface {
	/**
	写日志
	*/
	write(log *LogRecord)
	/**
	同步写日志
	*/
	syncWrite(log *LogRecord)
	/**
	  获取日志级别
	*/
	getLevel() Level
	/**
	  初始化配置
	*/
	initConfig(config LoggerAppenderConfig)
}
