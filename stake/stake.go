package stake

import (
	"context"
	"github.com/sero-cash/go-sero/seroclient"
	"github.com/xyths/sss/sero"
	"log"
	"math/big"
)

type Result struct {
	id string
}

type StakeDetail struct {
	Id       string `json:"id"`
	Company  string `json:"company"`
	Tx       string `json:"tx"`
	Addr     string `json:"addr"`
	Pool     string `json:"pool"`
	VoteAddr string `json:"voteAddr"`
	Fee      int    `json:"fee"`

	Price  string `json:"price"`
	Total  int    `json:"total"`
	Missed int    `json:"missed"`
	Profit string `json:"profit"`

	// remaining 和 expired 二选一只有一个
	Remaining int `json:remaining`
	Expired   int `json:"expired"`

	Status int `json:"status"`
}

func Stat(id string) (result Result, err error) {
	log.Printf("start process share id = %s\n", id)
	// 华南(成都): http://148.70.169.73:8545
	// 华南(广州): http://129.204.197.105:8545
	// Japan:     http://52.199.145.159:8545
	client, err := seroclient.Dial("http://148.70.169.73:8545")
	defer client.Close()

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("we have a connection now.")
	}

	ctx := context.Background()
	sero.GetStakeInfo(client, ctx, id)
	sero.GetTransactionReceipt(client, ctx, id)

	log.Printf("end process share id = %s\n", id)
	return
}

var SERO = big.NewFloat(1e18)

func Sum(sd StakeDetail) (remain *big.Float, err error) {
	remain = big.NewFloat(0)

	if sd.Remaining == 0 {
		return
	}
	remaining := big.NewFloat(0).SetInt64(int64(sd.Remaining))

	price := big.NewFloat(0)
	price.SetString(sd.Price)

	price.Quo(price, SERO)

	remain.Mul(price, remaining)

	return
}
