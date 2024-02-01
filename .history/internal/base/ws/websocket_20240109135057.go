package ws

import (
	"net/rpc"
	"sync"
	"time"

	// "net/rpc/jsonrpc"

	"wstester/pkg/jsonrpc"
	"wstester/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type WsConn struct {
	id       string
	isClose  bool // 通道closeChan是否已经关闭
	mutex    sync.Mutex
	conn     *websocket.Conn
	liveTime time.Time
}

func (conn *WsConn) Read(p []byte) (n int, err error) {
	_, data, err := conn.conn.ReadMessage()
	if err != nil {
		return 0, err
	}

	conn.liveTime = time.Now().Add(5 * time.Minute)

	copy(p, data)
	n = len(data)
	return
}

func (conn *WsConn) Write(p []byte) (n int, err error) {
	if conn.isClose {
		return 0, errors.New("ws is close")
	}
	out := make([]byte, len(p))
	copy(out[0:], p)
	err = conn.conn.WriteMessage(websocket.TextMessage, out)
	n = len(p)
	return
}

func (conn *WsConn) Close() error {
	log.Info("wsconn close")

	conn.isClose = true
	return nil
}

func (conn *WsConn) TryClose() error {
	return conn.conn.Close()
}

func (conn *WsConn) Start() {
	conn.serve()
}

func (conn *WsConn) serve() {
	rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
}
