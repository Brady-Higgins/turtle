package docker

import (
	"context"
	"net/netip"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

type dockerClient struct {
	Cli *client.Client
}

// New : Creates a new docker client
func New() (*dockerClient, error) {
	cli, err := client.New(client.FromEnv)
	c := &dockerClient{Cli: cli}
	return c, err
}

// GetContainerID : return container ID if a container running the image already exists, else empty string
func (d *dockerClient) GetContainerID(imageName string, ctx context.Context) string {
	containers, _ := d.Cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	for _, c := range containers.Items {
		if imageName == c.Image {
			return c.ID
		}
	}
	return ""
}

// StartContainer : build a docker image provided an image name
// imageName : image-name:tag
func (d *dockerClient) BuildContainer(imageName string, ctx context.Context) (string, error) {
	containerPort, _ := network.ParsePort("80/tcp")
	localAdd := netip.MustParseAddr("0.0.0.0")
	localPortBindings := []network.PortBinding{
		{
			HostIP:   localAdd,
			HostPort: "8080",
		},
	}
	resp, err := d.Cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: imageName,
		},
		HostConfig: &container.HostConfig{
			PortBindings: network.PortMap{
				containerPort: localPortBindings,
			},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// StartContainer : Starts container given container ID
func (d *dockerClient) StartContainer(containerID string, ctx context.Context) error {
	_, err := d.Cli.ContainerStart(ctx, containerID, client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return err
}

// StopContainer : stops container given container ID
func (d *dockerClient) StopContainer(containerID string, ctx context.Context) error {
	_, err := d.Cli.ContainerStop(ctx, containerID, client.ContainerStopOptions{})
	if err != nil {
		return err
	}
	return err
}
