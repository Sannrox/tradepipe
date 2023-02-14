package main

import (
	"fmt"
	"os"

	"github.com/Sannrox/tradepipe/cmd/cli"
	"github.com/Sannrox/tradepipe/cmd/grpc"
	"github.com/Sannrox/tradepipe/cmd/rest"
	"github.com/Sannrox/tradepipe/logger"
	_ "github.com/Sannrox/tradepipe/logger"
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

type RootOptions struct {
	Verbose bool
	GRPC    bool
	HTTP    bool
	Debug   bool
	LogFile string
	Done    chan struct{}
}

func NewRootCmd() *cobra.Command {
	rootOptions := &RootOptions{}
	cmd := &cobra.Command{
		Use:              "tradepipe",
		Short:            "tradepipe is a command line tool for interacting with the TradeMe API",
		Long:             `tradepipe is a command line tool for interacting with the TradeMe API.`,
		TraverseChildren: true,
		Version:          fmt.Sprintf("%s, built: %s ", Version, GitCommit),
		RunE: func(cmd *cobra.Command, args []string) error {
			if rootOptions.Debug {
				logger.Enable()
			}
			if err := logger.SetLogFile(rootOptions.LogFile); err != nil {
				return err
			}
			switch {
			case rootOptions.GRPC:
				// Run GRPC server

				server := grpc.NewGRPCServer()
				return server.Run()
			case rootOptions.HTTP:
				// Run HTTP server
				server := rest.NewRestServer()
				return server.Run(rootOptions.Done, "8080")
			default:
				// Run CLI
				if len(os.Args) < 2 {
					_ = cmd.Help()
					return nil
				}
				return cli.ExecuteCLI(os.Args)
			}
		},
	}

	cmd.Flags().StringVarP(&rootOptions.LogFile, "out", "o", "tradepip.txt", "Outputfile")
	cmd.Flags().BoolVarP(&rootOptions.Verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&rootOptions.GRPC, "grpc", "g", false, "run grpc server")
	cmd.Flags().BoolVarP(&rootOptions.HTTP, "http", "r", false, "run http server")
	cmd.Flags().BoolVarP(&rootOptions.Debug, "debug", "D", false, "debug mode")

	return cmd
}

func main() {
	cmd := NewRootCmd()
	cmd.Execute()
}
