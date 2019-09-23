package stake

import (
	"context"
	"github.com/xyths/sss/sero"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"math/big"
	"strconv"
)

// 每 42336 块一个支付点，支付本金和利息
// 过期是181,440块，但是需要等211,680块才会返回。
const (
	PayTime       = 42336
	ExpireTime    = 181440
	ExpirePayTime = 211680
)

type Result struct {
	Id      string
	ShortId string

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

type ShareReport struct {
	Id string

	// Share Number
	TotalShare    string // 全部票数
	ReturnedShare string // 返还票数
	MortgageShare string // 剩余票数

	// Principle
	TotalPrinciple    string // 所有本金
	ReturnedPrinciple string // 返还本金
	MortgagePrinciple string // 剩余本金

	// Interest
	TotalInterest    string // 所有利润
	ReturnedInterest string // 返还利润
	MortgageInterest string // 剩余利润
}

type Report struct {
	Shares  []ShareReport
	Summary ShareReport
}

func Stat(id string) (result Result, err error) {
	log.Printf("start process share id = %s\n", id)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("we have a connection now.")
	}

	ctx := context.Background()
	sero.GetStakeInfo(ctx, id)
	sero.GetTransactionReceipt(ctx, id)

	log.Printf("end process share id = %s\n", id)
	return
}

var SERO = big.NewFloat(1e18)

func Sum(sd StakeDetail) (r Result) {
	r.Id = sd.Id
	r.ShortId = r.short(r.Id)
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

	price := big.NewFloat(0)
	price.SetString(sd.Price)
	price.Quo(price, SERO)

	totalPrinciple := big.NewFloat(0)
	totalPrinciple.SetInt64(int64(r.TotalShare))
	totalPrinciple.Mul(totalPrinciple, price)
	r.TotalPrinciple, _ = totalPrinciple.Float64()
	mortgagePrinciple := big.NewFloat(0)
	mortgagePrinciple.SetInt64(int64(r.MortgageShare))
	mortgagePrinciple.Mul(mortgagePrinciple, price)
	r.MortgagePrinciple, _ = mortgagePrinciple.Float64()
	returnedPrinciple := big.NewFloat(0)
	returnedPrinciple.SetInt64(int64(r.ReturnedShare))
	returnedPrinciple.Mul(returnedPrinciple, price)
	r.ReturnedPrinciple, _ = returnedPrinciple.Float64()

	// Interest
	totalInterest := big.NewFloat(0)
	totalInterest.SetString(sd.Profit)
	totalInterest.Quo(totalInterest, SERO)
	r.TotalInterest, _ = totalInterest.Float64()
	returnedInterest := big.NewFloat(0)
	returnedInterest.SetString(sd.ReturnProfit)
	returnedInterest.Quo(returnedInterest, SERO)
	r.ReturnedInterest, _ = returnedInterest.Float64()
	r.MortgageInterest = r.TotalInterest - r.ReturnedInterest

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
	//	Expire:		%d
	//TotalPrinciple:		%.2f
	//ReturnedPrinciple:	%.2f
	//MortgagePrinciple:	%.2f
	//TotalInterest:		%.2f
	//ReturnedInterest:	%.2f
	//MortgageInterest:	%.2f
	//`,
	//	r.Id, r.BuyBlock, r.LastPayBlock,
	//	r.TotalShare,
	//	r.ReturnedShare, r.CheckedShare, r.ExpiredAndReturnShare,
	//	r.MortgageShare, r.RemainShare, r.IncomeShare, r.ExpireNoReturnShare,
	//	r.TotalPrinciple, r.ReturnedPrinciple, r.MortgagePrinciple,
	//	r.TotalInterest, r.ReturnedInterest, r.MortgageInterest)

	return
}

func FormatReport(results []Result) (report Report) {
	var s Result
	for _, r := range results {
		report.Shares = append(report.Shares, FormatShare(r))
		s.TotalShare += r.TotalShare
		s.ReturnedShare += r.ReturnedShare
		s.MortgageShare += r.MortgageShare
		s.TotalPrinciple += r.TotalPrinciple
		s.ReturnedPrinciple += r.ReturnedPrinciple
		s.MortgagePrinciple += r.MortgagePrinciple
		s.TotalInterest += r.TotalInterest
		s.ReturnedInterest += r.ReturnedInterest
		s.MortgageInterest += r.MortgageInterest
	}
	report.Summary = FormatShare(s)
	return report
}

func FormatShare(r Result) (sr ShareReport) {
	sr.Id = r.ShortId
	sr.TotalShare = formatInt(r.TotalShare)
	sr.ReturnedShare = formatInt(r.ReturnedShare)
	sr.MortgageShare = formatInt(r.MortgageShare)
	sr.TotalPrinciple = formatFloat64(r.TotalPrinciple)
	sr.ReturnedPrinciple = formatFloat64(r.ReturnedPrinciple)
	sr.MortgagePrinciple = formatFloat64(r.MortgagePrinciple)
	sr.TotalInterest = formatFloat64(r.TotalInterest)
	sr.ReturnedInterest = formatFloat64(r.ReturnedInterest)
	sr.MortgageInterest = formatFloat64(r.MortgageInterest)
	return sr
}

func formatInt(i int) string {
	if i > 0 {
		return strconv.Itoa(i)
	} else {
		return "-"
	}
}

func formatFloat64(f float64) string {
	if f > 0 {
		p := message.NewPrinter(language.English)
		return p.Sprintf("%.2f", f)
	} else {
		return "-"
	}
}

func (r Result) short(Id string) (shortId string) {
	shortId = Id[0:6] + "..." + Id[len(Id)-4:len(Id)]
	return shortId
}
