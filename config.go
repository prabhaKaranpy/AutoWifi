package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	PrimaryInterface   string `json:"primary_interface"`
	BackupSSID         string `json:"backup_ssid"`
	PingTarget         string `json:"ping_target"`
	LatencyThresholdMs int    `json:"latency_threshold_ms"`
	CheckIntervalSec   int    `json:"check_interval_sec"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
