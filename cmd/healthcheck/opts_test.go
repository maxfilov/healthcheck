package main

import "testing"

func TestParseOpts_ValidOpts(t *testing.T) {
	for _, args := range [][]string{
		{"--config=hello"},
		{"--config", "hello"},
		{"-c", "hello"}} {
		opts, err := ParseOpts(args)
		if err != nil {
			t.Errorf("Unexpected error occured during options parsing: '%s'", err.Error())
		}
		if opts.ConfigPath != "hello" {
			t.Errorf("Unexpected 'config' option value '%s', expected='%s'", opts.ConfigPath, "hello")
		}
	}
}

func TestParseOpts_InvalidOpts(t *testing.T) {
	opts, err := ParseOpts([]string{"--random-invalid-option=random-value"})
	if err == nil {
		t.Errorf("Expected an error, got nothing")
	}
	if opts != nil {
		t.Errorf("No opts should be produced from invalid input")
	}
}
