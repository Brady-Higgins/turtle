package main

import (
	"context"

	"github.com/Brady-Higgins/turtle/cmd"
	"github.com/Brady-Higgins/turtle/internal/cloudflare_client"
)

func main() {
	cmd.Execute()
	c := cloudflare_client.New()
	ctx := context.Background()
	c.NewDNSRecord(ctx)
}
