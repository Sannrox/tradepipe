package container

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/Sannrox/tradepipe/helper/testhelpers/docker"
	"github.com/Sannrox/tradepipe/helper/testhelpers/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/docker/go-connections/nat"
)

func SetUpScylla(ctx context.Context, startPort, endPort int) (string, int, error) {

	id := uuid.New().String()
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	containerName := "scylla-" + id + "-" + filepath.Base(callerName)
	dockerClient, err := docker.GetDockerClient()
	if err != nil {
		return "", 0, err
	}

	imageExists, err := docker.CheckIfImageExists(ctx, dockerClient, "scylladb/scylla:latest")
	if err != nil {
		return "", 0, err
	}

	if !imageExists {
		err = docker.PullDockerImage(ctx, dockerClient, "scylladb/scylla:latest")
		if err != nil {
			return "", 0, err
		}
	}

	containerExists, err := docker.CheckIfContainerExists(ctx, dockerClient, containerName)
	if err != nil {
		return "", 0, err
	}

	if containerExists {
		err = docker.RemoveDockerContainer(ctx, dockerClient, containerName)
		if err != nil {
			return "", 0, err
		}
	}

	freePort := 0
	for i := startPort; i <= endPort; i++ {
		freePort, err = utils.FindFreePort(startPort, endPort)
		if err != nil {
			return "", 0, err
		}

		isAllocated, err := docker.CheckIfContainerAllocatesPort(ctx, dockerClient, freePort)

		if err != nil {
			return "", 0, err
		}

		if !isAllocated {
			break
		}

		if i == endPort {
			return "", 0, fmt.Errorf("No free port found")
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
					HostPort: strconv.Itoa(freePort),
				},
			},
		},
	}

	if err := docker.CreateDockerContainer(ctx, dockerClient, containerName, containerConfig, hostConfig); err != nil {
		return "", 0, err
	}

	if err := docker.RunDockerContainer(ctx, dockerClient, containerName, types.ContainerStartOptions{}); err != nil {
		return "", 0, err
	}

	return containerName, freePort, nil

}

func TearDownScylla(containerName string, ctx context.Context) error {
	dockerClient, err := docker.GetDockerClient()
	if err != nil {
		return err
	}
	if err := docker.StopDockerContainer(ctx, dockerClient, containerName); err != nil {
		return err
	}

	if err := docker.RemoveDockerContainer(ctx, dockerClient, containerName); err != nil {
		return err
	}

	logrus.Debug("Scylla stopped")

	return nil
}
