package service

import (
	"context"
	"fmt"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/repository"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/hash"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/jwtoken"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/time_manager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type ClientService struct {
	cfg              *config.ServiceConfig
	repoUsers        repository.Users
	repoJWTokens     repository.JWTokens
	repoTransactions repository.Transactions
	jwtManager       jwtoken.JWTokenManager
	hashManager      hash.HashManager
	timeManager      time_manager.TimeManager
	logging          logger.Logger
}

func NewClientService(
	cfg *config.ServiceConfig,
	repoUsers repository.Users,
	repoTransactions repository.Transactions,
	timeManager time_manager.TimeManager,
	logging logger.Logger,
) Client {

	return &ClientService{
		cfg:              cfg,
		repoUsers:        repoUsers,
		repoTransactions: repoTransactions,
		timeManager:      timeManager,
		logging:          logging,
	}
}

func (s *ClientService) FixingPurchaseDataSet(ctx context.Context, tx repository.Transaction, from common.Address, value decimal.Decimal) error {
	userId, userRole, err := s.repoUsers.GetUserIdByAddress(ctx, tx, from.String())
	if err != nil {
		return fmt.Errorf("FixingPurchaseDataSet/GetUserIdByAddress: %w", err)
	}
	if userId == 0 {
		s.logging.Info("FixingPurchaseDataSet: client not exist!!!")
		return nil
	}
	if userRole != domain.RoleClient {
		return fmt.Errorf("FixingPurchaseDataSet: user not client")
	}
	if value.Cmp(s.cfg.DatasetCost) == -1 {
		s.logging.Info("FixingPurchaseDataSet: not enough Ethereum to buy!!!!")
		return nil
	}

	if err := s.repoUsers.UpdateDatasetPurchaseData(ctx, tx, true, userId); err != nil {
		return fmt.Errorf("FixingPurchaseDataSet/UpdateDatasetPurchaseData: %w", err)
	}

	return nil
}
