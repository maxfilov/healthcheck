package main

type Lifecycle interface {
	StartAsync()
	Shutdown() error
	AwaitShutdown() error
}
