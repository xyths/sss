package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xyths/sss/cmd/utils"
	"os"
	"path/filepath"
)

var app *cli.App

func init() {
	app = &cli.App{
		Name:    filepath.Base(os.Args[0]),
		Usage:   "The sss(SERO stake statistics) command line interface",
		Version: "0.1.0",
		Action:  sss,
	}
	utils.Info(app)
	app.Commands = []*cli.Command{
		appendCommand,
		sumCommand,
		snapshotCommand,
		profitCommand,
		mailCommand,
		testCommand,
		extractCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func sss(ctx *cli.Context) error {
	//filename := ctx.String(utils.FileFlag.Name)
	//file, err := os.Open(filename)
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer file.Close()
	//
	//reader := bufio.NewReader(file)
	//
	//var results []stake.Result
	//for {
	//	line, _, err := reader.ReadLine()
	//
	//	if err == io.EOF {
	//		break
	//	}
	//	if res, err := stake.Stat(string(line)); err != nil {
	//		results = append(results, res)
	//	}
	//	//log.Printf("%s\n", line)
	//}
	return nil
}
