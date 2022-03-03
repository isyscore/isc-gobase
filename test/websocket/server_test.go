package websocket

import (
	"fmt"
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
	server.InitServer()
	server.RegisterWebSocketRoute("/ws", ws)
	server.StartServer()
}

func handleConnection(c websocket.Connection) {
	fmt.Println("client connected,id=", c.ID())
	c.Write(1, []byte("welcome client"))
	c.On("chat", func(msg string) {
		fmt.Printf("%s sent: %s\n", c.Context().ClientIP(), msg)
		c.To(websocket.All).Emit("chat", msg)
	})
	c.OnMessage(func(msg []byte) {
		fmt.Println("received msg:", string(msg))
		c.Write(1, []byte("hello aa"))
		c.To(websocket.All).Emit("chat", msg)
	})
	c.OnDisconnect(func() {
		fmt.Println("client Disconnect,id=", c.ID())
	})
}
