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
	Outpath     string
	Historyfile string
}

func NewTradePipeCmd() *cobra.Command {
	opts := &TradePipeOptions{}
	cmd := &cobra.Command{
		Use:     "tradepipe [number] [pin]",
		Short:   "Download all files from the Trade Republic API",
		Long:    `tradepipe is a command line tool for interacting with the Trade Republic API.`,
		Version: fmt.Sprintf("%s, built: %s ", Version, GitCommit),
		Args:    cobra.ExactArgs(2),
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(opts.Outpath) == 0 {
				if path, err := os.Getwd(); err != nil {
					return err
				} else {
					logrus.Debug("No path set using current directory", path)
					opts.Outpath = path
				}
			}
			if len(opts.Historyfile) == 0 {
				return fmt.Errorf("historyfile cannot be empty")
			}

			return ExecuteCLI(args, opts.Outpath, opts.Historyfile)
		},
	}
	cmd.Flags().BoolVarP(&opts.Debug, "debug", "d", false, "Enable debug logging")
	cmd.Flags().StringVarP(&opts.LogFile, "logfile", "l", "", "If set, write logs to this file instead of stdout")
	cmd.Flags().StringVarP(&opts.Outpath, "outpath", "o", "", "Path to store the downloaded files")
	cmd.Flags().StringVarP(&opts.Historyfile, "historyfile", "f", "history.txt", "Path to store the history file")

	return cmd
}

func main() {
	cmd := NewTradePipeCmd()
	cmd.Execute()
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
