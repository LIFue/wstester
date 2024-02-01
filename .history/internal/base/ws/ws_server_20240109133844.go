package ws

import "github.com/gin-gonic/gin"

type WsServer struct {
	node *wsNode
}

func NewWsServer() *WsServer {
	return &WsServer{}
}

func (server *WsServer) InitByGinContext(ctx *gin.Context) error {

}
