package ws

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"wstester/pkg/log"

	"github.com/gorilla/websocket"
)

type wsNode struct {
	conn *websocket.Conn

	connected bool

	writeLocker sync.Mutex
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

func (node *wsNode) UpgradeHttp(w http.ResponseWriter, r *http.Request) (err error) {
	node.conn, err = upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("http upgrade to websocket error: %s", err.Error())
		return err
	}
	node.connected = true
	return nil
}

func (node *wsNode) ConnectWsServer(url string, requestHeaders http.Header) error {
	var wsClient = websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}

	conn, resp, err := wsClient.Dial(url, requestHeaders)
	if err != nil {
		log.Errorf("connect to web socket server error: %s", err.Error())
		return err
	}
	if resp.StatusCode != 200 {
		log.Errorf("connect to web socket server error, response status is wrong, status code: %d", resp.StatusCode)
		return errors.New("response status is wrong")
	}
	if conn != nil {
		node.conn = conn
		node.connected = true
	}

	return nil
}

func (node *wsNode) Read(p []byte) (n int, err error) {
	var readBytes []byte

	if !node.Valid() {
		return 0, errors.New("web socket connection is closed")
	}

	_, readBytes, err = node.conn.ReadMessage()
	copy(p, readBytes)
	n = len(readBytes)
	return
}

func (node *wsNode) Write(p []byte) (n int, err error) {
	if !node.Valid() {
		return 0, errors.New("web socket connection is closed")
	}
	err = node.conn.WriteMessage(websocket.TextMessage, p)
	n = len(p)
	return
}

func (node *wsNode) Close() error {
	log.Info("wsconn close")
	if node.conn != nil {
		if err := node.conn.Close(); err != nil {
			log.Errorf("close web socket connection error: %s", err.Error())
		}
	}
	node.conn = nil
	node.connected = false
	return nil
}

func (node *wsNode) Valid() bool {
	return node.connected && node.conn != nil
}
