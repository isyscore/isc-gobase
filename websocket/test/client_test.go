package test

import (
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait        = 10 * time.Second
	maxMessageSize   = 8192
	pongWait         = 15 * time.Second
	pingPeriod       = (pongWait * 9) / 10
	closeGracePeriod = 10 * time.Second

	WSURL = "127.0.0.1:8082"
)

func TestWebSocketClient(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: WSURL, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	defer func(ws *websocket.Conn) { _ = ws.Close() }(ws)

	done := make(chan bool)

	// ping
	go ping(ws, done)
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	message := []byte("hello, world")
	err = ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Send: %s\n", message)

	time.Sleep(time.Second * 30)
	done <- true

	log.Println("over")
}

func ping(ws *websocket.Conn, done chan bool) {
	ticker := time.NewTicker(pingPeriod)
	go func(t *time.Ticker) {
		defer t.Stop()
		for {
			select {
			case <-ticker.C:
				if err := ws.WriteControl(
					websocket.PingMessage,
					[]byte("ping"),
					time.Now().Add(writeWait)); err != nil {
					log.Println("ping:", err)
				} else {
					log.Println("ping ok")
				}

			case <-done:
				log.Println("done")
				return
			}
		}
	}(ticker)
}
