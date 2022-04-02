package main

type Service interface {
	Check() error
	Print() string
}
