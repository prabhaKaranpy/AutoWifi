package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func runSetup() {
	fmt.Println("Running AutoWifi Setup...")

	// 1. Check if config exists
	if _, err := os.Stat("config.json"); err == nil {
		return // Config exists, skip setup
	}

	// 2. Auto-detect Primary Interface (WiFi)
	primaryIface := detectPrimaryInterface()
	fmt.Printf("Detected Primary Interface: %s\n", primaryIface)

	// 3. Prompt for Backup SSID
	backupSSID := promptBackupSSID()
	fmt.Printf("User Selected Backup SSID: %s\n", backupSSID)

	if backupSSID == "" {
		// Fallback default
		backupSSID = "MyBackupWifi"
	}

	// 4. Create and Save Config
	newConfig := Config{
		PrimaryInterface:   primaryIface,
		BackupSSID:         backupSSID,
		PingTarget:         "8.8.8.8",
		LatencyThresholdMs: 200,
		CheckIntervalSec:   5,
	}

	saveConfig(newConfig)
}

func detectPrimaryInterface() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "Wi-Fi" // Safe default
	}

	// Priority heuristics
	for _, iface := range interfaces {
		if strings.Contains(strings.ToLower(iface.Name), "wifi") ||
			strings.Contains(strings.ToLower(iface.Name), "wi-fi") ||
			strings.Contains(strings.ToLower(iface.Name), "wireless") {
			return iface.Name
		}
	}
	return "Wi-Fi"
}

func promptBackupSSID() string {
	// Use PowerShell to show a simple InputBox
	// Requires: Add-Type -AssemblyName Microsoft.VisualBasic
	script := `
	Add-Type -AssemblyName Microsoft.VisualBasic
	$result = [Microsoft.VisualBasic.Interaction]::InputBox("Enter the exact NAME (SSID) of your Backup Wi-Fi Network. (Note: You must have connected to it before so Windows knows the password).", "AutoWifi Setup", "MyHotspot")
	Write-Output $result
	`

	cmd := exec.Command("powershell", "-Command", script)
	// We want to HIDE the console window of powershell too, but we need the GUI input box
	// We will try running it. If we hide the console, the input box might still show.
	// But standard exec should be fine for a popup.

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "MyHotspot"
	}

	return strings.TrimSpace(string(output))
}

func saveConfig(config Config) {
	file, _ := os.Create("config.json")
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(config)
}
