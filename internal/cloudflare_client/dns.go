package cloudflare_client

import (
	"context"
	"os"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/option"
)

type cloudflareClient struct {
	Cli *cloudflare.Client
}

// New : Creates a new docker client
func New() *cloudflareClient {
	client := cloudflare.NewClient(
		option.WithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN")),
	)
	c := &cloudflareClient{Cli: client}
	return c
}

func (c *cloudflareClient) NewDNSRecord(ctx context.Context) error {
	_, err := c.Cli.DNS.Records.New(ctx, dns.RecordNewParams{
		ZoneID: cloudflare.F(os.Getenv("CLOUDFLARE_ZONE_ID")),
		Body: dns.ARecordParam{
			Name:    cloudflare.F(os.Getenv("WEBSITE_DOMAIN")),
			TTL:     cloudflare.F(dns.TTL1), // automatic
			Type:    cloudflare.F(dns.ARecordTypeA),
			Content: cloudflare.F(os.Getenv("CLOUD_IP")),
			Proxied: cloudflare.F(true), // proxy it
		},
	})
	if err != nil {
		return err
	}
	return nil
}
