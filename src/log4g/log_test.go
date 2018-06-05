package log4g

import (
	"bufio"
//	"fmt"
	"os"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	//	path := "/home/lkclkc88/goWorkspace/log4g/config/logConfig.json"
	//	f, error := os.Open(path)
	//	if nil == error {
	//		LoadConfig(f)
	//	}
	//	log := GetLogger()
	//	for {
	//		log.Debug(" test")
	//	}

	path := "/home/lkclkc88/Desktop/t"
	f, _ := os.Open(path)
	//	buff := make([]string, 0)
	r := bufio.NewReader(f)
	out,_ := os.Create("/home/lkclkc88/Desktop/t1")
	for {
		//		tmp := make([]byte, 1024)
		//		n, _ := f.Read(tmp)
		tmp, _, _ := r.ReadLine()
		if nil != tmp && len(tmp) > 0 {
			strs := strings.Split(string(tmp), ",")
			//			fmt.Println(strs[0])
			//			buff = append(buff, strs[0])
			out.WriteString("\"" + strs[0] + "\",\n")
		}
		//		if n > 0 {
		//			buff = append(buff, tmp[:n]...)
		//			if n < 1024 {
		//				break
		//			}
		//		} else {
		//			break
		//		}
	}
	defer out.Close()
	//	for _, v := range buff {
	//		fmt.Println(v)
	//	}
	//	fmt.Println(string(buff))

}
