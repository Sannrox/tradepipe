package main

import (
	"fmt"
	"os"

	"github.com/Sannrox/tradepipe/cmd/cli"
	"github.com/Sannrox/tradepipe/cmd/grpc"
	"github.com/Sannrox/tradepipe/cmd/rest"
	"github.com/Sannrox/tradepipe/pkg/logger"
	_ "github.com/Sannrox/tradepipe/pkg/logger"
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
			switch {
			case rootOptions.GRPC:
				// Run GRPC server

				server := grpc.NewGRPCServer()
				if err := server.Run(); err != nil {
					return err
				}
			case rootOptions.HTTP:
				// Run HTTP server
				server := rest.NewRestServer()
				if err := server.Run(); err != nil {
					return err
				}
			default:
				// Run CLI
				if len(os.Args) < 2 {
					_ = cmd.Help()
					return nil
				}
				cli.ExecuteCLI(os.Args)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&rootOptions.Verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&rootOptions.GRPC, "grpc", "g", false, "run grpc server")
	cmd.Flags().BoolVarP(&rootOptions.HTTP, "http", "r", false, "run http server")
	cmd.Flags().BoolVarP(&rootOptions.Debug, "debug", "D", false, "debug mode")

	return cmd
}

func main() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
