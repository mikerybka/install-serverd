package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Install serverd
	cmd := exec.Command("go", "install", "github.com/mikerybka/serverd@latest")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	// Locate executable
	cmd = exec.Command("which", "serverd")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	path := strings.TrimSpace(string(out))

	// Get tailsacale IP
	cmd = exec.Command("tailscale", "ip", "--4")
	out, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	ip := strings.TrimSpace(string(out))

	// Write systemd unit
	err = os.WriteFile("/etc/systemd/system/serverd.service", fmt.Appendf(nil, `[Unit]
Description=serverd
After=network.target

[Service]
ExecStart=%s
Environment="LISTEN_ADDR=%s:4321"
Restart=on-failure

[Install]
WantedBy=multi-user.target
`, path, ip), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Start and enable service
	cmd = exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	cmd = exec.Command("systemctl", "enable", "--now", "serverd")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
