package cmd

import (
	"context"
	"os"
	"os/signal"

	"wstester/internal/base/conf"
	"wstester/pkg/log"

	"github.com/spf13/cobra"
)

var wsCommand = &cobra.Command{
	Use:   "ws",
	Short: "start the websocket server",
	Run:   runWsServer,
}

func init() {
	wsCommand.Flags().StringP("config", "c", "./config/config.yaml", "set the config file path")
	rootCommand.AddCommand(wsCommand)
}

func runWsServer(cmd *cobra.Command, args []string) {
	configFilePath, err := cmd.Flags().GetString("config")
	if err != nil {
		return
	}

	allConfig, err := conf.ReadConfig(configFilePath)
	if err != nil {
		return
	}

	log.Infof("allconfig: %+v", allConfig)
	ctx, cancleFunc := context.WithCancel(context.Background())
	r, err := InitializeWsServer(allConfig.Debug, ctx, allConfig.Data.Database)
	if err != nil {
		panic(err)
	}
	go func() {
		if err = r.Run(allConfig.Server.WebSocket.Addr); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	for range quit {
		cancleFunc()
		log.Info("end...")
		return
	}
}
