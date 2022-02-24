package file

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return os.IsExist(err)
	}
	return true
}

func DirectoryExists(dirPath string) bool {
	if s, err := os.Stat(dirPath); err != nil {
		return false
	} else {
		return s.IsDir()
	}
}

func MkDirs(path string) bool {
	if !DirectoryExists(path) {
		return os.MkdirAll(path, os.ModePerm) == nil
	} else {
		return false
	}
}

func DeleteDirs(path string) bool {
	return os.RemoveAll(path) == nil
}

func DeleteFile(filePath string) bool {
	return os.Remove(filePath) == nil
}

func ReadFile(filePath string) string {
	var ret = ""
	if b, err := ioutil.ReadFile(filePath); err == nil {
		ret = string(b)
	}
	return ret
}

func ReadFileBytes(filePath string) []byte {
	var ret []byte
	if b, err := ioutil.ReadFile(filePath); err == nil {
		ret = b
	}
	return ret
}

func ReadFileLines(filePath string) []string {
	var ret []string
	if b, err := ioutil.ReadFile(filePath); err == nil {
		ret = strings.Split(string(b), "\n")
	}
	return ret
}

func WriteFile(filePath string, text string) bool {
	return ioutil.WriteFile(filePath, []byte(text), fs.ModePerm) == nil
}

func WriteFileBytes(filePath string, data []byte) bool {
	return ioutil.WriteFile(filePath, data, fs.ModePerm) == nil
}

func AppendFile(filePath string, text string) bool {
	return ioutil.WriteFile(filePath, []byte(text), fs.ModeAppend) == nil
}

func AppendFileBytes(filePath string, data []byte) bool {
	return ioutil.WriteFile(filePath, data, fs.ModeAppend) == nil
}

func CopyFile(srcFilePath string, destFilePath string) bool {
	src, _ := os.Open(srcFilePath)
	defer func(src *os.File) { _ = src.Close() }(src)
	dst, _ := os.OpenFile(destFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer func(dst *os.File) { _ = dst.Close() }(dst)
	_, err := io.Copy(dst, src)
	return err == nil
}

func RenameFile(srcFilePath string, destFilePath string) bool {
	return os.Rename(srcFilePath, destFilePath) == nil
}
