package docker

import (
	"bufio"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func GetDockerClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return cli, nil
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

func GetDockerContainers(ctx context.Context, cli *client.Client) ([]types.Container, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func CheckIfContainerExists(ctx context.Context, cli *client.Client, containerName string) (bool, error) {
	containers, err := GetDockerContainers(ctx, cli)
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

func RunDockerContainer(ctx context.Context, cli *client.Client, containerName string, opts types.ContainerStartOptions) error {
	err := cli.ContainerStart(ctx, containerName, opts)
	if err != nil {
		return err
	}
	return nil
}

func CreateDockerContainer(ctx context.Context, cli *client.Client, containerName string, config *container.Config, hostConfig *container.HostConfig) error {
	_, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}
	return nil
}

func StopDockerContainer(ctx context.Context, cli *client.Client, containerName string) error {
	err := cli.ContainerStop(ctx, containerName, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}

func RemoveDockerContainer(ctx context.Context, cli *client.Client, containerName string) error {
	err := cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GetDockerImages(ctx context.Context, cli *client.Client) ([]types.ImageSummary, error) {
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
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
