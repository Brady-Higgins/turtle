package cloudflare_client

import (
	"fmt"
	"os/exec"
)

func RunCloudflared(tunnelName string) error {
	_, err := exec.Command("cloudflared", "tunnel", "run", tunnelName).Output()
	if err != nil {
		return err
	}
	return nil
}

func CreateTunnelDNSRecord(tunnelName string, hostName string) error {
	_, err := exec.Command("cloudflared", "tunnel", "route", "dns", tunnelName, hostName).Output()
	if err != nil {
		return err
	}
	fmt.Println("Tunnel DNS Record Created Successfully")
	return nil
}

func isCloudflaredUp() bool {
	return false
}
