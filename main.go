package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Brady-Higgins/turtle/cmd"
	"github.com/Brady-Higgins/turtle/internal/docker"
)

func main() {
	cmd.Execute()

	d, err := docker.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	imgName := "example-site"
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

	time.Sleep(time.Second * 10)
	err = d.StopContainer(id, ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}
