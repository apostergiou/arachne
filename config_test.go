package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSetupConfig(t *testing.T) {
	listen := "0.0.0.0:53"
	upstream := "1.1.1.1:853"
	network := "tcp"

	expectedCfg := new(Config)
	expectedCfg.Listen = listen
	expectedCfg.Upstream = upstream
	expectedCfg.Network = network

	actualCfg, err := SetupConfig(listen, upstream, network)
	if err != nil {
		t.Fatalf("Expected to not get any errors, got %v", err)
	}

	if !cmp.Equal(expectedCfg, actualCfg) {
		t.Errorf("Expeced config to equal %v, got %v", expectedCfg, actualCfg)
	}
}

func TestMissingConfig(t *testing.T) {
	listen := "0.0.0.0:53"
	upstream := ""
	network := "tcp"

	expected := "upstream cannot be empty"
	_, actual := SetupConfig(listen, upstream, network)

	if actual.Error() != expected {
		t.Errorf("Expected %s, got %s.", expected, actual)
	}
}
