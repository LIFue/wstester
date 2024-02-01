package ws

import (
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

func (server *WsServer) InitByGinContext(ctx *gin.Context) error {
	if err := server.node.UpgradeHttp(ctx.Writer, ctx.Request); err != nil {
		log.Errorf("upgrade http to websocket error: %s", err.Error())
		return err
	}
}
