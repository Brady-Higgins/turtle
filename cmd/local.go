package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Brady-Higgins/turtle/internal/cloudflare_client"
	"github.com/Brady-Higgins/turtle/internal/docker"
	"github.com/spf13/cobra"
)

type SelfHosting struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func activateSelfHosting(imgName string, ctx context.Context) {
	d, err := docker.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	//err = d.StartContainer("example-site:latest", ctx)
	id := d.GetContainerID(imgName, ctx)
	// container already exists for image
	if id != "" {
		err = d.StartContainer(id, ctx)
	} else { // need to build container first
		id, err = d.BuildContainer(imgName, ctx)
		d.StartContainer(id, ctx)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd := cloudflare_client.CreateCloudflaredCommand(ctx, os.Getenv("CLOUDFLARE_TUNNEL_NAME"))
	go cloudflare_client.RunCloudflared(cmd)
}

func checkSelfHosting() bool {
	return false
}

func deactivateSelfHosting(imgName string, cmd *exec.Cmd, cancel context.CancelFunc, ctx context.Context) {
	d, err := docker.New()
	id := d.GetContainerID(imgName, ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = d.StopContainer(id, ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cloudflare_client.StopCloudflared(cmd, cancel)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "switch to local hosting",
	Long:  "switch to local hosting",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		//imgName := "example-site"
		imgName, err := cmd.Flags().GetString("image")
		if imgName == "" {
			fmt.Println("Image flag required")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(imgName)
		fmt.Println("Starting Self Hosting...")
		activateSelfHosting(imgName, ctx)
		fmt.Printf("Self Hosting Started! Using %s\n", imgName)
		time.Sleep(time.Second * 10)
		deactivateSelfHosting(imgName, nil, cancel, ctx)
		fmt.Println("Self Hosting Turned Off!")
	},
}

func init() {
	onCmd.Flags().String("image", "", "docker image name to run")
}
