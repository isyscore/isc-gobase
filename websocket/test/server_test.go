package test

import (
	"log"
	"testing"

	"github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/websocket"
)

func TestWebSocketServer(t *testing.T) {
	ws := websocket.NewWSServer(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	ws.OnConnection(handleConnection)
	server.RegisterWebSocketRoute("/ws", ws)
	server.StartServer()
}

func handleConnection(c websocket.Connection) {
	log.Println("client connected,id=", c.ID())
	_ = c.Write(1, []byte("welcome client"))
	c.On("chat", func(msg string) {
		log.Printf("%s sent: %s\n", c.Context().ClientIP(), msg)
		_ = c.To(websocket.All).Emit("chat", msg)
	})
	c.OnMessage(func(msg []byte) {
		log.Println("received msg:", string(msg))
		_ = c.Write(1, []byte("hello aa"))
		_ = c.To(websocket.All).Emit("chat", msg)
	})
	c.OnDisconnect(func() {
		log.Println("client Disconnect,id=", c.ID())
	})
}
