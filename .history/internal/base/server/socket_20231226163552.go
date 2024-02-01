package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/rpc"
	"wstester/internal/base/ws"
	"wstester/internal/controller"
	"wstester/internal/middleware"
)

func NewWebSocketServer(
	platformController *controller.Platform,
	arith *controller.Arith,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	if err := rpc.Register(platformController); err != nil {
		return nil
	}

	if err := rpc.Register(arith); err != nil {
		return nil
	}

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
