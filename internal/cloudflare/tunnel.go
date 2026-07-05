package cloudflare

import (
	"os/exec"
)

func RunCloudflared(tunnelName string) error {
	_, err := exec.Command("cloudflared", "tunnel", "run", tunnelName).Output()
	if err != nil {
		return err
	}
	return nil
}
