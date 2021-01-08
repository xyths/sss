package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_GetPools(t *testing.T) {
	server := "https://sero-light-node.ririniannian.com/"
	pools, err := getPools(server)
	require.NoError(t, err)
	for i, p := range pools {
		t.Logf("[%d] %v", i, p)
	}
}

func Test_GetPrice(t *testing.T) {
	server := "https://sero-light-node.ririniannian.com/"
	price, err := getPrice(server)
	require.NoError(t, err)
	t.Logf("now price is %f", price)
}
