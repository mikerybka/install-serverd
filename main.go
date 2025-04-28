package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Install serverd
	fmt.Printf("Installing serverd...")
	cmd := exec.Command("go", "install", "github.com/mikerybka/serverd@latest")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(" Done.")

	// Locate executable
	fmt.Printf("Locating exe...")
	cmd = exec.Command("which", "serverd")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	path := strings.TrimSpace(string(out))
	fmt.Println(" Done.")

	// Get tailsacale IP
	fmt.Printf("Getting tailscale IP...")
	cmd = exec.Command("tailscale", "ip", "--4")
	out, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	ip := strings.TrimSpace(string(out))
	fmt.Println(" Done.")

	// Write systemd unit
	fmt.Printf("Writing systemd unit...")
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
	fmt.Println(" Done.")

	// Start and enable service
	fmt.Printf("Reloading systemd...")
	cmd = exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(" Done.")
	fmt.Printf("Starting serverd...")
	cmd = exec.Command("systemctl", "enable", "--now", "serverd")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(" Done.")
}
