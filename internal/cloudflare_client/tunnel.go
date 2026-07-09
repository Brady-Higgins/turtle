package cloudflare_client

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func CreateCloudflaredCommand(ctx context.Context, tunnelName string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "cloudflared", "tunnel", "run", tunnelName)
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	return cmd
}

func RunCloudflared(cmd *exec.Cmd) {
	_ = cmd.Start()
	return
}

func StopCloudflared(cmd *exec.Cmd, cancel context.CancelFunc) error {
	cancel()
	err := cmd.Wait()
	//err := cmd.Process.Signal(syscall.SIGTERM)
	return err
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
