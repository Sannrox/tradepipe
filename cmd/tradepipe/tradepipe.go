package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sannrox/tradepipe/logger"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/sirupsen/logrus"
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
	Debug       bool
	LogFile     string
	OutputPath  string
	HistoryFile string
}

func NewTradePipeCmd() *cobra.Command {
	opts := &TradePipeOptions{}
	cmd := &cobra.Command{
		Use:              "tradepipe [number] [pin]",
		Short:            "tradepipe is a command line tool for interacting with the TradeMe API",
		Long:             `tradepipe is a command line tool for interacting with the TradeMe API.`,
		TraverseChildren: true,
		Args:             cobra.ExactArgs(2),
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

			if len(opts.OutputPath) == 0 {
				if path, err := os.Getwd(); err != nil {
					return err
				} else {
					logrus.Debug("No path set using current directory", path)
					opts.OutputPath = path
				}
			}
			return ExecuteCLI(args, opts.OutputPath, opts.HistoryFile)
		},
	}
	cmd.Flags().BoolVarP(&opts.Debug, "debug", "d", false, "Enable debug logging")
	cmd.Flags().StringVarP(&opts.LogFile, "logfile", "l", "", "Log file to write to")
	cmd.Flags().StringVarP(&opts.OutputPath, "out", "o", "", "Directory in which you want to write the files")
	cmd.Flags().StringVar(&opts.HistoryFile, "history", "history.txt", "From which you want to write and read the history")
	return cmd
}

func ExecuteCLI(args []string, outpath string, historyfile string) error {

	number := args[0]
	pin := args[1]

	client := tr.NewAPIClient()

	client.SetCredentials(number, pin)

	err := client.Login()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter token: ")
	text, _ := reader.ReadString('\n')

	string_b := strings.Replace(text, "\n", "", -1)
	intVar, err := strconv.Atoi(string_b)
	if err != nil {
		return err
	}
	err = client.VerifyLogin(intVar)
	if err != nil {
		return err
	}
	data := make(chan tr.Message)
	ctx := context.Background()

	err = client.NewWebSocketConnection(data)
	if err != nil {
		return err
	}

	time.Sleep(20 * time.Second)

	dl := tr.NewDownloader(client)

	dl.SetHistoryFile(historyfile)
	dl.SetOutputPath(outpath)
	err = dl.DownloadAll(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	cmd := NewTradePipeCmd()
	cmd.Execute()
}
