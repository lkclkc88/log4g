package log4g

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

// 时间转字符串
func timeToString(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

//构建内容
func buildContent(args ...interface{}) string {
	tmp := fmt.Sprintln(args...)
	return tmp
}

//异常处理
func recoverErr() {
	if err := recover(); err != nil {
		write := bufio.NewWriter(os.Stdout)
		write.Write(debug.Stack())
		log.Error(debug.Stack())
	}
}
