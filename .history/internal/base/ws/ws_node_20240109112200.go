package ws

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type wsNode struct {
	conn       *websocket.Conn
	writeMutex *sync.Mutex

	connected bool
}

var upgrade = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,

	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWsNode(w http.ResponseWriter, r *http.Request) *wsNode {
	c, err := upgrade.Upgrade(w, r, nil)
}

func (w *wsNode) Read(p []byte) (n int, err error) {
	var readBytes []byte
	_, readBytes, err = w.conn.ReadMessage()
	copy(p, readBytes)
	n = len(readBytes)
	return
}

func (conn *wsNode) Write(p []byte) (n int, err error) {
	if !conn.Valid() {
		return 0, errors.New("web socket connection is closed")
	}
	err = conn.conn.WriteMessage(websocket.TextMessage, p)
	n = len(p)
	return
}

func (conn *wsNode) Valid() bool {
	return conn.connected && conn.conn != nil
}
