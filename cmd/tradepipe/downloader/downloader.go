package downloader

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sannrox/tradepipe/tr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewDownloaderCommand() *cobra.Command {
	var outpath string
	var historyfile string
	var cmd = &cobra.Command{
		Use:   "downloader [number] [pin]",
		Short: "Download all files from the Trade Republic API",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(outpath) == 0 {
				if path, err := os.Getwd(); err != nil {
					return err
				} else {
					logrus.Debug("No path set using current directory", path)
					outpath = path
				}
			}
			return ExecuteCLI(args, outpath, historyfile)
		},
	}
	cmd.Flags().StringVarP(&outpath, "outpath", "o", "", "Path to store the downloaded files")
	cmd.Flags().StringVarP(&historyfile, "historyfile", "f", "", "Path to store the history file")
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
