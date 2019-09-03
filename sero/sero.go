package sero

import (
	"context"
	"github.com/sero-cash/go-sero/common"
	"github.com/sero-cash/go-sero/seroclient"
	"log"
)

func GetStakeInfo(client *seroclient.Client, ctx context.Context, id string) {
	log.Printf("GetStakeInfo share id = %s\n", id)

	if networkID, err := client.NetworkID(ctx); err == nil {
		log.Printf("Network ID is: %s\n", networkID)
	} else {
		log.Println(err)
	}
}

func GetTransactionReceipt(client *seroclient.Client, ctx context.Context, id string) {
	log.Printf("GetTransactionReceipt id = %s\n", id)
	tx := "0xa27ba4fed3c6767a6943ec0dbdd57e4330cb6d2b771119111e7de0acc819cdb7"
	if receipt, err := client.TransactionReceipt(ctx, common.HexToHash(tx)); err == nil {
		log.Println(receipt)
	} else {
		log.Println(err)
	}

}
