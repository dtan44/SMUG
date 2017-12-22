package service

import (
	"strconv"
	"testing"
)

func setupServiceDiscovery() {
	serviceMap = make(map[string]string)
	mapPrint = mapToString
}

func TestDiscoveryList(t *testing.T) {
	setupServiceDiscovery()
	var ds DiscoveryService
	serviceMap["test"] = "test"

	// invalid URL
	res := ds.List()
	if len(res) != 1 {
		t.Errorf("service returned unexpected error: got %v want %v",
			strconv.Itoa(len(res)), strconv.Itoa(0))
	}
}
