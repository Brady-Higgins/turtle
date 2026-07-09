package main

import (
	"context"
	"fmt"

	"github.com/Brady-Higgins/turtle/cmd"
	"github.com/Brady-Higgins/turtle/internal/cloudflare_client"
	"github.com/cloudflare/cloudflare-go/v7/dns"
)

func main() {
	cmd.Execute()
	c := cloudflare_client.New()
	ctx := context.Background()
	//err := c.NewDNSRecord(ctx)
	//fmt.Println(err)
	//time.Sleep(time.Second * 8)
	d, err := c.GetDNSRecord(dns.RecordListParamsTypeCNAME, false, ctx)
	fmt.Println(err)
	if d == nil {
		fmt.Println("No record")
		return
	}

	fmt.Println("record exists")
	//err = c.CommentDNSRecord(d, ctx)
	//fmt.Println(err)
	////fmt.Println("other")
	////fmt.Println(d)
	//err = c.DeleteDNSRecord(d, ctx)
	//fmt.Println(err)
	//err := cloudflare_client.CreateTunnelDNSRecord()
	//if err != nil {
	//	fmt.Println(err)
	//}

}
