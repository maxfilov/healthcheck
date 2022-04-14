package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimpleService_Check_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(`{"status": "UP"}`))
		if err != nil {
			t.Errorf("Unexpected error on response writing: %s", err)
		}
	}))
	defer server.Close()
	service, err := MakeSimpleService(server.URL, server.Client())
	err = service.Check()
	if err != nil {
		t.Errorf("Unexpected error: '%s'", err.Error())
	}
}

func TestSimpleService_Check_CanNotMakeRequest(t *testing.T) {
	service, err := MakeSimpleService(
		// I hope that this URL will be actually invalid
		"%$",
		&http.Client{})
	if service != nil {
		t.Errorf("Service should be nil when URL is invalid")
	}
	if err == nil {
		t.Errorf("Error is expected when URL is invalid")
	}
	if err.Error() != `can not make request object: parse "%$": invalid URL escape "%$"` {
		t.Errorf("Unexpected error message '%s'", err.Error())
	}
}

func TestSimpleService_Check_NoServer(t *testing.T) {
	service, err := MakeSimpleService("http://__I_invalid__url", &http.Client{})
	err = service.Check()
	if err == nil {
		t.Errorf("Error is expected when server is not available")
	}
}

func TestSimpleService_Check_Not200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusAlreadyReported)
		_, err := rw.Write([]byte(`{"status": "UP"}`))
		if err != nil {
			t.Errorf("Unexpected error on response writing: %s", err)
		}
	}))
	defer server.Close()
	service, err := MakeSimpleService(server.URL, server.Client())
	err = service.Check()
	if err == nil {
		t.Errorf("Error is expected when server responds not 200 OK")
	}
	if err.Error() != fmt.Sprintf("received '208 Already Reported' from '%s'", server.URL) {
		t.Errorf("Unexpected error received: '%s'", err.Error())
	}
}

func TestSimpleService_Print(t *testing.T) {
	service, _ := MakeSimpleService("http://someurl", &http.Client{})
	if service.Print() != "service at 'http://someurl'" {
		t.Errorf("Unexpected service name: '%s'", service.Print())
	}
}
