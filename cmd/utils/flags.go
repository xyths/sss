package utils

import "gopkg.in/urfave/cli.v2"

var (
	FileFlag = &cli.StringFlag{
		Name:    "file",
		Aliases: []string{"f"},
		Value:   "stake.txt",
		Usage:   "Read stake list from `FILE`",
	}
	StakeFlag = &cli.StringFlag{
		Name:    "stake",
		Aliases: []string{"s"},
		Value:   "stake.json",
		Usage:   "Read stake detail info from `STAKE`",
	}
)
