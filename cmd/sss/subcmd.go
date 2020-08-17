package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xyths/sss/cmd/utils"
	"github.com/xyths/sss/extract"
	. "github.com/xyths/sss/mail"
	"github.com/xyths/sss/mail/client"
	"github.com/xyths/sss/share"
	"github.com/xyths/sss/stake"
	"gopkg.in/urfave/cli.v2"
	"html/template"
	"io"
	"log"
	"os"
	"sort"
)

var (
	appendCommand = &cli.Command{
		Action:  appendShare,
		Name:    "append",
		Aliases: []string{"a"},
		Usage:   "Append stake share to db",
		Flags: []cli.Flag{
			utils.AppendConfigFlag,
			utils.ShareListFlag,
		},
	}
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
	mailCommand = &cli.Command{
		Action:  mail,
		Name:    "mail",
		Aliases: []string{"m"},
		Usage:   "mail to investors",
		Flags: []cli.Flag{
			utils.MailConfigFlag,
			utils.DateFlag,
		},
	}
	testCommand = &cli.Command{
		Action:  test,
		Name:    "test",
		Aliases: []string{"t"},
		Usage:   "test some demo program",
		Flags: []cli.Flag{
			utils.MailConfigFlag,
		},
	}
	extractCommand = &cli.Command{
		Action: extractAction,
		Name:   "extract",
		//Aliases: []string{"t"},
		Usage: "extract address from blockchain",
		Flags: []cli.Flag{
			utils.ConfigFlag,
			utils.StartBlockFlag,
			utils.EndBlockFlag,
		},
	}
)

func appendShare(ctx *cli.Context) (err error) {
	filename := ctx.String(utils.MailConfigFlag.Name);
	if filename == "" {
		return errors.New("no config")
	}
	config := utils.LoadAppendConfig(filename)
	csvfile := ctx.String(utils.ShareListFlag.Name)
	return share.AppendShare(config, csvfile)
}

func sum(ctx *cli.Context) (err error) {
	stakeList := readStakeList(ctx)

	var balance float64
	var results []stake.Result
	for _, sd := range stakeList {
		res := stake.Sum(sd)
		results = append(results, res)
		log.Printf("%s,%d,%d,%d,%f,%f,%f,%f,%f,%f\n", res.Id,
			res.TotalShare, res.ReturnedShare, res.MortgageShare,
			res.TotalPrinciple, res.ReturnedPrinciple, res.MortgagePrinciple,
			res.TotalInterest, res.ReturnedInterest, res.MortgageInterest)
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

func mail(ctx *cli.Context) error {
	filename := ctx.String(utils.MailConfigFlag.Name);
	if filename == "" {
		return errors.New("no config")
	}
	config := utils.Load(filename)
	date := ctx.String(utils.DateFlag.Name)

	return Mail(config, date)
}

func test(ctx *cli.Context) error {
	user := "pangu_sero_pos@163.com"
	password := "pgsp20190916"
	host := "smtp.163.com:25"
	to := "xing_yongtao@163.com"

	subject := "Test send email by Golang"

	tpl1 := template.New("template.html")
	tpl, err := tpl1.ParseFiles("mail/template/template.html");
	if err != nil {
		log.Println(err)
		return err
	}

	filename := "data/share_20190912.json"
	com := "xh"
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
			if com == sd.Company {
				stakeList = append(stakeList, sd)
			}
		}
	}

	sort.Slice(stakeList, func(i, j int) bool {
		return stakeList[i].At < stakeList[j].At
	})

	var results []stake.Result

	for _, sd := range stakeList {
		res := stake.Sum(sd)
		results = append(results, res)
		log.Printf("%s,%d,%d,%d,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f\n", res.ShortId,
			res.TotalShare, res.ReturnedShare, res.MortgageShare,
			res.TotalPrinciple, res.ReturnedPrinciple, res.MortgagePrinciple,
			res.TotalInterest, res.ReturnedInterest, res.MortgageInterest)
	}

	report := stake.FormatReport(results)

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, report); err != nil {
		log.Println(err)
	}

	fmt.Println("send email")
	err = client.SendMail(user, password, host, to, subject, buf.String(), "html")
	if err != nil {
		fmt.Println("send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("send mail success!")
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

	sort.Slice(stakeList, func(i, j int) bool {
		return stakeList[i].At < stakeList[j].At
	})
	return stakeList
}

func extractAction(ctx *cli.Context) error {
	config := ctx.String(utils.ConfigFlag.Name)
	start := ctx.Uint64(utils.StartBlockFlag.Name)
	end := ctx.Uint64(utils.EndBlockFlag.Name)
	extractor, err := extract.New(ctx.Context, config)
	if err != nil {
		log.Fatal(err)
	}
	defer extractor.Close()

	log.Printf("start extract address from block %d to %d", start, end)
	err = extractor.Extract(ctx, start, end)
	log.Println("finish extract")
	return err
}
