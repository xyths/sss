package types

import "github.com/sero-cash/go-sero/core/types"

type Balance struct {
	SERO string              `json:"SERO"`
	Tkn  map[string]string   `json:"tkn"`
	Tkt  map[string][]string `json:"tkt"`
}

type Reception struct {
	Addr     string `json:"Addr"`
	Currency string `json:"Currency"`
	Value    string  `json:"Value"`
}

type PreTxParam struct {
	From       string      `json:"From"`
	RefundTo   string      `json:"RefundTo"`
	Gas        int64       `json:"Gas"`
	GasPrice   int64       `json:"GasPrice"`
	Receptions []Reception `json:"Receptions"`
	Roots      []string    `json:"Roots"`
}

type TransactionReceipt struct {
	BlockHash         string       `json:"blockHash"`
	BlockNumber       string       `json:"blockNumber"`
	TransactionHash   string       `json:"transactionHash"`
	TransactionIndex  string       `json:"transactionIndex"`
	From              string       `json:"from"`
	To                string       `json:"to"`
	GasUsed           string       `json:"gasUsed"`
	CumulativeGasUsed string       `json:"cumulativeGasUsed"`
	ContractAddress   string       `json:"contractAddress"`
	Logs              []*types.Log `json:"logs"`
	LogsBloom         string       `json:"logsBloom"`
	ShareId           string       `json:"shareId"`
	PoolId            string       `json:"poolId"`
}
