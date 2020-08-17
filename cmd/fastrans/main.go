package main

import (
	"context"
	"fmt"
	"github.com/xyths/hs"
	"github.com/xyths/hs/logger"
	"github.com/xyths/sero-go"
	"github.com/xyths/sss/cmd/utils"
	"gopkg.in/urfave/cli.v2"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

var (
	app *cli.App

	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "config.json",
		Usage:   "Read config from `FILE`",
	}
)

func init() {
	app = &cli.App{
		Name:    filepath.Base(os.Args[0]),
		Usage:   "fast transfer money, detect and transfer in specific period",
		Version: "0.1.3",
		Action:  fastTransfer,
		Flags: []cli.Flag{
			ConfigFlag,
		},
	}
	utils.Info(app)
}

type Config struct {
	Interval string // "10s"
	Source   string
	Cache    string
	ReFund   string
	Wait     int

	Gas int
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func fastTransfer(ctx *cli.Context) error {
	cfg := Config{}

	if err := hs.ParseJsonConfig(ctx.String(ConfigFlag.Name), &cfg); err != nil {
		logger.Sugar.Fatalf("error when open config file: %s", err)
	}

	interval, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		logger.Sugar.Fatalf("interval format error: %s", err)
	}

	logger.Sugar.Infof("source: %s...., cache: %s....", cfg.Source, cfg.Cache)

	for {
		select {
		case <-ctx.Context.Done():
			logger.Sugar.Info("fastrans cancelled")
			return nil
		case <-time.After(interval):
			doWork(ctx.Context, cfg.Source, cfg.Cache, cfg.ReFund, cfg.Wait)
		}
	}

}

func doWork(ctx context.Context, source, dest, refund string, wait int) {
	//log.Printf("doWork, try to transfer: %s -> %s", source, dest)
	api, err := sero.New("http://127.0.0.1:8545")
	defer api.Close()

	// 1. check balance
	b, err := api.Balance(ctx, source)
	if err != nil {
		logger.Sugar.Errorf("error when call exchange_getBalances: %s", err)
		return
	}
	if b.SERO == "" {
		logger.Sugar.Infof("no balance")
		return
	}
	logger.Sugar.Infof("%s total balance is %v", source, b)
	balance, ok := big.NewInt(0).SetString(b.SERO, 10)
	if !ok {
		logger.Sugar.Errorf("balance format error: %v", b)
		return
	}
	gas := big.NewInt(25000000000000)
	if balance.Cmp(gas) <= 0 {
		logger.Sugar.Infof("balance too low: %s", balance)
		return
	}
	balance.Sub(balance, gas)
	logger.Sugar.Infof("will send %s wei SERO", balance)

	// 2. try to transfer
	if _, err = api.SendAndWait(ctx, source, refund, dest, "SERO", balance, wait); err != nil {
		logger.Sugar.Errorf("transfer error: %s", err)
	}
}
