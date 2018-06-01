package log4g

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	path := "/home/lkclkc88/goWorkspace/log4g/config/logConfig.json"
	f, error := os.Open(path)
	if nil == error {
		LoadConfig(f)
	}
	log := GetLogger()
	for {
		log.Debug(" test")
	}
}
