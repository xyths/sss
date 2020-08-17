package extract

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/xyths/sero-go"
	"os"
	"strconv"
	"testing"
)

func TestExtractor_Extract(t *testing.T) {
	config := os.Getenv("CONF")
	blockNumber, err := strconv.Atoi(os.Getenv("BLOCK"))
	require.NoError(t, err)
	ctx := context.Background()
	extractor, err := New(ctx, config)
	require.NoError(t, err)
	api, err := sero.New(extractor.config.SeroNode)
	require.NoError(t, err)
	block, pkrs, err := extractor.extractOneBlock(ctx, api, int64(blockNumber))
	require.NoError(t, err)
	t.Logf("block = %#v, pkrs = %s", block, pkrs)
}
