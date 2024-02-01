package server

import (
	"net/http"
	"wstester/internal/base/ws"
	"wstester/internal/controller"
	"wstester/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewWebSocketServer(
	platformController *controller.Platform,
	arith *controller.Arith,
	messageController *controller.Message,
	controllerRegister *controller.ControllerRegister,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	controllerRegister.AddService(platformController)
	controllerRegister.AddService(arith)
	controllerRegister.AddService(messageController)
	if err := controllerRegister.Register(); err != nil {
		return nil
	}
	r.GET("/ws", func(ctx *gin.Context) {
		wsConn, err := ws.InitWs(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		go wsConn.Start()
	})

	return r
}
