package sero

import (
	"context"
	"log"
)

func GetStakeInfo(ctx context.Context, id string) {
	log.Printf("GetStakeInfo share id = %s\n", id)
}

func GetTransactionReceipt(ctx context.Context, id string) {
	log.Printf("GetTransactionReceipt id = %s\n", id)
}
