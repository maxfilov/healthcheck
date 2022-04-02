package main

type FragileService interface {
	IsOk() bool
	Check() error
	Print() string
}
