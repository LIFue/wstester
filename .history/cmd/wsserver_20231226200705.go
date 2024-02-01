package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"wstester/internal/base/conf"
	_ "wstester/internal/controller"
	"wstester/pkg/log"
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
	r, err := InitializeWsServer(allConfig.Debug, allConfig.Data.Database)
	if err != nil {
		panic(err)
	}
	go func() {
		if err = r.Run(allConfig.Server.WebSocket.Addr); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	for range quit {
		log.Info("end...")
		return
	}
}
