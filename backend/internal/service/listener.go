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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type ListenerService struct {
	cfg                 *config.ServiceConfig
	repoBlockCounter    repository.BlockCounterRepo
	repoTransactions    repository.Transactions
	repoEthTransactions repository.EthTransactions

	Client Client

	ethClient ethclient.Client

	logging logger.Logger

	stopCh chan struct{}
}

func NewListenerService(
	cfg *config.ServiceConfig,
	repoBlockCounter repository.BlockCounterRepo,
	repoTransactions repository.Transactions,
	repoEthTransactions repository.EthTransactions,

	Client Client,

	ethClient ethclient.Client,

	logging logger.Logger,

	stopCh chan struct{},
) Listener {
	return &ListenerService{
		cfg:                 cfg,
		repoBlockCounter:    repoBlockCounter,
		repoTransactions:    repoTransactions,
		repoEthTransactions: repoEthTransactions,

		Client: Client,

		ethClient: ethClient,

		logging: logging,

		stopCh: stopCh,
	}
}

func (s *ListenerService) listenBlockchain(ctx context.Context) error {
	lastBlock, err := s.repoBlockCounter.GetLastBlock(ctx)
	if err != nil {
		if err == redis.Nil {
			latestBlock, err := s.ethClient.GetLatestBlockNumber(ctx)
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
				current, err := s.checkBlock(ctx, lastBlock)
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

func (s *ListenerService) checkBlock(ctx context.Context, latest *big.Int) (*big.Int, error) {
	ctx2, cancel := context.WithTimeout(ctx, 13*time.Second)
	defer cancel()
	blockNum, err := s.ethClient.GetLatestBlockNumber(ctx2)
	if err != nil {
		s.logging.Error("get latest block failed: ", err)
		return latest, err
	}
	if blockNum.Cmp(latest) != 0 {
		s.logging.Info("processing block difference: ", latest.String(), " ", blockNum.String())
	}
	for blockNum.Cmp(latest) != 0 {
		latest, err = s.checkSingleBlock(ctx, latest)
		if err != nil {
			return latest, err
		}
	}
	return latest, nil
}

func (s *ListenerService) checkSingleBlock(ctx context.Context, latest *big.Int) (*big.Int, error) {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	pending := big.NewInt(0).Add(latest, big.NewInt(1))
	block, err := s.ethClient.BlockByNumber(ctx2, pending)
	if err != nil {
		s.logging.Error("get pending block failed ", pending.String(), " ", err)
		return latest, err
	} else {
		if err := s.processBlock(ctx, block); err != nil {
			s.logging.Error("process block failed ", err)
			return latest, err
		}
	}
	if err := s.repoBlockCounter.SetLastBlock(ctx, pending); err != nil {
		s.logging.Error("set last block failed ", err)
	}
	return pending, err
}

func (s *ListenerService) processBlock(ctx context.Context, block *types.Block) error {
	ctx2, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()

	for _, t := range block.Transactions() {
		if t.To() == nil {
			continue
		}
		to := *t.To()

		if s.cfg.Wallet.Cmp(to) == 0 {
			from, err := types.Sender(types.LatestSignerForChainID(t.ChainId()), t)
			if err != nil {
				return fmt.Errorf("get sender error: %w", err)
			}
			if err := s.processTransaction(
				ctx2,
				&domain.Transaction{
					Id:        strings.ToLower(t.Hash().String()),
					State:     domain.TransactionStateConfirmed,
					Timestamp: time.Now().UnixMilli(),
				},
				from,
				t.Value(),
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
	from common.Address,
	value *big.Int,
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

	if err := s.Client.FixingPurchaseDataSet(ctx, tx, from, decimal.NewFromBigInt(value, 0)); err != nil {
		s.logging.Error("processTransaction/FixingPurchaseDataSet: ", err)
		return fmt.Errorf("processTransaction/FixingPurchaseDataSet: %w", err)
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
