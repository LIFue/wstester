package server

import (
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
	wsManager ws.WsManager,
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
		if err := wsManager.UpgradeHttpToWsAndServer(ctx.Writer, ctx.Request); err != nil {
			ctx.JSON(500, gin.H{
				"result": "failed",
				"error":  err.Error(),
			})
		}

		ctx.JSON(200, gin.H{"result": "success"})
	})

	return r
}
