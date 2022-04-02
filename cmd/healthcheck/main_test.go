package main

import "testing"

func TestMain_assembleConfiguration(t *testing.T) {
	configuration, err := assembleConfiguration()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if configuration.Server.Port != 8080 {
		t.Errorf("Unexpected server.port: %d", configuration.Server.Port)
	}
}
