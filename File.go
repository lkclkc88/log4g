package log4g

import (
	"bufio"
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
	QUEUE_SIZE      = 1024
	DEFAULT_MAX_BAK = 10
)

//文件输出工具
type fileAppender struct {
	level    Level         // 日志级别
	out      *bufio.Writer //输出
	fileName string        // 输出文件名
	//	filePattern string          //备份文件路径
	MaxBak   int             //最大备份书 默认10
	bakLevel int             //备份级别, 1 天,2 小时 默认天
	async    bool            //是否异步
	queue    chan *LogRecord //队列
	lock     sync.Mutex
	count    int
}

func newFileAppender() *fileAppender {
	tmp := fileAppender{bakLevel: 1, MaxBak: DEFAULT_MAX_BAK}
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
			if c.async {
				c.queue = make(chan *LogRecord, QUEUE_SIZE)
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

	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if c.async {
		c.out.WriteString(data)
		if c.count > 15 {
			c.out.Flush()
			c.count = 0
		} else {
			c.count++
		}
		c.out.Flush()
	} else {
		c.out.WriteString(data)
		c.out.Flush()
	}
}

//写入异步队列
func (c *fileAppender) write(log *LogRecord) { //写日志
	if log.level >= c.level {
		if c.async {
			c.queue <- log
		} else {
			c.lock.Lock()
			defer c.lock.Unlock()
			c.writeString(log.toString())
		}
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

//异步写
func (c *fileAppender) asyncWrite() {
	for {
		lr := <-c.queue
		if nil != lr {
			c.writeString(lr.toString())
		}
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
