package trader

import (
	"context"
	"errors"
	"github.com/sero-cash/go-sero/rpc"
	"github.com/sero-cash/go-sero/zero/txtool"
	"github.com/xyths/sero-go"
	"github.com/xyths/sss/sero/types"
	"log"
	"math/big"
	"time"
)

const (
	currency      = "VIRTUE"
	confirmNumber = 32
	blockTime     = 14 * time.Second
)

type SeroConfig struct {
	Rpc string

	From     string
	Refund   string
	Gas      *big.Int
	GasPrice *big.Int
}

type Trader struct {
	client   *rpc.Client
	Rpc      string
	From     string
	Refund   string
	Gas      *big.Int
	GasPrice *big.Int
}

func NewTrader(cfg SeroConfig) (t *Trader, err error) {
	t = &Trader{
		Rpc:      cfg.Rpc,
		From:     cfg.From,
		Refund:   cfg.Refund,
		Gas:      cfg.Gas,
		GasPrice: cfg.GasPrice,
	}
	t.client, err = rpc.Dial(t.Rpc)

	return
}

func (t *Trader) BulkTransfer(ctx context.Context, records [][]string) {
	unit := big.NewInt(1e18)
	api, err := sero.New(t.Rpc)
	defer api.Close()
	if err != nil {
		log.Fatalf("error when new api: %s", err)
	}
	for i, r := range records {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			return
		default:
			log.Printf("[%d] %s %s %s", i, r[0], r[2], r[3])
			to := r[2]

			amount, ok := big.NewInt(0).SetString(r[3], 10)
			if !ok {
				log.Printf("error when parse amount: %s", r[3])
				continue
			}
			amount.Mul(amount, unit)
			//hash, err := t.transfer(ctx, to, currency, amount)
			//if err != nil {
			//	log.Printf("error when transfer %s", to)
			//	continue
			//}
			//now := time.Now().Format("2006-01-02 15:04:05")
			//log.Printf("[%d] %s sent to %s(%s) %d %s, hash is %s", i, now, r[0], to, amount, currency, hash)
			if trans, err := api.SendAndWait(ctx, t.From, t.Refund, to, currency, amount, 32); err == nil {
				log.Printf("[INFO] to: %s, amount: %s, tx: %s, time: %s", to, amount, trans.TransactionHash, time.Now())
			} else {
				log.Printf("error when SendAndWait: %s", err)
			}
		}
	}
	return
}

func (t *Trader) transfer(ctx context.Context, to string, currency string, amount *big.Int) (hash string, err error) {
	log.Printf("to: %s, currency: %s, amount: %s", to, currency, amount.String())

	// 0. clear used flag
	if err = t.ClearUsedFlag(ctx); err != nil {
		return
	}

	// 1. check balance
	totalBalance, err := t.GetBalance(ctx)
	if err != nil {
		return
	}
	if totalBalance.Cmp(amount) <= 0 {
		return hash, errors.New("no enough balance")
	}

	// 2. check available balance
	maxAvailableBalance, err := t.GetMaxAvailable(ctx)
	if err != nil {
		return
	}
	if maxAvailableBalance.Cmp(amount) <= 0 {
		return hash, errors.New("no enough balance")
	}

	// 3. gen tx

	gtx, err := t.GenTx(ctx, to, currency, amount)
	var txHash string
	if hByte, err := gtx.Hash.MarshalText(); err == nil {
		txHash = string(hByte)
	} else {
		log.Printf("MarshalText error: %s", err)
		return hash, err // inter err
	}
	hash = txHash
	log.Printf("exchange_genTxWithSign get tx: %s", txHash)

	// 4. commit tx
	var result interface{}
	err = t.client.Call(&result, "exchange_commitTx", &gtx);
	if err != nil {
		log.Printf("error when call exchange_commitTx: %s", err)
		return
	}
	log.Println("commit the transaction")

	// 5. check tx blockNumber
	blockNumber := uint64(0)
	for blockNumber == 0 {
		select {
		case <-ctx.Done():
			//log.Println(ctx.Err())
			return hash, nil
		case <-time.After(blockTime):
			log.Printf("check for tx %s", txHash)
			if blockNumber, err = t.CheckTx(ctx, txHash); err == nil && blockNumber > 0 {
				log.Printf("tx %s is at block %d", txHash, blockNumber)
			}
		}
	}

	// 6. wait for confirm & balance analysis
	currentBlock := uint64(0)
	for blockNumber+confirmNumber > currentBlock {
		select {
		case <-ctx.Done():
			//log.Println(ctx.Err())
			return hash, nil
		case <-time.After(blockTime):
			log.Printf("check for tx %s", txHash)
			if currentBlock, err = t.GetBlockNumber(ctx); err == nil && currentBlock > 0 {
				log.Printf("current block number is %d, confirm need at least %d", currentBlock, blockNumber+confirmNumber)
			}
		}
	}

	return
}

