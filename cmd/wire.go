//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"wstester/internal/base"
	"wstester/internal/base/data"
	"wstester/internal/base/server"
	"wstester/internal/controller"
	"wstester/internal/repo"
	"wstester/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitializeWsServer(debug bool, ctx context.Context, dbConf *data.Database) (*gin.Engine, error) {
	wire.Build(
		base.BaseSet,
		repo.RepoSet,
		controller.ControllerSet,
		service.ServiceSet,
		server.NewWebSocketServer,
	)
	return nil, nil
}
