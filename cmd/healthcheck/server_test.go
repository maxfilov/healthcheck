package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestServer_Lifecycle(t *testing.T) {
	logger := log.New()
	logger.SetOutput(bytes.NewBufferString(""))
	server := MakeServer(8080, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte("hello"))
		if err != nil {
			t.Errorf("Unexpected error on response writing: %s", err)
		}
	}), logger)
	server.StartAsync()
	client := http.Client{}
	resp, err := client.Get("http://localhost:8080/health")
	if err != nil {
		t.Errorf("Unexpected non nil error: '%s'", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d, expected=%d", resp.StatusCode, http.StatusOK)
	}
	buf := strings.Builder{}
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		t.Errorf("Unexpected error on conversion: %s", err)
	}
	s := buf.String()
	if s != "hello" {
		t.Errorf("Unexpected body '%s', expected='%s'", s, "hello")
	}
	http.DefaultServeMux = new(http.ServeMux)
	err = server.Shutdown()
	if err != nil {
		t.Errorf("Unexpected error on server shutdown: %s", err)
	}
	err = server.AwaitShutdown()
	if err != nil {
		t.Errorf("Unexpected error on server shutdown awaiting: %s", err)
	}
}
