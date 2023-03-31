package main

import (
	"fmt"

	"github.com/Sannrox/tradepipe/cmd/tradepipe/ctrl"
	"github.com/Sannrox/tradepipe/cmd/tradepipe/downloader"
	"github.com/Sannrox/tradepipe/logger"
	"github.com/spf13/cobra"
)

var (
	PlatformName = ""
	Version      = "unknown-version"
	GitCommit    = "unknown-commit"
	BuildTime    = "unknown-buildtime"
	BuildArch    = "unknown-buildarch"
	BuildOs      = "unknown-buildos"
)

type TradePipeOptions struct {
	Debug   bool
	LogFile string
}

func NewTradePipeCmd() *cobra.Command {
	opts := &TradePipeOptions{}
	cmd := &cobra.Command{
		Use:              "tradepipe",
		Short:            "tradepipe is a command line tool for interacting with the Trade Republic API",
		Long:             `tradepipe is a command line tool for interacting with the Trade Republic API.`,
		TraverseChildren: true,
		Version:          fmt.Sprintf("%s, built: %s ", Version, GitCommit),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if opts.Debug {
				logger.Enable()
			}
			if len(opts.LogFile) != 0 {
				if err := logger.SetLogFile(opts.LogFile); err != nil {
					panic(err)
				}
			}
		},
	}
	cmd.PersistentFlags().BoolVarP(&opts.Debug, "debug", "d", false, "Enable debug logging")
	cmd.PersistentFlags().StringVarP(&opts.LogFile, "logfile", "l", "", "Log file to write to")
	AddCommands(cmd)
	return cmd
}

func main() {
	cmd := NewTradePipeCmd()
	cmd.Execute()
}

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(downloader.NewDownloaderCommand())
	cmd.AddCommand(ctrl.NewCtrlCommand())
}
