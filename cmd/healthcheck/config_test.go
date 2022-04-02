package main

import (
	"testing"
)

func TestConfig_AsJson(t *testing.T) {
	config := Config{}
	expected := `{
  "Pod": {
    "Namespace": ""
  },
  "Server": {
    "Port": 0
  },
  "Logging": {
    "Level": {
      "Root": ""
    }
  },
  "Schedule": {
    "Enabled": false,
    "Delay": 0
  },
  "ClientServices": {
    "Services": null
  },
  "Geo": null,
  "FailureThreshold": 0
}`
	jsonString := config.AsJson()
	if jsonString != expected {
		t.Errorf("Unexpected json: %s", jsonString)
	}
}
