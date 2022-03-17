package goid

import (
	"net"
	"os"
	"sync/atomic"

	"github.com/isyscore/isc-gobase/time"
)

var seq uint64 = 0
var digits = []uint8{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

const max = 8000

func GenerateTraceID() string {
	buffer := make([]byte, 16)

	// 计算当前session的咋一序号
	atomic.AddUint64(&seq, 1)
	current := seq
	var next uint64
	if current >= max {
		next = 1
	} else {
		next = current + 1
	}
	seq = next
	bs := shortToBytes(uint16(current))
	putBuffer(&buffer, bs, 0)

	// 计算时间
	t0 := time.TimeInMillis()
	bt0 := int64ToBytes(t0)
	putBuffer(&buffer, bt0, 2)

	// 计算IP地址
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.To4()
				putBuffer(&buffer, ip, 10)
				break
			}
		}
	}

	// 计算PID
	pid := os.Getpid()
	bp := shortToBytes(uint16(pid))
	putBuffer(&buffer, bp, 14)

	hex := encodeHex(buffer, digits)

	return string(hex)
}

func putBuffer(buf *[]byte, b []byte, from int) {
	idx := from
	for _, e := range b {
		(*buf)[idx] = e
		idx++
	}
}

func shortToBytes(s uint16) []byte {
	b := make([]byte, 2)
	b[0] = byte(s >> 8)
	b[1] = byte(s)
	return b
}

func int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	b[0] = byte(i >> 56)
	b[1] = byte(i >> 48)
	b[2] = byte(i >> 40)
	b[3] = byte(i >> 32)
	b[4] = byte(i >> 24)
	b[5] = byte(i >> 16)
	b[6] = byte(i >> 8)
	b[7] = byte(i)
	return b
}

func encodeHex(data []byte, dig []uint8) []uint8 {
	l := len(data)
	out := make([]uint8, l<<1)
	var j int = 0
	for i := 0; i < l; i++ {
		out[j] = dig[(0xf0&data[i])>>4]
		j++
		out[j] = dig[0x0f&data[i]]
		j++
	}
	return out
}
