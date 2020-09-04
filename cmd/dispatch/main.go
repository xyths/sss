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
		Usage:   "dispatch the coin to all destinations",
		Version: "0.1.3",
		Action:  dispatch,
		Flags: []cli.Flag{
			ConfigFlag,
		},
	}
	utils.Info(app)
}

type Config struct {
	Interval     string
	Source       string
	Destinations []string
	ReFund       string
	Wait         int
	Reserved     string
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func dispatch(ctx *cli.Context) error {
	cfg := Config{}

	if err := hs.ParseJsonConfig(ctx.String(ConfigFlag.Name), &cfg); err != nil {
		logger.Sugar.Fatalf("error when open config file: %s", err)
	}

	interval, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		logger.Sugar.Fatalf("interval format error: %s", err)
	}

	logger.Sugar.Infof("source: %s...., destinations: %v....", cfg.Source, cfg.Destinations)

	doWork(ctx.Context, cfg.Source, cfg.Destinations, cfg.ReFund, cfg.Wait, cfg.Reserved)
	for {
		select {
		case <-ctx.Context.Done():
			logger.Sugar.Info("dispatch cancelled")
			return nil
		case <-time.After(interval):
			doWork(ctx.Context, cfg.Source, cfg.Destinations, cfg.ReFund, cfg.Wait, cfg.Reserved)
		}
	}

}

func doWork(ctx context.Context, source string, destinations []string, refund string, wait int, reserved string) {
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
	keep, _ := big.NewInt(0).SetString(reserved, 0)
	if balance.Cmp(keep.Add(keep, gas)) <= 0 {
		logger.Sugar.Infof("balance too low: %s", balance)
		return
	}
	balance.Sub(balance, keep)
	amount := big.NewInt(0).Div(balance, big.NewInt(2))
	logger.Sugar.Infof("will send %s wei SERO, %s wei each", balance, amount)

	// 2. try to transfer

	if _, err = api.MultiSendAndWait(ctx, source, refund, destinations, []string{"SERO", "SERO"}, []*big.Int{amount, amount}, wait); err != nil {
		logger.Sugar.Errorf("transfer error: %s", err)
	}
}
