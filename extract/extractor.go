package extract

import (
	"context"
	"github.com/xyths/hs"
	"github.com/xyths/sero-go"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/big"
	"time"
)

type Config struct {
	Mongo    hs.MongoConf
	SeroNode string
}

type Extractor struct {
	config Config
	db     *mongo.Database
}

func New(ctx context.Context, configFile string) (*Extractor, error) {
	e := &Extractor{}
	if err := hs.ParseJsonConfig(configFile, &e.config); err != nil {
		return nil, err
	}

	db, err := hs.ConnectMongo(ctx, e.config.Mongo)
	if err != nil {
		return nil, err
	}
	e.db = db

	return e, nil
}

func (e *Extractor) Close() {
	if e.db != nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_ = e.db.Client().Disconnect(ctx)
	}
}

func (e *Extractor) Extract(ctx context.Context, start, end uint64) error {
	if start > end {
		return nil
	}
	api, err := sero.New(e.config.SeroNode)
	if err != nil {
		log.Println(err)
		return err
	}
	defer api.Close()

	for current := start; current <= end; current++ {
		//block, pkrs, err1 := e.extractOneBlock(ctx, api, current)
	}
	return nil
}

func (e *Extractor) extractOneBlock(ctx context.Context, api *sero.API, blockNumber uint64) (block sero.Block, pkrs []string, err error) {
	block, err = api.GetBlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return
	}
	for _, txHash := range block.Transactions {
		tp, err1 := api.GetTransactionReceipt(ctx, txHash)
		if err1 != nil {
			log.Printf("[ERROR] when extract transaction %s from block %d", txHash, blockNumber)
			continue
		}
		if tp.ShareId != "" {
			pkrs = append(pkrs, tp.From)
		}
	}
	return
}
