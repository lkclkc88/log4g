package log4g

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	H_FORAMAT       = "2006-01-02_15"
	D_FORAMAT       = "2006-01-02"
	M_FORAMAT       = "2006-01-02_15:04"
	S_FORAMAT       = "2006-01-02_15:04:05"
	QUEUE_SIZE      = 8196
	DEFAULT_MAX_BAK = 10
)

//文件输出工具，实现Appender 接口
type fileAppender struct {
	level    Level         // 日志级别
	out      *bufio.Writer //输出
	fileName string        // 输出文件名
	MaxBak   int           //最大备份书 默认10
	bakLevel int           //备份级别, 1 天,2 小时 默认天
	async    bool          //是否异步
	queue    *Queue        //队列,异步写文件时使用
	lock     *sync.Mutex   //锁
	count    int           //计数，用于异步写数据时，刷新缓存区
}

func newFileAppender() *fileAppender {
	tmp := fileAppender{bakLevel: 1, MaxBak: DEFAULT_MAX_BAK, lock: &sync.Mutex{}}
	return &tmp
}

//初始化文件工具
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
		if checkFile(c.fileName, false) {
			f, _ := os.Create(c.fileName)
			c.out = bufio.NewWriter(f)
			//			如果是异步，初始化队列
			if c.async {
				c.queue = NewQueue(QUEUE_SIZE)
				go c.asyncWrite()
			}
			go c.bakTimer()
		} else {
			fmt.Println(" init file failed ", c.fileName)
		}
	}
}

// 写字符串到文件，如果是异步写，写１５次之后，刷新缓冲区
func (c *fileAppender) writeString(data string) {
	//异常处理
	defer recoverErr()
	//同步调用
	c.out.WriteString(data)
	c.out.Flush()
}

//写入异步队列
func (c *fileAppender) write(log *LogRecord) { //写日志
	defer recoverErr()
	if log.level >= c.level {
		if c.async {
			//写入队列
			c.queue.Offer(log, time.Second)
		} else {
			c.lock.Lock()
			defer c.lock.Unlock()
			c.writeString(log.toString())
		}
	}
}

//写日志
func (c *fileAppender) syncWrite(log *LogRecord) { //写日志
	defer recoverErr()
	if log.level >= c.level {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.writeString(log.toString())
	}
}

func (c *fileAppender) getLevel() Level { //获取日志级别
	return DEBUG
}

//备份日志文件，后去需要增加数据压缩
func (c *fileAppender) bakFile(end string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.out.Flush()
	os.Rename(c.fileName, c.fileName+end)
	f, _ := os.Create(c.fileName)
	c.out = bufio.NewWriter(f)

}

//备份数据时，获取时间
func getTimeStr(t time.Time, foramat string, bakLevel int) string {
	t = t.Add(-10 * time.Second)
	return t.Format(foramat)
}

//获取下次一次备份的时间
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

func softFile(fs *[]os.FileInfo) {
	tmp := *fs
	size := len(tmp)
	if size > 2 {
		for i := 1; i < size; i++ {
			for j := i; j > 0; j-- {
				if tmp[j].ModTime().Unix() > tmp[j-1].ModTime().Unix() {
					t := tmp[j]
					tmp[j] = tmp[j-1]
					tmp[j-1] = t
				} else {
					break
				}
			}
		}
	}
}

//清理历史备份
func (c *fileAppender) cleanHistoryBak() {
	parentPath := getParentDirectory(c.fileName)
	fileName := strings.Replace(c.fileName, parentPath+"/", "", 1)
	files, _ := ioutil.ReadDir(parentPath)
	fs := make([]os.FileInfo, 0)
	for _, f := range files {
		if strings.Contains(f.Name(), fileName) && f.Name() != fileName {
			fs = append(fs, f)
		}
	}
	size := len(fs)
	if size > c.MaxBak {
		softFile(&fs)
		for i := 10; i < size; i++ {
			f := fs[i]
			os.Remove(parentPath + "/" + f.Name())
		}
	}

}

//备份的timer
func (c *fileAppender) bakTimer() {
	defer recoverErr()
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
		c.cleanHistoryBak()
		time.Sleep(time.Duration(t) * time.Second)
	}

}

func (c *fileAppender) getDataFromQueueAndWrite(size int32) {
	var buff bytes.Buffer
	for size > 0 {
		size--
		x := c.queue.Get()
		if nil != x {
			lr := x.(*LogRecord)
			buff.WriteString(lr.toString())
		}
		if buff.Len() > 4096 {
			c.writeString(buff.String())
			buff.Reset()
		}
	}
}

func (c *fileAppender) getDataAndWrite() {
	defer recoverErr()
	size := c.queue.Size()
	if size > 0 {
		c.getDataFromQueueAndWrite(size)
	} else {
		x := c.queue.Poll(10 * time.Second)
		if nil != x {
			lr := x.(*LogRecord)
			c.writeString(lr.toString())
		}
	}
}

//异步写
func (c *fileAppender) asyncWrite() {
	for {
		c.getDataAndWrite()
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
				} else {
					os.Create(filePath)
				}
			}

		}
	}
	return true
}
