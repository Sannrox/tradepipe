package main

import (
	"fmt"
	"time"

	"github.com/Sannrox/tradepipe/gear/server"
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
	DBPort     int
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

			s := server.NewServer()

			return s.Run(opts.Done, opts.DB, opts.DBPort, opts.DBAttempts, opts.DBTimeouts)
		},
	}

	cmd.Flags().IntVarP(&opts.DBAttempts, "dbattempts", "a", 10, "Number of attempts to connect to the database")
	cmd.Flags().DurationVarP(&opts.DBTimeouts, "dbtimeouts", "t", 10*time.Second, "Timeout for connecting to the database")
	cmd.Flags().BoolVarP(&opts.Debug, "debug", "d", false, "Enable debug logging")
	cmd.Flags().StringVarP(&opts.LogFile, "logfile", "l", "tradegrpc.log", "Log file to write to")
	cmd.Flags().StringVarP(&opts.DB, "db", "b", "localhost", "Database host")
	cmd.Flags().IntVarP(&opts.DBPort, "dbport", "p", 9042, "Database port")

	return cmd
}

func main() {
	cmd := NewTadeGrpcCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
