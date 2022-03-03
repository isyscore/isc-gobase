package websocket

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
)

const (
	DefaultWebsocketWriteTimeout     = 0
	DefaultWebsocketReadTimeout      = 0
	DefaultWebsocketPongTimeout      = 60 * time.Second
	DefaultWebsocketPingPeriod       = (DefaultWebsocketPongTimeout * 9) / 10
	DefaultWebsocketMaxMessageSize   = 1024
	DefaultWebsocketReadBufferSize   = 4096
	DefaultWebsocketWriterBufferSize = 4096
	DefaultEvtMessageKey             = "gin-websocket-message:"

	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var (
	DefaultIDGenerator = func(*gin.Context) string {
		id, err := uuid.NewV4()
		if err != nil {
			return randomString(64)
		}
		return id.String()
	}

	src = rand.NewSource(time.Now().UnixNano())
)

type Config struct {
	IDGenerator       func(ctx *gin.Context) string
	EvtMessagePrefix  []byte
	Error             func(w http.ResponseWriter, r *http.Request, status int, reason error)
	CheckOrigin       func(r *http.Request) bool
	HandshakeTimeout  time.Duration
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	PongTimeout       time.Duration
	PingPeriod        time.Duration
	MaxMessageSize    int64
	BinaryMessages    bool
	ReadBufferSize    int
	WriteBufferSize   int
	EnableCompression bool
	Subprotocols      []string
}

func (c Config) Validate() Config {
	if c.WriteTimeout < 0 {
		c.WriteTimeout = DefaultWebsocketWriteTimeout
	}
	if c.ReadTimeout < 0 {
		c.ReadTimeout = DefaultWebsocketReadTimeout
	}
	if c.PongTimeout < 0 {
		c.PongTimeout = DefaultWebsocketPongTimeout
	}
	if c.PingPeriod <= 0 {
		c.PingPeriod = DefaultWebsocketPingPeriod
	}
	if c.MaxMessageSize <= 0 {
		c.MaxMessageSize = DefaultWebsocketMaxMessageSize
	}
	if c.ReadBufferSize <= 0 {
		c.ReadBufferSize = DefaultWebsocketReadBufferSize
	}
	if c.WriteBufferSize <= 0 {
		c.WriteBufferSize = DefaultWebsocketWriterBufferSize
	}
	if c.Error == nil {
		c.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			//empty
		}
	}
	if c.CheckOrigin == nil {
		c.CheckOrigin = func(r *http.Request) bool {
			// allow all connections by default
			return true
		}
	}
	if len(c.EvtMessagePrefix) == 0 {
		c.EvtMessagePrefix = []byte(DefaultEvtMessageKey)
	}
	if c.IDGenerator == nil {
		c.IDGenerator = DefaultIDGenerator
	}
	return c
}

func random(n int) []byte {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func randomString(n int) string {
	return string(random(n))
}
