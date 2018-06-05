package log4g

import (
	"bufio"
	"os"
)

//工作台输出工具
type consoleAppender struct {
	level Level         // 日志级别
	out   *bufio.Writer //输出
}

func newConsoleAppender() *consoleAppender {
	tmp := consoleAppender{level: DEBUG}
	tmp.out = bufio.NewWriter(os.Stdout)
	return &tmp
}

func (c *consoleAppender) initConfig(config LoggerAppenderConfig) {
	if "" != config.Level {
		c.level = stringToLevel(config.Level)
	}
}

func (c *consoleAppender) write(log *LogRecord) { //写日志
	if log.level >= c.level {
		c.out.WriteString(log.toString())
		c.out.Flush()
	}
}

func (c *consoleAppender) getLevel() Level { //获取日志级别
	return DEBUG
}
