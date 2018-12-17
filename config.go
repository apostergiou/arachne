package main

import "errors"

// Config holds the configuration values.
type Config struct {
	Listen   string
	Upstream string
	Network  string
}

// SetupConfig accepts the listening address, the upstream address and a
// network type (tcp or udp).
func SetupConfig(listen string, upstream string, network string) (*Config, error) {
	cfg := new(Config)

	if listen == "" {
		return cfg, errors.New("listen cannot be empty")
	}
	cfg.Listen = listen

	cfg.Upstream = upstream
	if upstream == "" {
		return cfg, errors.New("upstream cannot be empty")
	}

	cfg.Network = network
	if network == "" {
		return cfg, errors.New("network cannot be empty")
	}

	return cfg, nil
}
