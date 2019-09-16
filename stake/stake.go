package stake

import (
	"context"
	"github.com/sero-cash/go-sero/seroclient"
	"github.com/xyths/sss/sero"
	"log"
	"math/big"
)

// 每 42336 块一个支付点，支付本金和利息
// 过期是181,440块，但是需要等211,680块才会返回。
const (
	PayTime       = 42336
	ExpireTime    = 181440
	ExpirePayTime = 211680
)

type Result struct {
	Id string

	BuyBlock       int
	LastPayBlock   int
	ExpireBlock    int
	ExpirePayBlock int

	// Share Number
	TotalShare int

	MortgageShare       int // Mortagage
	RemainShare         int //   Remaining
	IncomeShare         int //   Income not return
	ExpireNoReturnShare int //   Expire not return

	ReturnedShare         int // Returned
	CheckedShare          int //   Checked
	ExpiredAndReturnShare int //   Expired and return

	// Principle
	TotalPrinciple    float64 // 所有本金
	MortgagePrinciple float64 // 抵押中本金
	ReturnedPrinciple float64 // 已返还本金

	// Interest
	TotalInterest    float64 // 所有本金
	MortgageInterest float64 // 抵押中本金
	ReturnedInterest float64 // 已返还本金
}

type StakeDetail struct {
	Id        string `json:"id"`
	Company   string `json:"company"`
	Tx        string `json:"tx"`
	At        int    `json:"at"`   // tx's blockNumber
	Addr      string `json:"addr"` // Owner
	Pool      string `json:"pool"`
	VoteAddr  string `json:"voteAddr"`
	Fee       int    `json:"fee"`
	Timestamp int    `json:"timestamp"` // buy time

	Price string `json:"price"`
	Total int    `json:"total"`

	// remaining 和 expired 二选一只有一个
	Remaining int    `json:"remaining"`
	Expired   int    `json:"expired"`
	Missed    int    `json:"missed"`
	Profit    string `json:"profit"` //all profit

	ReturnNum    int    `json:"returnNum"`    // blockNumber
	LastPayTime  int    `json:"lastPayTime"`  // blockNumber
	ReturnProfit string `json:"returnProfit"` // returned profit

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

func Sum(sd StakeDetail) (r Result) {
	r.Id = sd.Id
	r.TotalShare = sd.Total
	r.CheckedShare = sd.ReturnNum

	start := sd.At
	lastPayTime := sd.LastPayTime

	if lastPayTime-start >= ExpirePayTime {
		//log.Println("\tnow expired should be return.")
		r.ExpiredAndReturnShare += sd.Expired
	} else if lastPayTime-start >= ExpirePayTime {
		//log.Println("\tnow expire happen, but sero not return")
		r.ExpireNoReturnShare += sd.Expired
	} else {
		r.RemainShare = sd.Remaining
		r.IncomeShare = sd.Total - sd.Remaining - sd.ReturnNum
	}

	r.MortgageShare = r.RemainShare + r.IncomeShare + r.ExpireNoReturnShare
	r.ReturnedShare = r.CheckedShare + r.ExpiredAndReturnShare

	//log.Printf(`	id: %s
	//buyBlock: %d
	//lastPayBlock: %d
	//TotalShare: %d
	//ReturnedShare:	%d
	//	Checked:	%d
	//	Expire:		%d
	//MortgageShare:	%d
	//	Remaining:	%d
	//	Income:		%d
	//	Expire:		%d`, r.Id, r.BuyBlock, r.LastPayBlock,
	//	r.TotalShare,
	//	r.ReturnedShare, r.CheckedShare, r.ExpiredAndReturnShare,
	//	r.MortgageShare, r.RemainShare, r.IncomeShare, r.ExpireNoReturnShare)

	price := big.NewFloat(0)
	price.SetString(sd.Price)
	price.Quo(price, SERO)

	totalPrinciple := big.NewFloat(0).SetInt64(int64(r.TotalShare))
	totalPrinciple.Mul(totalPrinciple, price)
	r.TotalPrinciple, _ = totalPrinciple.Float64()
	mortgagePrinciple := big.NewFloat(0).SetInt64(int64(r.MortgageShare))
	mortgagePrinciple.Mul(mortgagePrinciple, price)
	r.MortgagePrinciple, _ = mortgagePrinciple.Float64()
	returnedPrinciple := big.NewFloat(0).SetInt64(int64(r.ReturnedShare))
	returnedPrinciple.Mul(returnedPrinciple, price)
	r.ReturnedPrinciple, _ = returnedPrinciple.Float64()

	return
}
