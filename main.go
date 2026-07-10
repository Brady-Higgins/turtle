package main

import (
	"fmt"
	"time"

	"github.com/Brady-Higgins/turtle/cmd"
	"github.com/Brady-Higgins/turtle/internal/windows_client"
)

func main() {
	cmd.Execute()
	//c := cloudflare_client.New()
	//ctx, cancel := context.WithCancel(context.Background())
	//cmd := cloudflare_client.CreateCloudflaredCommand(ctx, "example-site")
	//go cloudflare_client.RunCloudflared(cmd)
	//
	//fmt.Println("Started cloudflared")
	//time.Sleep(time.Second * 10)
	//cloudflare_client.StopCloudflared(cmd, cancel)
	//
	//fmt.Println("Stopped cloudflared")
	//time.Sleep(time.Second * 10)
	w, err := windows_client.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = w.StartService()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(10 * time.Second)
	err = w.StopService()
	if err != nil {
		fmt.Println(err)
		return
	}
}
