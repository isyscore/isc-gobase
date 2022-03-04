package websocket

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type connectionValue struct {
	key   []byte
	value any
}

type ConnectionValues []connectionValue

func (r *ConnectionValues) Set(key string, value any) {
	args := *r
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if string(kv.key) == key {
			kv.value = value
			return
		}
	}

	c := cap(args)
	if c > n {
		args = args[:n+1]
		kv := &args[n]
		kv.key = append(kv.key[:0], key...)
		kv.value = value
		*r = args
		return
	}

	kv := connectionValue{}
	kv.key = append(kv.key[:0], key...)
	kv.value = value
	*r = append(args, kv)
}

func (r *ConnectionValues) Get(key string) any {
	args := *r
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if string(kv.key) == key {
			return kv.value
		}
	}
	return nil
}

func (r *ConnectionValues) Reset() {
	*r = (*r)[:0]
}

type UnderlineConnection interface {
	SetWriteDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetReadLimit(limit int64)
	SetPongHandler(h func(appData string) error)
	SetPingHandler(h func(appData string) error)
	WriteControl(messageType int, data []byte, deadline time.Time) error
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	NextWriter(messageType int) (io.WriteCloser, error)
	Close() error
}

type DisconnectFunc func()
type LeaveRoomFunc func(roomName string)
type ErrorFunc func(error)
type NativeMessageFunc func([]byte)
type MessageFunc any
type PingFunc func()
type PongFunc func()

// Connection 接口
type Connection interface {
	Emitter
	Err() error
	ID() string
	Server() *Server
	Write(websocketMessageType int, data []byte) error
	Context() *gin.Context
	OnDisconnect(DisconnectFunc)
	OnError(ErrorFunc)
	OnPing(PingFunc)
	OnPong(PongFunc)
	FireOnError(err error)
	To(string) Emitter
	OnMessage(NativeMessageFunc)
	On(string, MessageFunc)
	Join(string)
	IsJoined(roomName string) bool
	Leave(string) bool
	OnLeave(roomLeaveCb LeaveRoomFunc)
	Wait()
	Disconnect() error
	SetValue(key string, value any)
	GetValue(key string) any
	GetValueArrString(key string) []string
	GetValueString(key string) string
	GetValueInt(key string) int
}

// Connection 实现
type connection struct {
	err                      error
	underline                UnderlineConnection
	id                       string
	messageType              int
	disconnected             bool
	onDisconnectListeners    []DisconnectFunc
	onRoomLeaveListeners     []LeaveRoomFunc
	onErrorListeners         []ErrorFunc
	onPingListeners          []PingFunc
	onPongListeners          []PongFunc
	onNativeMessageListeners []NativeMessageFunc
	onEventListeners         map[string][]MessageFunc
	started                  bool
	self                     Emitter
	broadcast                Emitter
	all                      Emitter
	ctx                      *gin.Context
	values                   ConnectionValues
	server                   *Server
	writerMu                 sync.Mutex
}

var _ Connection = &connection{}

const CloseMessage = websocket.CloseMessage

func newConnection(ctx *gin.Context, s *Server, underlineConn UnderlineConnection, id string) *connection {
	c := &connection{
		underline:                underlineConn,
		id:                       id,
		messageType:              websocket.TextMessage,
		onDisconnectListeners:    make([]DisconnectFunc, 0),
		onRoomLeaveListeners:     make([]LeaveRoomFunc, 0),
		onErrorListeners:         make([]ErrorFunc, 0),
		onNativeMessageListeners: make([]NativeMessageFunc, 0),
		onEventListeners:         make(map[string][]MessageFunc, 0),
		onPongListeners:          make([]PongFunc, 0),
		started:                  false,
		ctx:                      ctx,
		server:                   s,
	}

	if s.config.BinaryMessages {
		c.messageType = websocket.BinaryMessage
	}

	c.self = newEmitter(c, c.id)
	c.broadcast = newEmitter(c, Broadcast)
	c.all = newEmitter(c, All)

	return c
}

func (c *connection) Err() error {
	return c.err
}

func (c *connection) Write(websocketMessageType int, data []byte) error {
	c.writerMu.Lock()
	if writeTimeout := c.server.config.WriteTimeout; writeTimeout > 0 {
		_ = c.underline.SetWriteDeadline(time.Now().Add(writeTimeout))
	}

	err := c.underline.WriteMessage(websocketMessageType, data)
	c.writerMu.Unlock()
	if err != nil {
		_ = c.Disconnect()
	}
	return err
}

func (c *connection) writeDefault(data []byte) {
	_ = c.Write(c.messageType, data)
}

const WriteWait = 1 * time.Second

func (c *connection) startPinger() {
	pingHandler := func(message string) error {
		err := c.underline.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(WriteWait))
		if err == websocket.ErrCloseSent {
			return nil
		} else if _, ok := err.(net.Error); ok {
			return nil
		}
		return err
	}

	c.underline.SetPingHandler(pingHandler)

	go func() {
		for {
			time.Sleep(c.server.config.PingPeriod)
			if c.disconnected {
				break
			}
			c.fireOnPing()
			err := c.Write(websocket.PingMessage, []byte{})
			if err != nil {
				break
			}
		}
	}()
}

func (c *connection) fireOnPing() {
	for i := range c.onPingListeners {
		c.onPingListeners[i]()
	}
}

func (c *connection) fireOnPong() {
	for i := range c.onPongListeners {
		c.onPongListeners[i]()
	}
}

