package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Brady-Higgins/turtle/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func prompt(r *bufio.Reader, label string) string {
	fmt.Printf("%s: ", label)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func promptSecret(label string) string {
	fmt.Printf("%s: ", label)
	b, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return strings.TrimSpace(string(b))
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "set up config for turtle",
	Long:  "set up config for turtle",
	Run: func(cmd *cobra.Command, args []string) {
		r := bufio.NewReader(os.Stdin)
		cfg := &config.Config{
			CloudflareAPIToken:   promptSecret("Cloudflare API Token"),
			CloudflareAccountID:  prompt(r, "Cloudflare Account ID"),
			CloudflareZoneID:     prompt(r, "Cloudflare Zone ID"),
			CloudflareTunnelName: prompt(r, "Cloudflare Tunnel Name"),
			HostName:             prompt(r, "Host Name"),
			AWSAccessKeyID:       prompt(r, "AWS Access Key ID"),
			AWSSecretAccessKey:   promptSecret("AWS Secret Access Key"),
		}
		config.SetupConfig()
		err := config.WriteConfig(cfg)
		if err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
}
