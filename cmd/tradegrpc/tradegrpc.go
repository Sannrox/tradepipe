package main

import (
	"fmt"

	"github.com/Sannrox/tradepipe/grpc"
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

type TradeGrpcOptions struct {
	Debug   bool
	LogFile string
	Done    chan struct{}
}

func NewTadeGrpcCmd() *cobra.Command {
	opts := &TradeGrpcOptions{}
	cmd := &cobra.Command{
		Use:              "tradegrpc",
		Short:            "tradegrpc is a microservice with protobuffer for interacting with the TradeMe API",
		Long:             `tradegrpc is a microservice with protobuffer for interacting with the TradeMe API.`,
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
			server := grpc.NewGRPCServer()
			return server.Run(opts.Done)
		},
	}
	cmd.Flags().BoolVarP(&opts.Debug, "debug", "d", false, "Enable debug logging")
	cmd.Flags().StringVarP(&opts.LogFile, "logfile", "l", "", "Log file to write to")

	return cmd
}

func main() {
	cmd := NewTadeGrpcCmd()
	cmd.Execute()
}
