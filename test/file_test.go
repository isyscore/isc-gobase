package test

import (
	"github.com/isyscore/isc-gobase/file"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/time"
	"os"
	"testing"
)

var osInfoFilePre = "/Users/zhouzhenyong/tem/proc/run_info_"

func TestWrite(t *testing.T) {
	dataMap := map[string]interface{}{}
	dataMap["a"] = 12233
	dataMap["b"] = "sdfasd"

	//ioutil.WriteFile("/Users/zhouzhenyong/tem/proc/test.txt", []byte(isc.ObjectToJson(dataMap)), fs.ModeAppend)
	//ioutil.WriteFile("/Users/zhouzhenyong/tem/proc/test.txt", []byte(isc.ObjectToJson(dataMap)), fs.ModeAppend)

	file.AppendFile("/Users/zhouzhenyong/tem/proc/test.txt", "dfs")

	//fd,_:=os.OpenFile("/Users/zhouzhenyong/tem/proc/test.txt",os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	//fd_content:=strings.Join([]string{isc.ObjectToJson(dataMap)},"")
	//buf:=[]byte(fd_content)
	//fd.Write(buf)
	//fd.Close()

	// 判断文件是否创建
	fileName := osInfoFilePre + time.TimeToStringFormat(time.Now(), time.FmtYMd) + ".log"
	fileExist := file.FileExists(fileName)
	if !fileExist {
		_, err := os.Create(fileName)
		if nil != err {
			logger.Error("文件创建异常, %v+", err.Error())
			return
		}
	} else {
		fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		buf := []byte("xxx")
		fd.Write(buf)
		fd.Close()
	}
}
