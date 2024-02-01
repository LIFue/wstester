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
		wsServer := ws.NewWsServer()
		wsServer.InitServerByGinContext(ctx)

		wsServer.Serve()
	})

	return r
}
