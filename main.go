package main

import (
	"context"
	"fmt"

	"github.com/Brady-Higgins/turtle/cmd"
	"github.com/Brady-Higgins/turtle/internal/cloudflare_client"
)

func main() {
	cmd.Execute()
	c := cloudflare_client.New()
	ctx := context.Background()
	//err := c.NewDNSRecord(ctx)
	//fmt.Println(err)
	d, err := c.GetDNSRecord(ctx)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(d)
	}

}
