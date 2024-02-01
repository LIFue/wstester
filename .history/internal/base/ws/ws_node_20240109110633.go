package ws

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type wsNode struct {
	conn       *websocket.Conn
	writeMutex *sync.Mutex
}

func (w *wsNode) Read(p []byte) (n int, err error) {
	var readBytes []byte
	_, readBytes, err = w.conn.ReadMessage()
	copy(p, readBytes)
	n = len(readBytes)
	return
}

func (conn *wsNode) Write(p []byte) (n int, err error) {
	if conn.isClose {
		return 0, errors.New("ws is close")
	}
	out := make([]byte, len(p))
	copy(out[0:], p)
	err = conn.conn.WriteMessage(websocket.TextMessage, out)
	n = len(p)
	return
}
