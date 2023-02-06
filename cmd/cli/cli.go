package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sannrox/tradepipe/pkg/tr"
)

func ExecuteCLI(args []string) error {

	number := args[1]
	pin := args[2]

	client := tr.NewAPIClient()

	client.SetCredentials(number, pin)

	err := client.Login()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Token: ")
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
	dl := tr.NewDownloader(*client)
	dl.DownloadAll(ctx, data)

	return nil
}
