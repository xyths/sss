package main

import (
	"encoding/csv"
	"fmt"
	"github.com/xyths/sss/cmd/utils"
	"github.com/xyths/sss/sero/trader"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"path/filepath"
)

var app *cli.App

func init() {
	app = &cli.App{
		Name:    filepath.Base(os.Args[0]),
		Usage:   "bulk transfer money",
		Version: "0.1.2",
		Action:  bulkTransfer,
		Flags: []cli.Flag{
			ConfigFlag,
		},
	}
	utils.Info(app)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bulkTransfer(ctx *cli.Context) error {
	cfg := btrConfig{}
	if err := loadConfig(ctx.String(ConfigFlag.Name), &cfg); err != nil {
		log.Fatalf("error when open config file: %s", err)
	}
	log.Println(cfg)
	records, err := readData(cfg.Input)
	if err != nil {
		log.Fatal("error when read data from csv: %s", err)
	}
	t, err := trader.NewTrader(cfg.Sero)
	if err != nil {
		log.Fatalf("error when create trader: %s", err)
	}
	defer t.Close()

	t.BulkTransfer(ctx.Context, records)
	return nil
}

func readData(file string) (records [][]string, err error) {
	log.Printf("open input data from file: %s", file)
	f, err := os.Open(file)
	if err != nil {
		return
	}
	records, err = csv.NewReader(f).ReadAll()
	return
}
