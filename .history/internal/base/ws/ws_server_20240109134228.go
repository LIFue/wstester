package ws

import (
	"net/rpc"
	"wstester/pkg/jsonrpc"
	"wstester/pkg/log"

	"github.com/gin-gonic/gin"
)

type WsServer struct {
	node *wsNode
}

func NewWsServer() *WsServer {
	return &WsServer{
		node: NewWsNode(),
	}
}

func (server *WsServer) InitServerByGinContext(ctx *gin.Context) error {
	if err := server.node.UpgradeHttp(ctx.Writer, ctx.Request); err != nil {
		log.Errorf("upgrade http to websocket error: %s", err.Error())
		return err
	}
	return nil
}

func (server *WsServer) Serve() {
	rpc.ServeCodec(jsonrpc.NewServerCodec(server.node))
}
