package log4g

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// 日志配置，包含所有配置信息
type LogConfig struct {
	appenders    map[string]*Appender
	Loggers      []LoggerConfig
	AppendersMap map[string]LoggerAppenderConfig
}

// 日志工具配置
type LoggerConfig struct {
	Name      string
	Level     string
	Appenders []string
}

//日志数据工具配置
type LoggerAppenderConfig struct {
	Appender string
	Level    string
	FileName string
	Async    bool
	BakLevel int
}

func (c *LogConfig) initConfig() {
	if nil != c.AppendersMap {
		appenders := make(map[string]*Appender)
		for k, v := range c.AppendersMap {
			switch strings.ToLower(v.Appender) {
			case "file":
				appender := newFileAppender()
				appender.initConfig(v)
				var tmp Appender = appender
				appenders[k] = &tmp
			case "console":
				appender := newConsoleAppender()
				appender.initConfig(v)
				var tmp Appender = appender
				appenders[k] = &tmp

			}
		}
		c.appenders = appenders
	}
}

func LoadConfig(file *os.File) {
	size := 1024
	buff := make([]byte, 0)
	for {
		tmp := make([]byte, size)
		n, err := file.Read(tmp)
		if err != nil {
			fmt.Println(err)
			break
		}
		if n == 0 {
			break
		} else if n < size {
			buff = append(buff, tmp[:n]...)
			break
		}
	}
	str := string(buff)
	str = strings.TrimSpace(str)
	buff = []byte(str)
	size = len(buff)
	if size > 0 {
		logConfig := &LogConfig{}
		err := json.Unmarshal(buff, &logConfig)
		if nil != err {
			fmt.Println(err)
		}
		logConfig.initConfig()
		GlobalConfig = logConfig
		loggerManager.initLoggerManager()

	}
}
