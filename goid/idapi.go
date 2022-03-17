package goid

import (
	"log"
	"runtime"
	"sync"
)

const (
	stackSize = 1024
)

var (
	anchor       = []byte("goroutine ")
	stackBufPool = sync.Pool{
		New: func() any {
			buf := make([]byte, 64)
			return &buf
		},
	}
)

// getGoidByStack parse the current goroutine's id from caller stack.
// This function could be very slow(like 3000us/op), but it's very safe.
func getGoidByStack() (goid int64) {
	bp := stackBufPool.Get().(*[]byte)
	defer stackBufPool.Put(bp)

	b := *bp
	b = b[:runtime.Stack(b, false)]
	goid, _ = findNextGoid(b, 0)
	return
}

// getAllGoidByStack find all goid through stack; WARNING: This function could be very inefficient
func getAllGoidByStack() (goids []int64) {
	count := runtime.NumGoroutine()
	size := count * stackSize // it's ok?
	buf := make([]byte, size)
	n := runtime.Stack(buf, true)
	buf = buf[:n]
	// parse all goids
	goids = make([]int64, 0, count+4)
	for i := 0; i < len(buf); {
		goid, off := findNextGoid(buf, i)
		if goid > 0 {
			goids = append(goids, goid)
		}
		i = off
	}
	return
}

// Find the next goid from `buf[off:]`
func findNextGoid(buf []byte, off int) (goid int64, next int) {
	i := off
	hit := false
	// skip to anchor
	acr := anchor
	for sb := len(buf) - len(acr); i < sb; {
		if buf[i] == acr[0] && buf[i+1] == acr[1] && buf[i+2] == acr[2] && buf[i+3] == acr[3] &&
			buf[i+4] == acr[4] && buf[i+5] == acr[5] && buf[i+6] == acr[6] &&
			buf[i+7] == acr[7] && buf[i+8] == acr[8] && buf[i+9] == acr[9] {
			hit = true
			i += len(acr)
			break
		}
		for ; i < len(buf) && buf[i] != '\n'; i++ {
		}
		i++
	}
	// return if not hit
	if !hit {
		return 0, len(buf)
	}
	// extract goid
	var done bool
	for ; i < len(buf) && !done; i++ {
		switch buf[i] {
		case '0':
			goid *= 10
		case '1':
			goid = goid*10 + 1
		case '2':
			goid = goid*10 + 2
		case '3':
			goid = goid*10 + 3
		case '4':
			goid = goid*10 + 4
		case '5':
			goid = goid*10 + 5
		case '6':
			goid = goid*10 + 6
		case '7':
			goid = goid*10 + 7
		case '8':
			goid = goid*10 + 8
		case '9':
			goid = goid*10 + 9
		case ' ':
			done = true
			break
		default:
			goid = 0
			log.Printf("should never be here, any bug happens\n")
		}
	}
	next = i
	return
}
