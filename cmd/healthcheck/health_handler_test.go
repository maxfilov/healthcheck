package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fragileStub struct {
	isOk bool
}

func (s *fragileStub) IsOk() bool { return s.isOk }

func TestMakeHealthHandler(t *testing.T) {
	services := []Fragile{&fragileStub{isOk: true}}
	handler, err := MakeHealthHandler("", services, nil)
	if err != nil {
		t.Error("Unexpected error occurred during health handler initialization")
	}
	if string(handler.success) != `{"code":"200","namespace":"","status":"success"}` {
		t.Errorf("Unexpected value for success status response: '%s'", string(handler.success))
	}
	if string(handler.error) != `{"code":"500","namespace":"","status":"error"}` {
		t.Errorf("Unexpected value for error status response: '%s'", string(handler.error))
	}
}

func TestHealthHandler_ServeHTTP_ValidRequest_ServicesOk(t *testing.T) {
	handler, _ := MakeHealthHandler("x-namespace-x", []Fragile{&fragileStub{isOk: true}}, nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, &http.Request{
		Method: http.MethodGet,
		Header: http.Header{
			http.CanonicalHeaderKey("content-type"): []string{"application/json"},
			http.CanonicalHeaderKey("accept"):       []string{"application/json"},
		},
	})
	if rr.Code != http.StatusOK {
		t.Errorf("Unexpected status code %d in the response, expected=%d", rr.Code, http.StatusOK)
	}
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Error("Unexpected Content-Type in the response")
	}
	var resp map[string]string
	err := json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Errorf("Unexpected error during parsion of the response body: '%s'", err.Error())
	}
	if resp["code"] != "200" {
		t.Errorf("Unexpected code '%s' in the body, expected='%s'", resp["code"], "200")
	}
	if resp["status"] != "success" {
		t.Errorf("Unexpected status '%s' in the body, expected='%s'", resp["status"], "success")
	}
	if resp["namespace"] != "x-namespace-x" {
		t.Errorf("Unexpected namespace '%s' in the body, expected='%s'", resp["namespace"], "x-namespace-x")
	}
}

func TestHealthHandler_ServeHTTP_ValidRequest_ServicesNotOk(t *testing.T) {
	handler, _ := MakeHealthHandler("x-namespace-x", []Fragile{&fragileStub{isOk: false}}, nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, &http.Request{
		Method: http.MethodGet,
		Header: http.Header{
			http.CanonicalHeaderKey("content-type"): []string{"application/json"},
			http.CanonicalHeaderKey("accept"):       []string{"application/json"},
		},
	})
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected status code %d in the response, expected=%d", rr.Code, http.StatusInternalServerError)
	}
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Error("Unexpected Content-Type in the response")
	}
	var resp map[string]string
	err := json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Errorf("Unexpected error during parsion of the response body: '%s'", err.Error())
	}
	if resp["code"] != "500" {
		t.Errorf("Unexpected code '%s' in the body, expected='%s'", resp["code"], "500")
	}
	if resp["status"] != "error" {
		t.Errorf("Unexpected status '%s' in the body, expected='%s'", resp["status"], "error")
	}
	if resp["namespace"] != "x-namespace-x" {
		t.Errorf("Unexpected namespace '%s' in the body, expected='%s'", resp["namespace"], "x-namespace-x")
	}
}

func TestHealthHandler_ServeHTTP_InvalidMethod(t *testing.T) {
	handler, _ := MakeHealthHandler("x-namespace-x", []Fragile{}, nil)
	rr := httptest.NewRecorder()

	for _, method := range []string{
		http.MethodPut,
		http.MethodPost,
		http.MethodHead,
		http.MethodTrace,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions} {
		handler.ServeHTTP(rr, &http.Request{
			Method: method,
			Header: http.Header{
				http.CanonicalHeaderKey("content-type"): []string{"application/json"},
				http.CanonicalHeaderKey("accept"):       []string{"application/json"},
			},
		})
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Unexpected status code %d in the response, expected=%d", rr.Code, http.StatusMethodNotAllowed)
		}
	}
}

func TestHealthHandler_ServeHTTP_InvalidAccept(t *testing.T) {
	handler, _ := MakeHealthHandler("x-namespace-x", []Fragile{}, nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, &http.Request{
		Method: http.MethodGet,
		Header: http.Header{
			http.CanonicalHeaderKey("content-type"): []string{"application/json"},
			http.CanonicalHeaderKey("accept"):       []string{"application/xml"},
		},
	})
	if rr.Code != http.StatusNotAcceptable {
		t.Errorf("Unexpected status code %d in the response, expected=%d", rr.Code, http.StatusMethodNotAllowed)
	}
}
