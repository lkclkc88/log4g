package log4g

import (
	//	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		//		log.Error(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func TestLog(t *testing.T) {
	//	path := GetCurrentDirectory()
	path := "/home/lkclkc88/git/log4g/logConfig.json"
	fmt.Println(path)
	//	Loger.LoadConfiguration(path, "json")
	file, err := os.Open(path)
	if nil == err {
		LoadConfig(file)
		log := GetLogger()
		log.Info("---------init log read config " + path + "--------")
		for i := 0; i < 1000; i++ {
			log.Info("test", i)
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
}
