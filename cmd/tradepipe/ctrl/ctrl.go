package ctrl

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Sannrox/tradepipe/grpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCtrlCommand() *cobra.Command {
	var host string
	var port string
	cmd := &cobra.Command{
		Use:   "ctrl [number] [pin]",
		Short: "Utility command for the tradepipe service",
		Long:  `Utility command for the tradepipe service`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteCLI(args)
		},
	}

	cmd.Flags().StringVar(&host, "host", "localhost", "Host to connect to")
	cmd.Flags().StringVarP(&port, "port", "p", "50051", "Port to connect to")

	return cmd
}

func ExecuteCLI(args []string) error {
	number := args[0]
	pin := args[1]

	client := grpc.NewClient()

	err := client.Connect(fmt.Sprintf("%s:%s", "localhost", "50051"))
	if err != nil {
		return err
	}

	defer client.Close()
	logrus.Info("Connected to server")

	processId, err := client.Login(number, pin)
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

	login, err := client.Verify(processId.ProcessId, int32(intVar))
	if err != nil {
		return err
	}
	if login.Error != "" {
		return fmt.Errorf("error: %s", login.Error)
	}

	// _, err = client.Positions(processId.ProcessId)
	// if err != nil {
	// 	return err
	// }

	// _, err = client.Timeline(processId.ProcessId)
	// if err != nil {
	// 	return err
	// }
	_, err = client.SavingsPlans(processId.ProcessId)
	if err != nil {
		return err
	}

	// var positionsJson []tr.Position

	// err = json.Unmarshal([]byte(positions.Postions), &positionsJson)
	// if err != nil {
	// 	return err
	// }

	return nil
}
