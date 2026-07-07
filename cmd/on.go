package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Brady-Higgins/turtle/internal/cloudflare_client"
	"github.com/Brady-Higgins/turtle/internal/docker"
	"github.com/spf13/cobra"
)

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
	go cloudflare_client.RunCloudflared("example")
}

func checkSelfHosting() bool {
	return false
}

func deactivateSelfHosting(imgName string, ctx context.Context) {
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
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "turn on self hosting",
	Long:  "turn on self hosting",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
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
		deactivateSelfHosting(imgName, ctx)
		fmt.Println("Self Hosting Turned Off!")
	},
}

func init() {
	onCmd.Flags().String("image", "", "docker image name to run")
}
