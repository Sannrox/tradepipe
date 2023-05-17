package testing

import (
	"context"

	"github.com/Sannrox/tradepipe/helper/testhelpers/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/sirupsen/logrus"

	"github.com/docker/go-connections/nat"
)

func SetUpScylla(ctx context.Context) error {
	dockerClient, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	imageExists, err := docker.CheckIfImageExists(ctx, dockerClient, "scylladb/scylla:latest")
	if err != nil {
		return err
	}

	if !imageExists {
		err = docker.PullDockerImage(ctx, dockerClient, "scylladb/scylla:latest")
		if err != nil {
			return err
		}
	}

	containerExists, err := docker.CheckIfContainerExists(ctx, dockerClient, "scylla")
	if err != nil {
		return err
	}

	if containerExists {
		err = docker.RemoveDockerContainer(ctx, dockerClient, "scylla")
		if err != nil {
			return err
		}
	}

	containerConfig := &container.Config{
		Image: "scylladb/scylla:latest",
		Tty:   true,
		ExposedPorts: nat.PortSet{
			"9042/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"9042/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "9042",
				},
			},
		},
	}

	if err := docker.CreateDockerContainer(ctx, dockerClient, "scylla", containerConfig, hostConfig); err != nil {
		return err
	}

	if err := docker.RunDockerContainer(ctx, dockerClient, "scylla", types.ContainerStartOptions{}); err != nil {
		return err
	}

	logrus.Debug("Waiting for Scylla to start")

	return nil

}

func TearDownScylla(ctx context.Context) error {
	dockerClient, err := docker.GetDockerClient()
	if err != nil {
		return err
	}
	if err := docker.StopDockerContainer(ctx, dockerClient, "scylla"); err != nil {
		return err
	}

	if err := docker.RemoveDockerContainer(ctx, dockerClient, "scylla"); err != nil {
		return err
	}

	logrus.Debug("Scylla stopped")

	return nil
}
