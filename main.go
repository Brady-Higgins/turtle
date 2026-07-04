package main

import (
	"context"
	"fmt"

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
	err = d.StartContainer("example-site:latest", ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}
