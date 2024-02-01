package ws

import (
	"net/http"
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

var upgrade = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,

	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWs(ctx *gin.Context) (*WsConn, error) {
	if ctx == nil || ctx.Writer == nil || ctx.Request == nil {
		return nil, errors.New("connection is wrong")
	}
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}

	ws := &WsConn{
		id:       uuid.New().String(),
		isClose:  false,
		conn:     conn,
		mutex:    sync.Mutex{},
		liveTime: time.Now().Add(1 * time.Minute),
	}

	AddWs(ws)
	return ws, nil
}

func NewWsConn(conn *websocket.Conn) *WsConn {
	return &WsConn{
		conn: conn,
	}
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
	conn.conn
	return conn.conn.Close()
}

func (conn *WsConn) Start() {
	conn.serve()
}

func (conn *WsConn) serve() {
	rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
}
