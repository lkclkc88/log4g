package log4g

import (
	"fmt"
	"time"
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

//日志记录转字符串
func (record *LogRecord) toString() string {
	format := "[%s] [%s]  %s(%d) %s "
	return fmt.Sprintf(format, timeToString(record.date), levelToString(record.level), record.method, record.line, record.content)
}
