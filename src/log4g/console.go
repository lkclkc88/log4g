package log4g

import (
	"bufio"
	"os"
	"sync"
)

//工作台输出工具
type consoleAppender struct {
	level Level         // 日志级别
	out   *bufio.Writer //输出
	async bool          //是否异步
	lock  sync.Mutex
	ch    chan *LogRecord
}

func newConsoleAppender() *consoleAppender {
	tmp := consoleAppender{level: ALL}
	tmp.out = bufio.NewWriter(os.Stdout)
	return &tmp
}

//初始化配置
func (c *consoleAppender) initConfig(config LoggerAppenderConfig) {
	if "" != config.Level {
		c.level = stringToLevel(config.Level)
		c.async = config.Async
		if c.async {
			c.ch = make(chan *LogRecord, 1024)
			go c.asyncWrite()
		}
	}
}

//写日志
func (c *consoleAppender) write(log *LogRecord) { //写日志
	if log.level >= c.level {
		if c.async {
			c.ch <- log
		} else {
			c.lock.Lock()
			defer c.lock.Unlock()
			c.writeString(log.toString())
		}
	}
}

//异步写数据
func (c *consoleAppender) asyncWrite() {
	for {
		tmp := <-c.ch
		if nil != tmp {
			c.writeString(tmp.toString())
		}
	}
}

//写数据到输出流
func (c *consoleAppender) writeString(data string) { //写日志
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	c.out.WriteString(data)
	c.out.Flush()

}

//工作台默认为debug级别
func (c *consoleAppender) getLevel() Level { //获取日志级别
	return DEBUG
}
