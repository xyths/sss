package main

import (
	"bufio"
	"encoding/json"
	"github.com/xyths/sss/sero/trader"
	"os"
)

type btrConfig struct {
	Input string
	Sero  trader.SeroConfig
}

func loadConfig(file string, cfg *btrConfig) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}
