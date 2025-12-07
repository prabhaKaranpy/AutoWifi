# ⚡ AutoWifi

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-gray)](https://github.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**A High-Availability Network Daemon that automatically reroutes system traffic to a backup interface (Hotspot) when the primary connection (WiFi) degrades.**

Built to solve the "Codeforces/CodeChef Lag" problem, ensuring 99.9% uptime for competitive programming submissions.

---

## 📖 The Problem
During competitive programming contests (Codeforces, LeetCode), a stable internet connection is critical. Standard OS behavior only switches networks if the WiFi **disconnects** completely. However, it fails to account for **Packet Loss**, **Jitter**, or **High Latency**, leaving the user stuck on a "zombie" connection that cannot load pages or submit code.

## 🚀 The Solution
**NetFailover** acts as a Layer-3 Network Sentinel. It continuously monitors the *quality* (not just connectivity) of your primary interface. If latency spikes or packet loss is detected, it programmatically manipulates the OS Routing Table to switch the Default Gateway to your backup interface in sub-seconds.

### ✨ New Features (v1.0.0)
As seen in the latest release:

* **🧙‍♂️ Smart Setup Wizard:** * **Auto-Discovery:** Automatically detects your active WiFi interface via OS syscalls.
  * **Zero-Config:** No more manual JSON editing. The app prompts you once for your backup interface (e.g., "Ethernet" or "iPhone USB") and saves it.
* **🎨 System Tray Integration:** * **Real-time Feedback:** A non-intrusive icon in your taskbar showing live health.
  * 🟢 **Green:** Stable Connection (Low Latency).
  * 🟡 **Yellow:** Degraded Service (Jitter detected).
  * 🔴 **Red:** Failover Active / Primary Down.
* **🤫 Silent Mode:** * Runs as a background daemon with low memory footprint (~5MB). No popping command terminals.

---

## 🛠️ Architecture & Tech Stack

The system follows a **Monitor-Decide-Act** architectural pattern:

| Component | Tech | Description |
| :--- | :--- | :--- |
| **Sentinel** | `prometheus/pro-bing` | Binds concurrent ICMP sockets to specific interfaces to measure RTT and Packet Loss independently. |
| **Arbiter** | **Go (Golang)** | Implements hysteresis logic (e.g., "Only switch if Loss > 10% for 5 seconds") to prevent network flapping. |
| **Switchman** | `syscall` / `netsh` | Directly interfaces with the Windows IP Helper API (or Linux Netlink) to modify Route Metrics. |
| **UI Layer** | `systray` | Cross-platform system tray implementation for status visualization. |

---

## 💻 Installation & Usage

### Option 1: Download Binary
Go to the [Release](https://github.com/prabhaKaranpy/AutoWifi/releases/) page and download `autowifi.exe`.

