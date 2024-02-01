package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type wsNode struct {
	conn       *websocket.Conn
	writeMutex *sync.Mutex

	connected bool
}

func (w *wsNode) Read(p []byte) (n int, err error) {
	var readBytes []byte
	_, readBytes, err = w.conn.ReadMessage()
	copy(p, readBytes)
	n = len(readBytes)
	return
}

func (conn *wsNode) Write(p []byte) (n int, err error) {
	if !conn.connected {

	}
	err = conn.conn.WriteMessage(websocket.TextMessage, p)
	n = len(p)
	return
}

func (conn *wsNode) Valid() bool {
	return conn.connected && conn.conn != nil
}
