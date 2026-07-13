package cloudflare_client

import (
	"context"
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

func StartCloudflared() {
	// linux or windows

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
	return nil
}

func isCloudflaredUp() bool {
	return false
}