func (t *Trader) Close() {
	t.client.Close()
}

func (t Trader) ClearUsedFlag(ctx context.Context) error {
	var result int
	err := t.client.Call(&result, "exchange_clearUsedFlag", &t.From)
	if err != nil {
		log.Printf("error when call exchange_clearUsedFlag: %s", err)
		return err
	}
	log.Printf("clear used flag result: %d", result)
	return nil
}

func (t Trader) GetBalance(ctx context.Context) (balance *big.Int, err error) {
	var b types.Balance
	err = t.client.Call(&b, "exchange_getBalances", &t.From)
	if err != nil {
		log.Printf("error when call exchange_getBalances: %s", err)
		return
	}
	balance, ok := big.NewInt(0).SetString(b.Tkn[currency], 10)
	if !ok {
		return balance, errors.New("error when parse balance")
	}
	log.Printf("%s total balance is %s", t.From, balance.String())
	return
}

func (t Trader) GetMaxAvailable(ctx context.Context) (balance *big.Int, err error) {
	var balanceStr string
	err = t.client.Call(&balanceStr, "exchange_getMaxAvailable", &t.From, currency)
	if err != nil {
		log.Printf("error when call exchange_getMaxAvailable: %s", err)
	}
	log.Printf("account %s has available %s %s", t.From, currency, balanceStr)
	balance, ok := big.NewInt(0).SetString(balanceStr, 10)
	if !ok {
		return balance, errors.New("error when parse max available balance")
	}
	log.Printf("%s max available balance is %s", t.From, balance.String())
	return
}

func (t Trader) GenTx(ctx context.Context, to string, currency string, amount *big.Int) (gtx txtool.GTx, err error) {
	receptions := []types.Reception{
		{
			Currency: currency,
			Addr:     to,
			Value:    amount.String(),
		},
	}
	preTxParam := types.PreTxParam{
		From:       t.From,
		RefundTo:   t.Refund,
		Gas:        25000,
		GasPrice:   1000000000,
		Receptions: receptions,
		Roots:      []string{},
	}

	err = t.client.Call(&gtx, "exchange_genTxWithSign", &preTxParam)
	if err != nil {
		log.Printf("error when call exchange_genTxWithSign: %s", err)
	}

	return
}

func (t Trader) CheckTx(ctx context.Context, txHash string) (blockNumber uint64, err error) {
	var transactionReceipt types.TransactionReceipt
	if err := t.client.Call(&transactionReceipt, "sero_getTransactionReceipt", &txHash); err != nil {
		log.Printf("error when call sero_getTransactionReceipt: %s", err)
	}
	log.Printf("sero_getTransactionReceipt of tx %s: %v", txHash, transactionReceipt)
	log.Printf("blockNumber is %s", transactionReceipt.BlockNumber)
	if transactionReceipt.BlockNumber == "" {
		return
	}
	b, ok := big.NewInt(0).SetString(transactionReceipt.BlockNumber, 0)
	if ok {
		blockNumber = b.Uint64()
	} else {
		return blockNumber, errors.New("bad blockNumber")
	}
	return
}

//curl -X POST 'http://47.92.64.129:8545' -H 'Content-Type:application/json' -d '{
//    "id": 0,
//    "jsonrpc": "2.0",
//    "method": "sero_blockNumber",
//    "params": []
//}'
//{"jsonrpc":"2.0","id":0,"result":"0x302a"}
func (t Trader) GetBlockNumber(ctx context.Context) (blockNumber uint64, err error) {
	var blockNumStr string
	err = t.client.Call(&blockNumStr, "sero_blockNumber")
	if err != nil {
		log.Printf("error when call sero_blockNumber: %s", err)
		return
	}

	b, ok := big.NewInt(0).SetString(blockNumStr, 0)
	if ok {
		blockNumber = b.Uint64()
	} else {
		return blockNumber, errors.New("bad blockNumber")
	}
	return
}
