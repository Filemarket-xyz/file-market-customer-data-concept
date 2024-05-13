package ethclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Client interface {
	bind.ContractBackend
	ethereum.ChainReader
	ethereum.ChainStateReader
	ethereum.TransactionReader
	GetLogs(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) ([]*types.Log, error)
	GetLatestBlockNumber(ctx context.Context) (*big.Int, error)
	CustomBlockByNumber(ctx context.Context, args ...interface{}) (*types.Block, error)
}

func HexToAddress(s string) common.Address {
	return common.HexToAddress(s)
}
