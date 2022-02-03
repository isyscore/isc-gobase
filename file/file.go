package file

import (
	"io/fs"
	"io/ioutil"
	"os"
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

func WriteFile(filePath string, text string) bool {
	return ioutil.WriteFile(filePath, []byte(text), fs.ModePerm) == nil
}

func WriteFileBytes(filePath string, data []byte) bool {
	return ioutil.WriteFile(filePath, data, fs.ModePerm) == nil
}
