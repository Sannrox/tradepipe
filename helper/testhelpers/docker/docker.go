package docker

import (
	"bufio"
	"context"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func GetDockerClient() (*client.Client, error) {
	return client.NewEnvClient()
}

func ExecDockerContainer(ctx context.Context, cli *client.Client, containerName string, cmd []string) ([]string, error) {
	exec, err := cli.ContainerExecCreate(ctx, containerName, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	})
	if err != nil {
		return nil, err
	}

	execStartCheck := types.ExecStartCheck{
		Tty: true,
	}
	hijackedResp, err := cli.ContainerExecAttach(ctx, exec.ID, execStartCheck)
	if err != nil {
		return nil, err
	}
	var output []string
	defer hijackedResp.Close()
	scanner := bufio.NewScanner(hijackedResp.Reader)
	for scanner.Scan() {
		output = append(output, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	cli.Close()
	return output, nil
}

func GetDockerContainers(ctx context.Context, cli *client.Client, listOptions types.ContainerListOptions) ([]types.Container, error) {
	return cli.ContainerList(ctx, listOptions)
}

func GetAllDockerContainers(ctx context.Context, cli *client.Client) ([]types.Container, error) {
	return GetDockerContainers(ctx, cli, types.ContainerListOptions{All: true})
}

func GetRunningDockerContainers(ctx context.Context, cli *client.Client) ([]types.Container, error) {
	return GetDockerContainers(ctx, cli, types.ContainerListOptions{All: false})
}

func CheckIfContainerExists(ctx context.Context, cli *client.Client, containerName string) (bool, error) {
	containers, err := GetAllDockerContainers(ctx, cli)
	if err != nil {
		return false, err
	}
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+containerName {
				return true, nil
			}
		}
	}
	return false, nil
}

func CheckIfContainerAllocatesPort(ctx context.Context, cli *client.Client, port int) (bool, error) {
	containers, err := GetAllDockerContainers(ctx, cli)
	if err != nil {
		return true, err
	}
	// Iterate over each container
	for _, container := range containers {
		// Get the container's port bindings
		containerPorts, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			return false, err
		}

		// Iterate over each port binding
		for _, portBinding := range containerPorts.NetworkSettings.Ports {
			// Check if the port is allocated
			if portBinding != nil && portBinding[0].HostPort == strconv.Itoa(port) {
				logrus.Warnf("Port %d is already allocated", port)
				return true, nil
			}
		}
	}

	// Port is not allocated
	return false, nil

}

func RunDockerContainer(ctx context.Context, cli *client.Client, containerName string, opts types.ContainerStartOptions) error {
	return cli.ContainerStart(ctx, containerName, opts)
}

func CreateDockerContainer(ctx context.Context, cli *client.Client, containerName string, config *container.Config, hostConfig *container.HostConfig) error {
	_, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}
	return nil
}

func StopDockerContainer(ctx context.Context, cli *client.Client, containerName string) error {
	return cli.ContainerStop(ctx, containerName, container.StopOptions{})
}

func RemoveDockerContainer(ctx context.Context, cli *client.Client, containerName string) error {
	return cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{})
}

func GetDockerImages(ctx context.Context, cli *client.Client) ([]types.ImageSummary, error) {
	return cli.ImageList(ctx, types.ImageListOptions{})
}

func CheckIfImageExists(ctx context.Context, cli *client.Client, imageName string) (bool, error) {
	images, err := GetDockerImages(ctx, cli)
	if err != nil {
		return false, err
	}
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == imageName {
				return true, nil
			}
		}
	}
	return false, nil
}

func PullDockerImage(ctx context.Context, cli *client.Client, imageName string) error {
	_, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	return nil
}
