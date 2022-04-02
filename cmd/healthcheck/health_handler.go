package main

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct {
	success  []byte
	error    []byte
	geo      Fragile
	fragiles []Fragile
}

func MakeHealthHandler(
	namespace string,
	fragiles []Fragile,
	geo Fragile) (*HealthHandler, error) {
	successString, _ := json.Marshal(map[string]string{
		"status":    "success",
		"code":      "200",
		"namespace": namespace,
	})
	errorString, _ := json.Marshal(map[string]string{
		"status":    "error",
		"code":      "500",
		"namespace": namespace,
	})
	return &HealthHandler{
		success:  successString,
		error:    errorString,
		geo:      geo,
		fragiles: fragiles,
	}, nil
}

func (hh *HealthHandler) isOk() bool {
	if hh.geo != nil && !hh.geo.IsOk() {
		return true
	}
	isOk := true
	for _, fragile := range hh.fragiles {
		isOk = isOk && fragile.IsOk()
	}
	return isOk
}

func (hh *HealthHandler) getCurrentResponse() ([]byte, int) {
	if hh.isOk() {
		return hh.success, http.StatusOK
	} else {
		return hh.error, http.StatusInternalServerError
	}
}

func (hh *HealthHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if req.Header.Get("Accept") != "application/json" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	response, statusCode := hh.getCurrentResponse()
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(response)
}
