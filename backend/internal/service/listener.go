package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/repository"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/ethclient"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type ListenerService struct {
	cfg                 *config.ServiceConfig
	repoBlockCounter    repository.BlockCounterRepo
	repoTransactions    repository.Transactions
	repoEthTransactions repository.EthTransactions

	ethClient ethclient.Client

	logging logger.Logger

	stopCh chan struct{}
}

func NewListenerService(
	cfg *config.ServiceConfig,
	repoBlockCounter repository.BlockCounterRepo,
	repoTransactions repository.Transactions,
	repoEthTransactions repository.EthTransactions,

	ethClient ethclient.Client,

	logging logger.Logger,

	stopCh chan struct{},
) Listener {
	return &ListenerService{
		cfg:                 cfg,
		repoBlockCounter:    repoBlockCounter,
		repoTransactions:    repoTransactions,
		repoEthTransactions: repoEthTransactions,

		ethClient: ethClient,

		logging: logging,

		stopCh: stopCh,
	}
}

func (s *ListenerService) listenBlockchain() error {
	lastBlock, err := s.repoBlockCounter.GetLastBlock(context.Background())
	if err != nil {
		if err == redis.Nil {
			latestBlock, err := s.ethClient.GetLatestBlockNumber(context.Background())
			if err != nil {
				log.Panicln(err)
				return err
			}
			lastBlock = latestBlock
		} else {
			return err
		}
	}
	go func() {
		for {
			select {
			case <-time.After(s.cfg.Periods.ListenerPeriod):
				current, err := s.checkBlock(lastBlock)
				if err != nil {
					s.logging.Errorf("process block failed: ", err)
				}
				lastBlock = current
			case <-s.stopCh:
				return
			}
		}
	}()
	return nil
}

func (s *ListenerService) checkBlock(latest *big.Int) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 13*time.Second)
	defer cancel()
	blockNum, err := s.ethClient.GetLatestBlockNumber(ctx)
	if err != nil {
		s.logging.Error("get latest block failed: ", err)
		return latest, err
	}
	if blockNum.Cmp(latest) != 0 {
		s.logging.Info("processing block difference: ", latest.String(), " ", blockNum.String())
	}
	for blockNum.Cmp(latest) != 0 {
		latest, err = s.checkSingleBlock(latest)
		if err != nil {
			return latest, err
		}
	}
	return latest, nil
}

func (s *ListenerService) checkSingleBlock(latest *big.Int) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pending := big.NewInt(0).Add(latest, big.NewInt(1))
	block, err := s.ethClient.CustomBlockByNumber(ctx, hexutil.EncodeBig(pending), true)
	if err != nil {
		s.logging.Error("get pending block failed ", pending.String(), " ", err)
		return latest, err
	} else {
		if err := s.processBlock(block); err != nil {
			s.logging.Error("process block failed ", err)
			return latest, err
		}
	}
	if err := s.repoBlockCounter.SetLastBlock(context.Background(), pending); err != nil {
		s.logging.Error("set last block failed ", err)
	}
	return pending, err
}

func (s *ListenerService) processBlock(block *types.Block) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	logs, err := s.ethClient.GetLogs(ctx, block.Number(), block.Number())
	if err != nil {
		return fmt.Errorf("process block get logs failed: %w", err)
	}

	for _, t := range block.Transactions() {
		if t.To() == nil {
			continue
		}
		txLogs := logsForTransaction(t.Hash(), logs)
		to := *t.To()

		if s.cfg.Wallet.Cmp(to) == 0 {
			if err := s.processTransaction(
				ctx,
				&domain.Transaction{
					Id:        strings.ToLower(t.Hash().String()),
					State:     domain.TransactionStateConfirmed,
					Timestamp: time.Now().UnixMilli(),
				},
				txLogs,
			); err != nil {
				s.logging.Error("process block tx failed", zap.String("id", t.Hash().String()), zap.Error(err))
			}
		}
	}

	return nil
}

func (s *ListenerService) processTransaction(
	ctx context.Context,
	ethTx *domain.Transaction,
	logs []*types.Log,
) error {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("processTransaction begin transaction failed: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = s.repoEthTransactions.GetTransaction(ctx, tx, ethTx.Id)
	if err != nil && !errors.Is(err, repository.ErrNoRows) {
		s.logging.Error("processTransaction/GetTransaction: ", err)
		return fmt.Errorf("processTransaction/GetTransaction: %w", err)
	}
	if err == nil {
		return nil
	}

	for i, l := range logs {
		var err error

		// TODO

		if err != nil {
			s.logging.Error("processTransaction failed", zap.String("id", ethTx.Id), zap.Int("log_index", i), zap.String("contract address", l.Address.String()), zap.Error(err))
			return fmt.Errorf("processTransaction/contract address(%s): %w", l.Address.String(), err)
		}
	}

	if err := s.repoEthTransactions.InsertTransaction(ctx, tx, &domain.Transaction{
		Id:        ethTx.Id,
		State:     ethTx.State,
		Timestamp: ethTx.Timestamp,
	}); err != nil {
		return fmt.Errorf("processTransaction/InsertTransaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("processTransaction commit failed: %w", err)
	}

	return nil
}
