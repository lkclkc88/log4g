package log4g

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	H_FORAMAT  = "2006-01-02_15"
	D_FORAMAT  = "2006-01-02"
	M_FORAMAT  = "2006-01-02_15:04"
	S_FORAMAT  = "2006-01-02_15:04:05"
	QUEUE_SIZE = 1024
)

type fileAppender struct {
	level    Level         // 日志级别
	out      *bufio.Writer //输出
	fileName string        // 输出文件名
	//	filePattern string          //备份文件路径
	//	maxBak      int             //最大备份书
	bakLevel int             //备份级别, 1 天,2 小时 默认天
	async    bool            //是否异步
	queue    chan *LogRecord //队列

	lock  sync.RWMutex
	count int
}

func newFileAppender() *fileAppender {
	tmp := fileAppender{bakLevel: 1}
	return &tmp
}

func (c *fileAppender) initConfig(config LoggerAppenderConfig) {
	if "" != config.Level {
		c.level = stringToLevel(config.Level)
		c.async = config.Async
		c.fileName = config.FileName
		if config.BakLevel > 0 && config.BakLevel < 4 {
			c.bakLevel = config.BakLevel
		} else {
			c.bakLevel = 1
		}
		c.initAppender()
	}
}

// 写字符串到文件
func (c *fileAppender) writeString(data string) {
	if c.async {
		c.lock.RLock()
		defer c.lock.RUnlock()
		c.out.WriteString(data)
		if c.count > 15 {
			c.out.Flush()
			c.count = 0
		} else {
			c.count++
		}
		c.out.Flush()
	} else {
		c.lock.RLock()
		defer c.lock.RUnlock()
		c.out.WriteString(data)
		c.out.Flush()
	}
}

//写入一部队列
func (c *fileAppender) write(log *LogRecord) { //写日志
	if log.level >= c.level {
		if c.async {
			c.queue <- log
		} else {
			c.writeString(log.toString())
		}
	}
}

func (c *fileAppender) getLevel() Level { //获取日志级别
	return DEBUG
}

func (c *fileAppender) bakFile(end string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.out.Flush()
	os.Rename(c.fileName, c.fileName+end)
	f, _ := os.Create(c.fileName)
	c.out = bufio.NewWriter(f)

}

func getTimeStr(t time.Time, foramat string, bakLevel int) string {
	t = t.Add(-10 * time.Second)
	return t.Format(foramat)
}

func getNextTime(bakLevel int, t time.Time) time.Time {
	num := 0
	switch bakLevel {
	case 1:
		num = (24-t.Hour())*3600 - (t.Minute()*60 + t.Second())

	case 2:
		num = 3600 - (t.Minute()*60 + t.Second())

	case 3:
		num = 60 - t.Second()
	}

	t = t.Add(time.Duration(num) * time.Second)
	return t
}

//备份的timer
func (c *fileAppender) bakTimer() {
	if c.bakLevel > 3 || c.bakLevel < 1 {
		c.bakLevel = 1
	}
	nTime := getNextTime(c.bakLevel, time.Now())
	for {
		t := nTime.Unix() - time.Now().Unix()
		if t <= 1 {
			switch c.bakLevel {
			case 1:
				end := getTimeStr(nTime, D_FORAMAT, c.bakLevel)
				c.bakFile(end)
				nTime = getNextTime(c.bakLevel, nTime)
				t = 3600 * 24
			case 2:
				end := getTimeStr(nTime, H_FORAMAT, c.bakLevel)
				c.bakFile(end)
				t = 3600
				nTime = getNextTime(c.bakLevel, nTime)

			case 3:
				end := getTimeStr(nTime, M_FORAMAT, c.bakLevel)
				c.bakFile(end)
				t = 60
				nTime = getNextTime(c.bakLevel, nTime)
			}

		}
		time.Sleep(time.Duration(t) * time.Second)
	}

}

func (c *fileAppender) asyncWrite() {
	for {
		lr := <-c.queue
		if nil != lr {
			c.writeString(lr.toString())
		}
	}
}

func (c *fileAppender) initAppender() {
	if checkFile(c.fileName, false) {
		//		f, _ := os.OpenFile(c.fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		f, _ := os.Create(c.fileName)
		c.out = bufio.NewWriter(f)
		//		c.out = f
		if c.async {
			c.queue = make(chan *LogRecord, QUEUE_SIZE)
			go c.asyncWrite()
		}
		go c.bakTimer()
	} else {
		fmt.Println(" init file failed ", c.fileName)
	}
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func checkFile(filePath string, path bool) bool {
	if "" != filePath {
		file, err := os.Stat(filePath)
		if nil == err {
			if path {
				if file.IsDir() {
					return true
				}
			} else if !file.IsDir() {
				return true
			}
			return false
		} else {
			//获取上级路径,然后
			if checkFile(getParentDirectory(filePath), true) {
				if path {
					err = os.MkdirAll(filePath, 0777)
					if err != nil {
						return false
					} else {
						return true
					}
					//				} else {
					//					os.Open(filePath)
				} else {
					os.Create(filePath)
				}
			}

		}
	}
	return true
}
