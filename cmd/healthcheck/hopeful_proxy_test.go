package main

import (
	"errors"
	"testing"
)

func TestMakeWatchfulDecorator(t *testing.T) {
	decorator := MakeHopefulProxy(&ServiceStub{}, 3)
	if decorator == nil {
		t.Errorf("Decorator should not be nil")
	}
}

func TestWatchfulDecorator_Lifecycle(t *testing.T) {
	service := ServiceStub{}
	threshold := 3
	decorator := MakeHopefulProxy(&service, threshold)
	for i := 0; i < 20; i++ {
		err := decorator.Check()
		if err != nil {
			t.Errorf("Unexpected error when backend does not return error")
		}
		if !decorator.IsOk() {
			t.Errorf("Unexpected state of the decorator")
		}
	}
	service.Err = errors.New("error")
	for i := 0; i < threshold; i++ {
		err := decorator.Check()
		if err == nil {
			t.Errorf("Unexpected nil when service returns error")
		}
		if !decorator.IsOk() {
			t.Errorf("Unexpected state of the decorator")
		}
	}
	err := decorator.Check()
	if err == nil {
		t.Errorf("Unexpected nil when service returns error")
	}
	if decorator.IsOk() {
		t.Errorf("Unexpected state of the decorator after the threshold has been exceeded")
	}
	service.Err = nil
	for i := 0; i < 20; i++ {
		err := decorator.Check()
		if err != nil {
			t.Errorf("Unexpected error when backend does not return error")
		}
		if !decorator.IsOk() {
			t.Errorf("Unexpected state of the decorator")
		}
	}
}

func TestWatchfulDecorator_Print(t *testing.T) {
	s := MakeHopefulProxy(&ServiceStub{}, 3).Print()
	if s != "watchful decorator for service stub" {
		t.Errorf("Unexpected description of the service")
	}
}
