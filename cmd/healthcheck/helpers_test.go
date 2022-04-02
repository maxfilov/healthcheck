package main

import "fmt"

type ServiceStub struct {
	Err error
}

func (s *ServiceStub) Check() error {
	return s.Err
}

func (s *ServiceStub) Print() string {
	return fmt.Sprintf("service stub")
}
