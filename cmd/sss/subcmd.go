package main

import (
	"encoding/json"
	"fmt"
	"github.com/xyths/sss/cmd/utils"
	"github.com/xyths/sss/stake"
	"gopkg.in/urfave/cli.v2"
	"io"
	"log"
	"os"
	"sort"
)

var (
	sumCommand = &cli.Command{
		Action:  sum,
		Name:    "sum",
		Aliases: []string{"s"},
		Usage:   "Sum staking SERO from all stake",
		Flags: []cli.Flag{
			utils.StakeFlag,
			utils.CsvFlag,
			utils.FilterCompanyFlag,
		},
	}
	snapshotCommand = &cli.Command{
		Action:  snapshot,
		Name:    "snapshot",
		Aliases: []string{"snap"},
		Usage:   "Snapshot staking status of all stake",
		Flags: []cli.Flag{
			utils.StakeFlag,
			utils.CsvFlag,
			utils.FilterCompanyFlag,
		},
	}
	profitCommand = &cli.Command{
		Action:  profit,
		Name:    "profit",
		Aliases: []string{"pf"},
		Usage:   "sum all profit of all stake",
		Flags: []cli.Flag{
			utils.StakeFlag,
			utils.CsvFlag,
			utils.FilterCompanyFlag,
		},
	}
)

func sum(ctx *cli.Context) (err error) {
	stakeList := readStakeList(ctx)

	var balance float64
	var results []stake.Result
	for _, sd := range stakeList {
		res := stake.Sum(sd)
		results = append(results, res)
		log.Printf("%s,%d,%d,%d,%f,%f,%f\n", res.Id,
			res.TotalShare, res.ReturnedShare, res.MortgageShare,
			res.TotalPrinciple, res.ReturnedPrinciple, res.MortgagePrinciple)
		balance += res.MortgagePrinciple
	}
	log.Printf("Mortgage SERO is: %f\n", balance)
	return
}

func snapshot(ctx *cli.Context) error {
	stakeList := readStakeList(ctx)

	for _, v := range stakeList {
		if v.Expired > 0 {
			fmt.Printf("Expired: %s,%d,%d,%d,%s\n", v.Id, v.At, v.Expired, v.Remaining, v.Profit)
		} else {
			fmt.Printf("%s,%d,%s,%d,%d,%s\n", v.Id, v.At, v.Price, v.Total, v.Remaining, v.Profit)
		}
	}

	return nil
}

func profit(ctx *cli.Context) error {
	stakeList := readStakeList(ctx)

	for _, v := range stakeList {
		if v.Expired > 0 {
			fmt.Printf("Expired: %s,%d,%d,%d,%s\n", v.Id, v.At, v.Expired, v.Remaining, v.Profit)
		} else {
			fmt.Printf("%s,%d,%s,%d,%d,%s\n", v.Id, v.At, v.Price, v.Total, v.Remaining, v.Profit)
		}
	}

	return nil
}

func readStakeList(ctx *cli.Context) []stake.StakeDetail {
	filename := ctx.String(utils.StakeFlag.Name)
	com := ctx.String(utils.FilterCompanyFlag.Name)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	var stakeList []stake.StakeDetail
	for {
		var sd stake.StakeDetail
		if err := decoder.Decode(&sd); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			if com == "" || com == sd.Company {
				stakeList = append(stakeList, sd)
			}
		}
	}
	//for k, v := range stakeList {
	//	log.Println(k, v.Id, v.At)
	//}
	sort.Slice(stakeList, func(i, j int) bool {
		return stakeList[i].At < stakeList[j].At
	})
	return stakeList
}
