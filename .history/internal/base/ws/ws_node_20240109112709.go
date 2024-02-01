package ws

import (
	"errors"
	"net/http"
	"time"
	"wstester/pkg/log"

	"github.com/gorilla/websocket"
)

type wsNode struct {
	conn *websocket.Conn

	connected bool
}

var upgrade = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,

	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWsNode() *wsNode {
	return &wsNode{
		connected: false,
	}
}

func (conn *wsNode) UpgradeHttp(w http.ResponseWriter, r *http.Request) (err error) {
	conn.conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("http upgrade to websocket error: %s", err.Error())
		return err
	}
	return nil
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
