package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/baselog"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func TestInfo(t *testing.T) {
	baselog.Info("hello1 %v", "isyscore")
	getCaller()
	fmt.Println("datda")
	//fmt.Println(isc.ToJsonString(rum))
}

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

var callerInitOnce sync.Once
var minimumCallerDepth = 0
var logrusPackage string

func getCaller()  {
	fmt.Println("========1=======")
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			//if strings.Contains(funcName, "getCaller") {
			//	logrusPackage = getPackageName(funcName)
			//	break
			//}

			fmt.Println(funcName)
		}

		minimumCallerDepth = knownLogrusFrames
	})
	fmt.Println("=======2========")

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	fmt.Println("=======3========")
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		fmt.Println(pkg)
		fmt.Println(f.Line)
	}
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}