func (c *connection) startReader() {
	conn := c.underline
	hasReadTimeout := c.server.config.ReadTimeout > 0

	conn.SetReadLimit(c.server.config.MaxMessageSize)
	conn.SetPongHandler(func(s string) error {
		if hasReadTimeout {
			_ = conn.SetReadDeadline(time.Now().Add(c.server.config.ReadTimeout))
		}
		go c.fireOnPong()
		return nil
	})

	defer func() { _ = c.Disconnect() }()

	for {
		if hasReadTimeout {
			_ = conn.SetReadDeadline(time.Now().Add(c.server.config.ReadTimeout))
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.FireOnError(err)
			}
			break
		} else {
			c.messageReceived(data)
		}
	}
}

func (c *connection) messageReceived(data []byte) {

	if bytes.HasPrefix(data, c.server.config.EvtMessagePrefix) {
		receivedEvt := c.server.messageSerializer.getWebsocketCustomEvent(data)
		listeners, ok := c.onEventListeners[string(receivedEvt)]
		if !ok || len(listeners) == 0 {
			return
		}

		customMessage, err := c.server.messageSerializer.deserialize(receivedEvt, data)
		if customMessage == nil || err != nil {
			return
		}

		for i := range listeners {
			if fn, ok := listeners[i].(func()); ok {
				fn()
			} else if fnString, ok := listeners[i].(func(string)); ok {

				if msgString, is := customMessage.(string); is {
					fnString(msgString)
				} else if msgInt, is := customMessage.(int); is {
					fnString(strconv.Itoa(msgInt))
				}

			} else if fnInt, ok := listeners[i].(func(int)); ok {
				fnInt(customMessage.(int))
			} else if fnBool, ok := listeners[i].(func(bool)); ok {
				fnBool(customMessage.(bool))
			} else if fnBytes, ok := listeners[i].(func([]byte)); ok {
				fnBytes(customMessage.([]byte))
			} else {
				listeners[i].(func(any))(customMessage)
			}

		}
	} else {
		for i := range c.onNativeMessageListeners {
			c.onNativeMessageListeners[i](data)
		}
	}

}

func (c *connection) ID() string {
	return c.id
}

func (c *connection) Server() *Server {
	return c.server
}

func (c *connection) Context() *gin.Context {
	return c.ctx
}

func (c *connection) Values() ConnectionValues {
	return c.values
}

func (c *connection) fireDisconnect() {
	for i := range c.onDisconnectListeners {
		c.onDisconnectListeners[i]()
	}
}

func (c *connection) OnDisconnect(cb DisconnectFunc) {
	c.onDisconnectListeners = append(c.onDisconnectListeners, cb)
}

func (c *connection) OnError(cb ErrorFunc) {
	c.onErrorListeners = append(c.onErrorListeners, cb)
}

func (c *connection) OnPing(cb PingFunc) {
	c.onPingListeners = append(c.onPingListeners, cb)
}

func (c *connection) OnPong(cb PongFunc) {
	c.onPongListeners = append(c.onPongListeners, cb)
}

func (c *connection) FireOnError(err error) {
	for _, cb := range c.onErrorListeners {
		cb(err)
	}
}

func (c *connection) To(to string) Emitter {
	if to == Broadcast {
		return c.broadcast
	} else if to == All {
		return c.all
	} else if to == c.id {
		return c.self
	}

	return newEmitter(c, to)
}

func (c *connection) EmitMessage(nativeMessage []byte) error {
	return c.self.EmitMessage(nativeMessage)
}

func (c *connection) Emit(event string, message any) error {
	return c.self.Emit(event, message)
}

func (c *connection) OnMessage(cb NativeMessageFunc) {
	c.onNativeMessageListeners = append(c.onNativeMessageListeners, cb)
}

func (c *connection) On(event string, cb MessageFunc) {
	if c.onEventListeners[event] == nil {
		c.onEventListeners[event] = make([]MessageFunc, 0)
	}

	c.onEventListeners[event] = append(c.onEventListeners[event], cb)
}

func (c *connection) Join(roomName string) {
	c.server.Join(roomName, c.id)
}

func (c *connection) IsJoined(roomName string) bool {
	return c.server.IsJoined(roomName, c.id)
}

func (c *connection) Leave(roomName string) bool {
	return c.server.Leave(roomName, c.id)
}

func (c *connection) OnLeave(roomLeaveCb LeaveRoomFunc) {
	c.onRoomLeaveListeners = append(c.onRoomLeaveListeners, roomLeaveCb)
}

func (c *connection) fireOnLeave(roomName string) {
	if c == nil {
		return
	}
	for i := range c.onRoomLeaveListeners {
		c.onRoomLeaveListeners[i](roomName)
	}
}

func (c *connection) Wait() {
	if c.started {
		return
	}
	c.started = true
	c.startPinger()
	c.startReader()
}

var ErrAlreadyDisconnected = errors.New("already disconnected")

func (c *connection) Disconnect() error {
	if c == nil || c.disconnected {
		return ErrAlreadyDisconnected
	}
	return c.server.Disconnect(c.ID())
}

func (c *connection) SetValue(key string, value any) {
	c.values.Set(key, value)
}

func (c *connection) GetValue(key string) any {
	return c.values.Get(key)
}

func (c *connection) GetValueArrString(key string) []string {
	if v := c.values.Get(key); v != nil {
		if arrString, ok := v.([]string); ok {
			return arrString
		}
	}
	return nil
}

func (c *connection) GetValueString(key string) string {
	if v := c.values.Get(key); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (c *connection) GetValueInt(key string) int {
	if v := c.values.Get(key); v != nil {
		if i, ok := v.(int); ok {
			return i
		} else if s, ok := v.(string); ok {
			if iv, err := strconv.Atoi(s); err == nil {
				return iv
			}
		}
	}
	return 0
}
