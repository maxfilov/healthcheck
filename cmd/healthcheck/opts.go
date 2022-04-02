package main

import "github.com/jessevdk/go-flags"

type Opts struct {
	ConfigPath string `short:"c" long:"config" description:"Path to the configuration file"`
}

func ParseOpts(args []string) (*Opts, error) {
	var opts Opts
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return nil, err
	}
	return &opts, nil
}
