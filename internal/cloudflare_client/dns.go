package cloudflare_client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/option"
)

type DnsRecord struct {
	Id         string `json:"id"`
	RecordType string `json:"type"`
}

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
			Proxied: cloudflare.F(true),         // proxy it
			Comment: cloudflare.F("Turtle CLI"), // Use comment to identify records created by turtle cli
		},
	})
	if err != nil {
		return err
	}
	fmt.Println("New DNS Record Created Successfully")
	return nil
}

// GetDNSRecord : returns a DNSRecord struct if a turtle created record exists else nil
func (c *cloudflareClient) GetDNSRecord(ctx context.Context) (*DnsRecord, error) {
	//_, err := c.Cli.DNS.Records.Get(ctx)
	resp, err := c.Cli.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(os.Getenv("CLOUDFLARE_ZONE_ID")),
		// Use comment to identify records created by turtle cli
		Comment: cloudflare.F(dns.RecordListParamsComment{
			Contains: cloudflare.F("Turtle CLI"),
		}),
	})
	if err != nil {
		return nil, err
	}
	// No DNS records exist that were created by turtle
	if len(resp.Result) == 0 {
		return nil, nil
	}
	d := DnsRecord{}
	err = json.Unmarshal([]byte(resp.Result[0].JSON.RawJSON()), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (c *cloudflareClient) DeleteDNSRecord(d *DnsRecord, ctx context.Context) error {
	_, err := c.Cli.DNS.Records.Delete(ctx, d.Id, dns.RecordDeleteParams{
		ZoneID: cloudflare.F(os.Getenv("CLOUDFLARE_ZONE_ID")),
	})
	return err
}
