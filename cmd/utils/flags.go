package utils

import "gopkg.in/urfave/cli.v2"

var (
	AppendConfigFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "config.json",
		Usage:   "Read append config from `FILE`",
	}
	ShareListFlag = &cli.StringFlag{
		Name:    "list",
		Aliases: []string{"l"},
		Value:   "shares.csv",
		Usage:   "Read share list from `FILE`",
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
		Value:   "shares.csv",
		Usage:   "Output stake profit snapshot to `CSV`",
	}
	FilterCompanyFlag = &cli.StringFlag{
		Name:  "company",
		Usage: "Only print `company`'s stake",
	}
	MailConfigFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "config.json",
		Usage:   "Read mail config from `FILE`",
	}
	DateFlag = &cli.StringFlag{
		Name:    "date",
		Aliases: []string{"d"},
		Usage:   "Report `DATE`, like 20190916",
	}
)
