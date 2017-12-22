package service

import "testing"

func setupServiceDiscovery() {
	serviceMap = make(map[string]string)
	mapPrint = mapToString
}

func TestDiscoveryList(t *testing.T) {
	setupServiceDiscovery()
	var ds DiscoveryService

	testText := `{}`

	// invalid URL
	res := ds.List()
	if res != testText {
		t.Errorf("service returned unexpected error: got %v want %v",
			res, testText)
	}
}
