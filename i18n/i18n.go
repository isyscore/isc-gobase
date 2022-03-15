package i18n

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	f0 "github.com/isyscore/isc-gobase/file"
)

// InitI18N 初始化国际化，传入的语言为默认语言
func InitI18N(language string) error {
	pwd, _ := os.Getwd()
	i18nDir := filepath.Join(pwd, "i18n")
	if !f0.DirectoryExists(i18nDir) {
		f0.MkDirs(i18nDir)
	}
	lngFile := filepath.Join(i18nDir, fmt.Sprintf("%s.po", language))
	if !f0.FileExists(lngFile) {
		return fmt.Errorf("没有找到指定的语言文件 %s", lngFile)
	}
	innerMap = NewI18NMap(language, lngFile)
	return nil
}

// LoadLanguage 加载指定语言，该语言将优先于默认语言，但是若指定语言中不存在某个key，将会从默认语言读取
func LoadI18NLanguage(language string) error {
	pwd, _ := os.Getwd()
	i18nDir := filepath.Join(pwd, "i18n")
	lngFile := filepath.Join(i18nDir, fmt.Sprintf("%s.po", language))
	if !f0.FileExists(lngFile) {
		return fmt.Errorf("没有找到指定的语言文件 %s", lngFile)
	}
	innerMap.Language = language
	innerMap.Data = loadPo(lngFile)
	return nil
}

func T(key string) string {
	if v, ok := innerMap.Data[key]; ok {
		return v
	} else if v, ok := innerMap.DefaultData[key]; ok {
		return v
	} else {
		log.Printf("没有找到key %s", key)
		return ""
	}
}

func Tf(key string, value ...any) string {
	if v, ok := innerMap.Data[key]; ok {
		return fmt.Sprintf(v, value...)
	} else if v, ok := innerMap.DefaultData[key]; ok {
		return fmt.Sprintf(v, value...)
	} else {
		log.Printf("没有找到key %s", key)
		return ""
	}
}
