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
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type WsConn struct {
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan []byte
	isClose   bool // 通道closeChan是否已经关闭
	mutex     sync.Mutex
	conn      *websocket.Conn
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

	return &WsConn{
		conn: conn,
	}, nil
}

func NewWsConn(conn *websocket.Conn) *WsConn {
	return &WsConn{
		conn: conn,
	}
}

func (conn *WsConn) Read(p []byte) (n int, err error) {
	messageType, data, err := conn.conn.ReadMessage()
	if err != nil {
		return 0, err
	}

	if messageType != websocket.TextMessage {
		return 0, errors.New("message type is wrong")
	}

	copy(p, data)
	n = len(data)
	return
}

func (conn *WsConn) Write(p []byte) (n int, err error) {
	out := make([]byte, len(p))
	copy(out[0:], p)
	err = conn.conn.WriteMessage(websocket.TextMessage, out)
	n = len(p)
	return
}

func (conn *WsConn) Close() error {
	log.Info("wsconn close")
	wc, err := conn.conn.NextWriter(websocket.CloseMessage)
	if err != nil {
		return err
	}
	n, err := wc.Write([]byte("ws close"))
	if err != nil {
		return err
	}
	if err := conn.conn.Close(); err != nil {
		log.Errorf("close websocket error: %s", err.Error())
		return err
	}
	return nil
}

func (conn *WsConn) Start() {
	conn.serve()
}

func (conn *WsConn) serve() {
	rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
}
