package main

import (
	"fmt"
	"time"

	"github.com/Sannrox/tradepipe/gear"
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
	Debug      bool
	LogFile    string
	Done       chan struct{}
	DBAttempts int
	DBTimeouts time.Duration
	DB         string
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

			server := gear.NewGRPCServer()

			if err := server.CreateKeySpaceConnection(opts.DB, opts.DBAttempts, opts.DBTimeouts); err != nil {
				return err
			}

			return server.Run(opts.Done)
		},
	}

	cmd.Flags().IntVarP(&opts.DBAttempts, "dbattempts", "a", 10, "Number of attempts to connect to the database")
	cmd.Flags().DurationVarP(&opts.DBTimeouts, "dbtimeouts", "t", 10*time.Second, "Timeout for connecting to the database")
	cmd.Flags().BoolVarP(&opts.Debug, "debug", "d", false, "Enable debug logging")
	cmd.Flags().StringVarP(&opts.LogFile, "logfile", "l", "tradegrpc.log", "Log file to write to")
	cmd.Flags().StringVarP(&opts.DB, "db", "b", "localhost", "Database host")

	return cmd
}

func main() {
	cmd := NewTadeGrpcCmd()
	cmd.Execute()
}
