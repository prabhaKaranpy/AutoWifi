package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/getlantern/systray"
)

var (
	config  *Config
	logFile *os.File
)

func logMsg(format string, v ...interface{}) {
	if logFile != nil {
		fmt.Fprintf(logFile, format+"\n", v...)
	}
}

func main() {
	var err error
	logFile, err = os.OpenFile("autowifi.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer logFile.Close()
	}
	logMsg("App Started")

	runSetup()

	config, err = LoadConfig("config.json")
	if err != nil {
		logMsg("Error loading config: %v", err)
		return
	}
	logMsg("Loaded Config: %v", config)

	systray.Run(onReady, onExit)
}

func onReady() {
	logMsg("onReady called")

	systray.SetIcon(iconGreen)
	systray.SetTitle("AutoWifi")
	systray.SetTooltip("AutoWifi: Connected")
	logMsg("Icon set to Green")

	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			monitorNetwork()
			time.Sleep(time.Duration(config.CheckIntervalSec) * time.Second)
		}
	}()

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// Cleanup here
}

func monitorNetwork() {
	// For SSID switching, we mainly care about the primary WiFi interface status
	// We want to verify we are connected to SOMETHING

	// Check Primary Connection Quality
	// Note: In single adapter mode, primaryIP is effectively just the adapter's IP.
	// We need to fetch it dynamically.

	iface, err := net.InterfaceByName(config.PrimaryInterface)
	if err != nil {
		systray.SetIcon(iconRed)
		systray.SetTooltip("WiFi Adapter Not Found")
		return
	}

	ip := getIPv4(iface)
	if ip == "" {
		systray.SetIcon(iconRed)
		systray.SetTooltip("WiFi Disconnected")
		// If disconnected, maybe try to connect to backup?
		// But for now, let's assume we only switch on poor connection,
		// or if completely disconnected we try backup.

		fmt.Println("[!] Disconnected. Attempting switch to Backup SSID...")
		switchSSID(config.BackupSSID, config.PrimaryInterface)
		return
	}

	// Check Latency
	latency, _, err := ping(ip, config.PingTarget)
	if err != nil {
		systray.SetIcon(iconRed)
		systray.SetTooltip("Ping Failed")

		fmt.Println("[!] Ping Failed. Attempting switch to Backup SSID...")
		switchSSID(config.BackupSSID, config.PrimaryInterface)
		return
	}

	if latency > config.LatencyThresholdMs {
		systray.SetIcon(iconYellow)
		systray.SetTooltip(fmt.Sprintf("Unstable: %dms", latency))
		fmt.Printf("[!] High Latency: %dms. Attempting switch to Backup SSID...\n", latency)

		switchSSID(config.BackupSSID, config.PrimaryInterface)
	} else {
		systray.SetIcon(iconGreen)
		systray.SetTooltip(fmt.Sprintf("Good: %dms", latency))
	}
}

func switchSSID(ssid string, interfaceName string) {
	// Command: netsh wlan connect name="SSID" interface="Interface"
	fmt.Printf("Executing Switch to %s...\n", ssid)

	// HIDE WINDOW
	cmd := exec.Command("netsh", "wlan", "connect", fmt.Sprintf("name=%s", ssid), fmt.Sprintf("interface=%s", interfaceName))
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Switch Failed: %v\n", err)
	} else {
		fmt.Println("Switch Command Sent. Waiting for reconnection...")
		// Wait a bit to avoid rapid flapping
		time.Sleep(15 * time.Second)
	}
}

func getIPv4(iface *net.Interface) string {
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}
		if ip.To4() != nil && !ip.IsLinkLocalUnicast() {
			return ip.String()
		}
	}
	return ""
}

func ping(sourceIP, targetIP string) (int, float64, error) {
	cmd := exec.Command("ping", "-n", "1", "-w", "1000", "-S", sourceIP, targetIP)

	// CREATE_NO_WINDOW = 0x08000000
	// This ensures the window is absolutely hidden
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 100, err
	}

	return parsePingOutput(string(output))
}

func parsePingOutput(output string) (int, float64, error) {
	lines := strings.Split(output, "\n")
	var avgLatency int
	var packetLoss float64

	for _, line := range lines {
		if strings.Contains(line, "Average =") {
			parts := strings.Split(line, "Average =")
			if len(parts) > 1 {
				latencyStr := strings.TrimSpace(parts[1])
				latencyStr = strings.TrimSuffix(latencyStr, "ms")
				fmt.Sscanf(latencyStr, "%d", &avgLatency)
			}
		}
		// Simplified loss parsing for single ping
		if strings.Contains(line, "Lost = 1") {
			packetLoss = 100
		}
	}

	return avgLatency, packetLoss, nil
}

func getGateway(ifaceName string) string {
	// ... (Same as before, omitted for brevity if not used immediately)
	return ""
}
