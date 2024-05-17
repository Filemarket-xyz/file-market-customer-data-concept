package service

import (
	"context"
	"time"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/repository"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/ethclient"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/hash"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/jwtoken"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/providers/dune"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/rand_manager"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/time_manager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
)

type Auth interface {
	GetUserByJWToken(ctx context.Context, purpose jwtoken.Purpose, token string) (*domain.UserWithTokenNumber, error)
	RefreshJWTokens(ctx context.Context, user *domain.UserWithTokenNumber) (*models.AuthResponse, *domain.PairOfTokens, error)
	Logout(ctx context.Context, user *domain.UserWithTokenNumber) error
	FullLogout(ctx context.Context, user *domain.UserWithTokenNumber) error
	GetAuthMessage(ctx context.Context, req *models.AuthMessageRequest) (*models.AuthMessageResponse, error)
	AuthByMessage(ctx context.Context, req *models.AuthBySignatureRequest) (*models.AuthResponse, *domain.PairOfTokens, error)
	AuthByToken(ctx context.Context, user *domain.UserWithTokenNumber) (*models.AuthResponse, error)
	DropExpiredTokens(ctx context.Context) error
}

type Client interface {
	UpdateClientAgreement(ctx context.Context, clientId int64, agreement bool) (*models.Client, error)
	GetUserDataset(ctx context.Context, user *domain.UserWithTokenNumber) ([][]string, error)
	GetDataset(ctx context.Context, clientId int64) (*models.GetClientDatasetResponse, error)
	FixingPurchaseDataSet(ctx context.Context, tx repository.Transaction, from common.Address, value decimal.Decimal) error
}

type Dataset interface {
	UploadDatasetsByAddress(clientId int64, address common.Address)
}

type Config interface {
	GetConfig(ctx context.Context) *models.Config
}

type EthTransactions interface {
	getTxLogs(ctx context.Context, hash string) ([]*types.Log, error)
}

type Listener interface {
	listenBlockchain(ctx context.Context) error
}

type Service interface {
	Auth
	Client
	Config
	Dataset
	Listener
	EthTransactions

	Shutdown()
}

type service struct {
	Auth
	Client
	Config
	Dataset
	Listener
	EthTransactions

	ethClient ethclient.Client

	ctxBackground context.Context
	cancelCtx     context.CancelFunc
	stopCh        chan struct{}

	cfg     *config.ServiceConfig
	logging logger.Logger
}

func NewService(
	ctx context.Context,
	repo *repository.Repository,
	ethClient ethclient.Client,
	jwtTokenManager jwtoken.JWTokenManager,
	hashManager hash.HashManager,
	timeManager time_manager.TimeManager,
	randManager rand_manager.RandManager,

	duneProvider dune.Dune,
	cfg *config.ServiceConfig,
	logging logger.Logger,
) (Service, error) {

	var (
		ctxBackground, cancelCtx = context.WithCancel(context.Background())

		stopCh = make(chan struct{})

		Dataset = NewDatasetsService(ctxBackground, repo.Users, repo.Datasets, repo.Transactions, duneProvider, logging)
		Auth    = NewAuthService(cfg, repo.Users, repo.JWTokens, repo.Transactions, Dataset,
			jwtTokenManager, hashManager, timeManager, randManager, logging)
		Config          = NewConfigService(cfg, logging)
		Client          = NewClientService(cfg, repo.Users, repo.Datasets, repo.Transactions, timeManager, logging)
		EthTransactions = NewEthTransactionsService(ethClient, logging)

		listener = NewListenerService(cfg, repo.BlockCounterRepo, repo.Transactions, repo.EthTransactions, Client,
			ethClient, logging, stopCh)
	)

	res := &service{
		Auth:            Auth,
		Client:          Client,
		Config:          Config,
		Dataset:         Dataset,
		Listener:        listener,
		EthTransactions: EthTransactions,

		ethClient: ethClient,

		cfg:           cfg,
		logging:       logging,
		stopCh:        stopCh,
		ctxBackground: ctxBackground,
		cancelCtx:     cancelCtx,
	}

	if err := listener.listenBlockchain(res.ctxBackground); err != nil {
		return nil, err
	}

	go res.dropExpiredTokens()

	return res, nil
}

func (s *service) dropExpiredTokens() {
	for {
		after := time.After(s.cfg.Periods.JwtCheckPeriod)
		select {
		case <-after:
			ctx, cancel := context.WithTimeout(s.ctxBackground, 30*time.Second)
			defer cancel()

			if err := s.Auth.DropExpiredTokens(ctx); err != nil {
				s.logging.Error("check payment status error:", err)
			}

		case <-s.stopCh:
			return
		}
	}
}

func (s *service) Shutdown() {
	s.cancelCtx()

	for i := 0; i < 1; i++ {
		s.stopCh <- struct{}{}
	}
	close(s.stopCh)
}
