package cmd

import (
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "wstest",
	Short: "Wstester is useless ws tester",
	Run:   func(cmd *cobra.Command, args []string) {},
}

// func init() {
// 	rootCommand.PersistentFlags().StringP("cfgfilename", "c", "config", "config file name")
// 	rootCommand.PersistentFlags().StringArrayP("path", "p", []string{"."}, "config file path")
// }

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}
