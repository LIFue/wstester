package server

import (
	"net/http"
	"net/rpc"
	"wstester/internal/base/ws"
	"wstester/internal/controller"
	"wstester/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewWebSocketServer(
	platformController *controller.Platform,
	arith *controller.Arith,
	messageController *controller.MessageController,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	if err := rpc.Register(platformController); err != nil {
		return nil
	}

	if err := rpc.Register(arith); err != nil {
		return nil
	}

	rpc.Register(messageController)

	r.GET("/ws", func(ctx *gin.Context) {
		wsConn, err := ws.InitWs(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		wsConn.Start()
	})

	return r
}
