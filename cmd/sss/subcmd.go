package main

import (
	"encoding/json"
	"github.com/xyths/sss/cmd/utils"
	"github.com/xyths/sss/stake"
	"gopkg.in/urfave/cli.v2"
	"io"
	"log"
	"math/big"
	"os"
)

var (
	sumCommand = &cli.Command{
		Action:  sum,
		Name:    "sum",
		Aliases: []string{"s"},
		Usage:   "Sum staking SERO from all stake",
		Flags: []cli.Flag{
			utils.StakeFlag,
		},
	}
)

func sum(ctx *cli.Context) (err error) {
	filename := ctx.String(utils.StakeFlag.Name)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	stakeList := json.NewDecoder(file)

	balanceXH := big.NewFloat(0)
	balanceJRXC := big.NewFloat(0)
	for {
		var sd stake.StakeDetail
		if err := stakeList.Decode(&sd); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if remain, err := stake.Sum(sd); err == nil {
			log.Printf("%s\tid %s remain sero: %v\n", sd.Company, sd.Id, remain)

			switch sd.Company {
			case "xh":
				balanceXH.Add(balanceXH, remain)
			case "jrxc":
				balanceJRXC.Add(balanceJRXC, remain)
			}
		} else {
			log.Printf("id: %s %s\n", sd.Id, err)
		}
	}
	log.Printf("xh balance is: %v", balanceXH)
	log.Printf("jrxc balance is: %v", balanceJRXC)
	return
}
