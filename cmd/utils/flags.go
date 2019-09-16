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
	CsvFlag = &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Value:   "stake.csv",
		Usage:   "Output stake profit snapshot to `CSV`",
	}
	FilterCompanyFlag = &cli.StringFlag{
		Name:    "company",
		Aliases: []string{"com"},
		Value:   "",
		Usage:   "Only print `company`'s stake",
	}
)
