package cloudflare_client

import (
	"github.com/cloudflare/cloudflare-go/v7"
)

type cloudflareClient struct {
	Cli *cloudflare.Client
}

// New : Creates a new docker client
func New() *cloudflareClient {
	client := cloudflare.NewClient()
	c := &cloudflareClient{Cli: client}
	return c
}
