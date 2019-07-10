package log4g

import (
"strings"
)

type Level uint8 //日志级别
//日志级别
const (
	ALL Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

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
		case "all":
			return ALL
		case "off":
			return OFF
		}
	}
	return OFF
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
	case OFF:
		return "OFF"
	case ALL:
		return "ALL"
	}
	return "OFF"
}
