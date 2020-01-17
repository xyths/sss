package trader

import (
	"context"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"testing"
)

func TestTrader_BulkTransfer(t *testing.T) {

}

// TDRPC=http://47.92.64.129:8545 TDFROM=3ZzFEq5Wy7LVV7QtVrwbTW4axrUKDrdp98ybX2VTM7QjQjqLXkBbGRdjjqd2Vq4LiRRdAGrDCTE21jLdWYDERTh5 TDREFUND=5kbeHtrdaYKZqeMyJRsjy8P1FwsN6Y26EWf18JCPUwGEpBp4Qwu9wKEAqki5AwBRvLvxcummk2NNHoVmtiorLXrmPccimvpPZuv4mmfRkAayEpPhkS9hNjqPAuwsw47vufk go test -v ./sero/trader -run=TestTrader_GetBalance
func TestTrader_GetBalance(t *testing.T) {
	trader, err := NewTrader(SeroConfig{
		Rpc:      os.Getenv("TDRPC"),
		From:     os.Getenv("TDFROM"),
		Refund:   os.Getenv("TDREFUND"),
		Gas:      big.NewInt(25000),
		GasPrice: big.NewInt(1000000000),
	})
	defer trader.Close()
	require.NoError(t, err)

	b, err := trader.GetBalance(context.Background())
	require.NoError(t, err)
	t.Logf("%s balance: %s", trader.From, b.String())
}

// TDRPC=http://47.92.64.129:8545 TDFROM=3ZzFEq5Wy7LVV7QtVrwbTW4axrUKDrdp98ybX2VTM7QjQjqLXkBbGRdjjqd2Vq4LiRRdAGrDCTE21jLdWYDERTh5 TDREFUND=5kbeHtrdaYKZqeMyJRsjy8P1FwsN6Y26EWf18JCPUwGEpBp4Qwu9wKEAqki5AwBRvLvxcummk2NNHoVmtiorLXrmPccimvpPZuv4mmfRkAayEpPhkS9hNjqPAuwsw47vufk go test -v ./sero/trader -run=TestTrader_GetMaxAvailable
func TestTrader_GetMaxAvailable(t *testing.T) {
	trader, err := NewTrader(SeroConfig{
		Rpc:      os.Getenv("TDRPC"),
		From:     os.Getenv("TDFROM"),
		Refund:   os.Getenv("TDREFUND"),
		Gas:      big.NewInt(25000),
		GasPrice: big.NewInt(1000000000),
	})
	defer trader.Close()
	require.NoError(t, err)

	b, err := trader.GetMaxAvailable(context.Background())
	require.NoError(t, err)
	t.Logf("%s balance: %s", trader.From, b.String())
}

func TestTrader_GenTx(t *testing.T) {
	trader, err := NewTrader(SeroConfig{
		Rpc:      os.Getenv("TDRPC"),
		From:     os.Getenv("TDFROM"),
		Refund:   os.Getenv("TDREFUND"),
		Gas:      big.NewInt(25000),
		GasPrice: big.NewInt(1000000000),
	})
	defer trader.Close()
	require.NoError(t, err)

	to := "ZBR8cDAWFE8ueQ5B4suci8t3m7mtYRLDAZ6H4EVJRDLzpbvAheVPdWMoqnziTHYE7jEFb3VYqnrQSxNG7oLpHDorhK8NMyGYUAtav4pQLTkiYQqviiYFHYexdJZmXgpH4kH"
	amount, _ := big.NewInt(0).SetString("100000000000000000000", 10) // 100

	gtx, err := trader.GenTx(context.Background(), to, currency, amount)
	require.NoError(t, err)
	hByte, err := gtx.Hash.MarshalText()
	require.NoError(t, err)
	txHash := string(hByte)
	t.Logf("%s gen tx hash: %s", trader.From, txHash)
}
