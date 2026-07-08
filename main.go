package main

import (
	"fmt"

	"github.com/Brady-Higgins/turtle/cmd"
	"github.com/Brady-Higgins/turtle/internal/cloudflare_client"
)

func main() {
	cmd.Execute()
	//c := cloudflare_client.New()
	//ctx := context.Background()
	//err := c.NewDNSRecord(ctx)
	//fmt.Println(err)
	//time.Sleep(time.Second * 8)
	//d, err := c.GetDNSRecord(ctx)
	////fmt.Println("other")
	////fmt.Println(d)
	//err = c.DeleteDNSRecord(d, ctx)
	//fmt.Println(err)
	err := cloudflare_client.CreateTunnelDNSRecord("example", "plotsearcher.com")
	if err != nil {
		fmt.Println(err)
	}

}
