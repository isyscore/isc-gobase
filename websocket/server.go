package websocket

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	w0 "github.com/gorilla/websocket"
)

var ClientSource []byte

type ConnectionFunc func(Connection)

type websocketRoomPayload struct {
	roomName     string
	connectionID string
}

type websocketMessagePayload struct {
	from string
	to   string
	data []byte
}

type Server struct {
	config                Config
	ClientSource          []byte
	messageSerializer     *messageSerializer
	connections           sync.Map
	rooms                 map[string][]string
	mu                    sync.RWMutex
	onConnectionListeners []ConnectionFunc
	upgrader              w0.Upgrader
}

func NewWSServer(cfg Config) *Server {
	cfg = cfg.Validate()
	return &Server{
		config:                cfg,
		ClientSource:          bytes.Replace(ClientSource, []byte(DefaultEvtMessageKey), cfg.EvtMessagePrefix, -1),
		messageSerializer:     newMessageSerializer(cfg.EvtMessagePrefix),
		connections:           sync.Map{}, // ready-to-use, this is not necessary.
		rooms:                 make(map[string][]string),
		onConnectionListeners: make([]ConnectionFunc, 0),
		upgrader: w0.Upgrader{
			HandshakeTimeout:  cfg.HandshakeTimeout,
			ReadBufferSize:    cfg.ReadBufferSize,
			WriteBufferSize:   cfg.WriteBufferSize,
			Error:             cfg.Error,
			CheckOrigin:       cfg.CheckOrigin,
			Subprotocols:      cfg.Subprotocols,
			EnableCompression: cfg.EnableCompression,
		},
	}
}

func (s *Server) Handler() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		c := s.Upgrade(ctx)
		if c.Err() != nil {
			return
		}
		for i := range s.onConnectionListeners {
			s.onConnectionListeners[i](c)
		}
		c.Wait()
	}
}

func (s *Server) Upgrade(ctx *gin.Context) Connection {
	conn, err := s.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Printf("websocket error: %v\n", err)
		ctx.AbortWithStatus(503)
		return &connection{err: err}
	}

	return s.handleConnection(ctx, conn)
}

func (s *Server) addConnection(c *connection) {
	s.connections.Store(c.id, c)
}

func (s *Server) getConnection(connID string) (*connection, bool) {
	if cValue, ok := s.connections.Load(connID); ok {
		if conn, ok := cValue.(*connection); ok {
			return conn, ok
		}
	}

	return nil, false
}

func (s *Server) handleConnection(ctx *gin.Context, websocketConn UnderlineConnection) *connection {
	cid := s.config.IDGenerator(ctx)
	c := newConnection(ctx, s, websocketConn, cid)
	s.addConnection(c)
	s.Join(c.id, c.id)
	return c
}

func (s *Server) OnConnection(cb ConnectionFunc) {
	s.onConnectionListeners = append(s.onConnectionListeners, cb)
}

func (s *Server) IsConnected(connID string) bool {
	_, found := s.getConnection(connID)
	return found
}

func (s *Server) Join(roomName string, connID string) {
	s.mu.Lock()
	s.join(roomName, connID)
	s.mu.Unlock()
}

func (s *Server) join(roomName string, connID string) {
	if s.rooms[roomName] == nil {
		s.rooms[roomName] = make([]string, 0)
	}
	s.rooms[roomName] = append(s.rooms[roomName], connID)
}

func (s *Server) IsJoined(roomName string, connID string) bool {
	s.mu.RLock()
	room := s.rooms[roomName]
	s.mu.RUnlock()

	if room == nil {
		return false
	}

	for _, connid := range room {
		if connID == connid {
			return true
		}
	}

	return false
}

func (s *Server) LeaveAll(connID string) {
	s.mu.Lock()
	for name := range s.rooms {
		s.leave(name, connID)
	}
	s.mu.Unlock()
}

func (s *Server) Leave(roomName string, connID string) bool {
	s.mu.Lock()
	left := s.leave(roomName, connID)
	s.mu.Unlock()
	return left
}

func (s *Server) leave(roomName string, connID string) (left bool) {
	if s.rooms[roomName] != nil {
		for i := range s.rooms[roomName] {
			if s.rooms[roomName][i] == connID {
				s.rooms[roomName] = append(s.rooms[roomName][:i], s.rooms[roomName][i+1:]...)
				left = true
				break
			}
		}
		if len(s.rooms[roomName]) == 0 {
			delete(s.rooms, roomName)
		}
	}

	if left {
		if c, ok := s.getConnection(connID); ok {
			c.fireOnLeave(roomName)
		}
	}
	return
}

func (s *Server) GetTotalConnections() (n int) {
	s.connections.Range(func(k, v interface{}) bool {
		n++
		return true
	})

	return n
}

func (s *Server) GetConnections() []Connection {
	length := s.GetTotalConnections()
	conns := make([]Connection, length, length)
	i := 0
	s.connections.Range(func(k, v interface{}) bool {
		conn, ok := v.(*connection)
		if !ok {
			return false
		}
		conns[i] = conn
		i++
		return true
	})

	return conns
}

func (s *Server) GetConnection(connID string) Connection {
	conn, ok := s.getConnection(connID)
	if !ok {
		return nil
	}

	return conn
}

func (s *Server) GetConnectionsByRoom(roomName string) []Connection {
	var conns []Connection
	s.mu.RLock()
	if connIDs, found := s.rooms[roomName]; found {
		for _, connID := range connIDs {
			if cValue, ok := s.connections.Load(connID); ok {
				if conn, ok := cValue.(*connection); ok {
					conns = append(conns, conn)
				}
			}
		}
	}

	s.mu.RUnlock()

	return conns
}

func (s *Server) emitMessage(from, to string, data []byte) {
	if to != All && to != Broadcast {
		s.mu.RLock()
		room := s.rooms[to]
		s.mu.RUnlock()
		if room != nil {
			for _, connectionIDInsideRoom := range room {
				if c, ok := s.getConnection(connectionIDInsideRoom); ok {
					c.writeDefault(data)
				} else {
					cid := connectionIDInsideRoom
					if c != nil {
						cid = c.id
					}
					s.Leave(cid, to)
				}
			}
		}
	} else {
		s.connections.Range(func(k, v interface{}) bool {
			connID, ok := k.(string)
			if !ok {
				return true
			}

			if to != All && to != connID {
				if to == Broadcast && from == connID {
					return true
				}

			}

			conn, ok := v.(*connection)
			if ok {
				conn.writeDefault(data)
			}

			return ok
		})
	}
}

func (s *Server) Disconnect(connID string) (err error) {
	s.LeaveAll(connID)
	if conn, ok := s.getConnection(connID); ok {
		conn.disconnected = true
		conn.fireDisconnect()
		err = conn.underline.Close()
		s.connections.Delete(connID)
	}
	return
}
