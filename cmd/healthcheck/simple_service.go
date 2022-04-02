package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type SimpleService struct {
	endpoint string
	client   *http.Client
	name     string
	request  *http.Request
}

func MakeSimpleService(endpoint string, client *http.Client) (Service, error) {
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("can not create request object: %s", err.Error())
	}
	request.Header.Set("Accept", "application/json")
	return &SimpleService{
		endpoint: endpoint,
		client:   client,
		name:     fmt.Sprintf("service at '%s'", endpoint),
		request:  request,
	}, nil
}

func (srv *SimpleService) Print() string {
	return srv.name
}

func (srv *SimpleService) Check() error {
	resp, err := srv.client.Do(srv.request)
	if err != nil {
		return fmt.Errorf("can not make request: %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("could not close the response: %s", err.Error())
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received '%s' from '%s'", resp.Status, srv.endpoint)
	}
	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response from '%s': '%s'", srv.endpoint, err.Error())
	}
	return nil
}
