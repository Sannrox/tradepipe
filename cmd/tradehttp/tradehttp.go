package main

import (
	"fmt"

	"github.com/Sannrox/tradepipe/logger"
	"github.com/Sannrox/tradepipe/rest"
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

type TradeHttpOptions struct {
	Debug   bool
	LogFile string
	Port    string
	Done    chan struct{}
}

func NewTradeHttpCmd() *cobra.Command {
	opts := &TradeHttpOptions{}
	cmd := &cobra.Command{
		Use:              "tradehttp",
		Short:            "tradehttp is a microservice with REST for interacting with the TradeMe API",
		Long:             `tradehttp is a microservice with REST for interacting with the TradeMe API.`,
		TraverseChildren: true,
		Version:          fmt.Sprintf("%s, built: %s ", Version, GitCommit),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Debug {
				logger.Enable()
			}
			if len(opts.LogFile) != 0 {
				if err := logger.SetLogFile(opts.LogFile); err != nil {
					return err
				}
			}
			server := rest.NewRestServer()
			return server.Run(opts.Done, opts.Port)
		},
	}
	cmd.Flags().BoolVarP(&opts.Debug, "debug", "d", false, "Enable debug logging")
	cmd.Flags().StringVarP(&opts.LogFile, "logfile", "l", "", "Log file to write to")
	cmd.Flags().StringVarP(&opts.Port, "port", "p", "8080", "Port to listen on")

	return cmd
}

func main() {
	cmd := NewTradeHttpCmd()
	cmd.Execute()
}
