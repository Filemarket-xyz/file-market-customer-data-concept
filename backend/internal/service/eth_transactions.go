package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/ethclient"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

type EthTransactionsService struct {
	ethClient ethclient.Client
	logging   logger.Logger
}

func NewEthTransactionsService(
	ethClient ethclient.Client,
	logging logger.Logger,
) EthTransactions {
	return &EthTransactionsService{
		ethClient: ethClient,
		logging:   logging,
	}
}

func (t *EthTransactionsService) getTxLogs(ctx context.Context, hash string) ([]*types.Log, error) {
	for i := 0; i < 30; i++ {
		time.Sleep(200 * time.Millisecond)
		rec, err := t.ethClient.TransactionReceipt(ctx, common.HexToHash(hash))
		if err != nil {
			t.logging.Error("get tx logs failed", zap.String("hash", hash), zap.Error(err))
			continue
		}
		if rec == nil {
			continue
		}
		return rec.Logs, nil
	}
	return nil, fmt.Errorf("failed to fetch tx logs")
}

func logsForTransaction(hash common.Hash, logs []*types.Log) []*types.Log {
	var res []*types.Log
	for _, l := range logs {
		if l.TxHash == hash {
			res = append(res, l)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Index < res[j].Index
	})
	return res
}
