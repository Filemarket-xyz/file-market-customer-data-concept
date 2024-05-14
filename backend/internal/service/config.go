package service

import (
	"context"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
)

type ConfigService struct {
	cfg *config.ServiceConfig

	logging logger.Logger
}

func NewConfigService(
	cfg *config.ServiceConfig,

	logging logger.Logger,
) Config {

	return &ConfigService{
		cfg: cfg,

		logging: logging,
	}
}

func (s *ConfigService) GetConfig(ctx context.Context) *models.Config {

	return &models.Config{
		WalletAddress: domain.P(s.cfg.Wallet.String()),
		DatasetCost:   domain.P(s.cfg.DatasetCost.String()),
	}
}
