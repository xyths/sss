package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xyths/hs"
	"github.com/xyths/sss/cmd/utils"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server string
	Names  string
	Output string
	Log    hs.LogConf
}

type Pool struct {
	Id     string  `json:"id"`
	Shares int64   `json:"shares"`
	Value  float64 `json:"value"`
}

type ResponsePools struct {
	JsonRpc string    `json:"jsonrpc"`
	Id      int       `json:"id"`
	Result  []RawPool `json:"result"`
}
type RawPool struct {
	Closed bool   `json:"closed"`
	Id     string `json:"id"`
	Shares string `json:"shareNum"`
}
type ResponsePrice struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

func snap(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}

	l, err := hs.NewZapLogger(cfg.Log)
	if err != nil {
		return err
	}
	sugar := l.Sugar()
	nameCache := make(map[string]string)
	err = cacheNames(cfg.Names, nameCache)
	if err != nil {
		sugar.Errorf("cache name error: %s", err)
		return err
	}
	sugar.Infof("cached %d names", len(nameCache))
	pools, err := getPools(cfg.Server)
	if err != nil {
		sugar.Errorf("get pool error: %s", err)
		return err
	}
	price, err := getPrice(cfg.Server)
	if err != nil {
		sugar.Errorf("get price error: %s", err)
		return err
	}
	for i := 0; i < len(pools); i++ {
		pools[i].Value = float64(pools[i].Shares) * price
		if name, p := nameCache[pools[i].Id]; p {
			pools[i].Id = name
		}
	}
	err = log(cfg.Output, pools)
	if err != nil {
		sugar.Errorf("get pool error: %s", err)
		return err
	}
	return nil
}

func cacheNames(names string, cache map[string]string) error {
	f, err := os.Open(names)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line:= scanner.Text()
		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			continue
		}
		cache[tokens[0]] = tokens[1]
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func getPools(server string) ([]Pool, error) {
	resp, err := http.Post(server,
		"application/json",
		strings.NewReader(`{
    "id": 0,
    "jsonrpc": "2.0",
    "method": "stake_stakePools",
    "params": []
}`))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r ResponsePools
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	var pools []Pool
	for _, p := range r.Result {
		if p.Closed {
			continue
		}
		shares, err := strconv.ParseInt(p.Shares, 0, 64)
		if err != nil {
			continue
		}
		if shares == 0 {
			continue
		}
		id := p.Id
		if len(p.Id) >= 6 {
			id = p.Id[0:6]
		}
		pools = append(pools, Pool{
			Id:     id,
			Shares: shares,
		})
	}
	return pools, nil
}

func getPrice(server string) (float64, error) {
	resp, err := http.Post(server,
		"application/json",
		strings.NewReader(`{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "stake_sharePrice",
    "params": []
}`))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var r ResponsePrice
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, err
	}
	//price_, ok := big.NewInt(0).SetString(r.Result, 16)
	//if !ok {
	//	return 0, errors.New("can't parse price as integer")
	//}
	//price_.Div(price_, big.NewInt(1e16))
	//price := float64(price_.Int64()) / 100
	price_, _, err := big.NewFloat(0).Parse(r.Result, 0)
	if err != nil {
		return 0, err
	}
	price_.Quo(price_, big.NewFloat(1e18))
	price, _ := price_.Float64()
	return price, nil
}

func log(output string, pools []Pool) error {
	f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, p := range pools {
		data, err2 := json.Marshal(p)
		if err2 != nil {
			continue
		}
		_, err1 := fmt.Fprintf(f, "%s\n", string(data))
		if err1 != nil {
			continue
		}
	}
	return nil
}
