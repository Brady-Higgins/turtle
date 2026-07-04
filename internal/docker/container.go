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

func New() (*dockerClient, error) {
	cli, err := client.New(client.FromEnv)
	//cli, err := client.NewClientWithOpts(client.FromEnv)
	c := &dockerClient{Cli: cli}
	return c, err
}

// StartContainer : build a docker image provided an image name
// imageName : image-name:tag
func (d *dockerClient) StartContainer(imageName string, ctx context.Context) error {
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
	_, err = d.Cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return err
}
