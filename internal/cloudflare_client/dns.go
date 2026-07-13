package cloudflare_client

import (
	"context"
	"encoding/json"
	"os"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/option"
)

type DnsRecord struct {
	Id         string `json:"id"`
	RecordType string `json:"type"`
}

const dnsRecordComment string = "Turtle CLI"

type CloudflareClient struct {
	Cli *cloudflare.Client
}

// New : Creates a new docker client
func New() *CloudflareClient {
	client := cloudflare.NewClient(
		option.WithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN")),
	)
	c := &CloudflareClient{Cli: client}
	return c
}

func (c *CloudflareClient) NewDNSRecord(ip string, ctx context.Context) error {
	_, err := c.Cli.DNS.Records.New(ctx, dns.RecordNewParams{
		ZoneID: cloudflare.F(os.Getenv("CLOUDFLARE_ZONE_ID")),
		Body: dns.ARecordParam{
			Name:    cloudflare.F(os.Getenv("WEBSITE_DOMAIN")),
			TTL:     cloudflare.F(dns.TTL1), // automatic
			Type:    cloudflare.F(dns.ARecordTypeA),
			Content: cloudflare.F(ip),
			Proxied: cloudflare.F(true),             // proxy it
			Comment: cloudflare.F(dnsRecordComment), // Use comment to identify records created by turtle cli
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// GetDNSRecord : returns a DNSRecord struct if a turtle created record exists else nil
// dns.RecordListParamsTypeA for A records and dns.RecordListParamsTypeCNAME for tunnel record
// commented : true if created by turtle CLI. Uses magic comment
func (c *CloudflareClient) GetDNSRecord(recordType dns.RecordListParamsType, commented bool, ctx context.Context) (*DnsRecord, error) {
	params := dns.RecordListParams{
		ZoneID: cloudflare.F(os.Getenv("CLOUDFLARE_ZONE_ID")),
		Type:   cloudflare.F(recordType),
		// contains no comment
		Comment: cloudflare.F(dns.RecordListParamsComment{
			Absent: cloudflare.F("comment"),
		}),
	}
	if commented {
		// Use comment to identify records created by turtle cli
		params.Comment = cloudflare.F(dns.RecordListParamsComment{
			Contains: cloudflare.F(dnsRecordComment),
		})
	}
	resp, err := c.Cli.DNS.Records.List(ctx, params)
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

func (c *CloudflareClient) DeleteDNSRecord(d *DnsRecord, ctx context.Context) error {
	_, err := c.Cli.DNS.Records.Delete(ctx, d.Id, dns.RecordDeleteParams{
		ZoneID: cloudflare.F(os.Getenv("CLOUDFLARE_ZONE_ID")),
	})
	return err
}

func (c *CloudflareClient) CommentDNSRecord(d *DnsRecord, ctx context.Context) error {
	_, err := c.Cli.DNS.Records.Edit(ctx, d.Id, dns.RecordEditParams{
		ZoneID: cloudflare.F(os.Getenv("CLOUDFLARE_ZONE_ID")),
		Body: dns.RecordEditParamsBody{
			Comment: cloudflare.F(dnsRecordComment),
		},
	})

	return err
}
