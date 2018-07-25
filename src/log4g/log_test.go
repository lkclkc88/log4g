package log4g

import (
	//	"bufio"
	"fmt"
	"os"
	//	"path/filepath"
	//	"strings"
	"testing"
	//	"time"
)

func TestLog(t *testing.T) {
	path := "/lkclkc88/git/log4g/logConfig.json"
	fmt.Println(path)
	file, err := os.Open(path)
	if nil == err {
		LoadConfig(file)
		log := GetLogger()
		log.Info("---------init log read config " + path + "--------")
		for i := 0; i < 150000; i++ {
			//			log.Info("test", i)
			log.Warn("test", i)
		}
		fmt.Println(log.IsDebug())
		fmt.Println(log.IsInfo())
		fmt.Println(log.IsWarn())

		fmt.Println(log.IsError())
		log.Warn(" execute complate")
	} else {
		fmt.Println("init log file failed")
	}
	//	time.Sleep(35 * time.Second)
}
