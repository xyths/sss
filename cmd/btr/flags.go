package main

import "github.com/urfave/cli/v2"

var (
	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "config.json",
		Usage:   "Read config from `FILE`",
	}
	InputFlag = &cli.StringFlag{
		Name:    "input",
		Aliases: []string{"i"},
		Value:   "input.csv",
		Usage:   "transfer data in `FILE`",
	}
)
