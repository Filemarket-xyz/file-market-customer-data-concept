package repository

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/go-redis/redis/v8"
)

const (
	lastBlockKey = "last_block"
)

type blockCounter struct {
	rdb *redis.Client
}

func NewBlockCounter(rdb *redis.Client) BlockCounterRepo {
	return &blockCounter{
		rdb: rdb,
	}
}

func (b *blockCounter) GetLastBlock(ctx context.Context) (*big.Int, error) {
	num := b.rdb.Get(ctx, lastBlockKey)
	if num.Err() != nil {
		log.Printf("key error. Key: %s, err: %s", lastBlockKey, num.Err())
		return nil, num.Err()
	}
	var blockNum string
	if err := num.Scan(&blockNum); err != nil {
		return nil, err
	}
	res, ok := big.NewInt(0).SetString(blockNum, 10)
	if !ok {
		return nil, fmt.Errorf("parse block num: %s failed", blockNum)
	}
	return res, nil
}

func (b *blockCounter) SetLastBlock(ctx context.Context, lastBlock *big.Int) error {
	return b.rdb.Set(ctx, lastBlockKey, lastBlock.String(), redis.KeepTTL).Err()
}
