package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync/atomic"
)

type HopefulProxy struct {
	backend   Service
	counter   int
	threshold int
	isOk      *atomic.Value
}

func MakeHopefulProxy(backend Service, threshold int) *HopefulProxy {
	isOk := atomic.Value{}
	isOk.Store(true)
	return &HopefulProxy{
		backend:   backend,
		counter:   0,
		threshold: threshold,
		isOk:      &isOk,
	}
}

func (decor *HopefulProxy) Check() error {
	err := decor.backend.Check()
	if err != nil {
		decor.failed()
	} else {
		decor.succeeded()
	}
	return err
}

func (decor *HopefulProxy) IsOk() bool {
	return decor.isOk.Load().(bool)
}

func (decor *HopefulProxy) Print() string {
	return fmt.Sprintf("watchful decorator for %s", decor.backend.Print())
}

func (decor *HopefulProxy) failed() {
	if decor.counter >= decor.threshold {
		decor.isOk.Store(false)
		return
	}
	decor.counter += 1
}

func (decor *HopefulProxy) succeeded() {
	if !decor.isOk.Load().(bool) {
		logrus.Infof("%s is OK now", decor.backend.Print())
	}
	decor.isOk.Store(true)
	decor.counter = 0
}
