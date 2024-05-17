package service

import (
	"context"
	"fmt"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/repository"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/time_manager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type ClientService struct {
	cfg              *config.ServiceConfig
	repoUsers        repository.Users
	repoDatasets     repository.Datasets
	repoTransactions repository.Transactions
	timeManager      time_manager.TimeManager
	logging          logger.Logger
}

func NewClientService(
	cfg *config.ServiceConfig,
	repoUsers repository.Users,
	repoDatasets repository.Datasets,
	repoTransactions repository.Transactions,
	timeManager time_manager.TimeManager,
	logging logger.Logger,
) Client {

	return &ClientService{
		cfg:              cfg,
		repoUsers:        repoUsers,
		repoDatasets:     repoDatasets,
		repoTransactions: repoTransactions,
		timeManager:      timeManager,
		logging:          logging,
	}
}

func (s *ClientService) UpdateClientAgreement(ctx context.Context, clientId int64, agreement bool) (*models.Client, error) {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("UpdateClientAgreement/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(ctx)

	client, err := s.repoUsers.GetClientById(ctx, tx, clientId)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("UpdateClientAgreement/GetClientById: %w", err), InternalError, "")
	}

	if client.Agreement {
		return nil, newServiceError(code400,
			fmt.Errorf("UpdateClientAgreement: there is already agreement"), "there is already agreement", "")
	}

	client.Agreement = agreement
	if agreement {
		client.PointBalance = client.PointBalance.Add(decimal.NewFromInt(500)) // TODO определить, сколько насыпать за согласие
	}
	if !agreement {
		client.Dataset = false
		if err := s.repoDatasets.DeleteByClientId(ctx, tx, client.Id); err != nil {
			return nil, newServiceError(code500,
				fmt.Errorf("UpdateClientAgreement/DeleteByClientId: %w", err), InternalError, "")
		}
	}
	if err := s.repoUsers.UpdateClient(ctx, tx, client); err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("UpdateClientAgreement/UpdateClient: %w", err), InternalError, "")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("UpdateClientAgreement/Commit: %w", err), InternalError, "")
	}

	return client.ToModel(), nil
}

func (s *ClientService) GetUserDataset(ctx context.Context, user *domain.UserWithTokenNumber) ([][]string, error) {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetUserDataset/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(ctx)

	if user.Role != domain.RoleClient {
		return nil, newServiceError(code400,
			fmt.Errorf("GetUserDataset: %s", UserNotClient), UserNotClient, "")
	}
	client, err := s.repoUsers.GetClientById(ctx, tx, user.Id)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetUserDataset/GetClientById: %w", err), InternalError, "")
	}
	// if !client.Agreement {
	// }
	if !client.Bought {
		return nil, newServiceError(code400,
			fmt.Errorf("GetUserDataset: client did not pay"), "The client did not pay for the data set", "")
	}

	dataset, err := s.repoDatasets.GetDatasetsByClientId(ctx, tx, user.Id)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetUserDataset/GetClientById: %w", err), InternalError, "")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetUserDataset/Commit: %w", err), InternalError, "")
	}

	res, err := domain.SliceDataToStrings([]*domain.Dataset{dataset})
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetUserDataset/DatasetsToStrings: %w", err), InternalError, "")
	}

	return res, nil
}

func (s *ClientService) GetDataset(ctx context.Context, clientId int64) (*models.GetClientDatasetResponse, error) {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetDataset/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(ctx)

	client, err := s.repoUsers.GetClientById(ctx, tx, clientId)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetDataset/GetClientById: %w", err), InternalError, "")
	}
	if !client.Dataset {
		return nil, newServiceError(code400,
			fmt.Errorf("GetDataset: the client does not have a dataset"), "the client does not have a dataset", "")
	}

	dataset, err := s.repoDatasets.GetDatasetsByClientId(ctx, tx, client.Id)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetDataset/GetClientById: %w", err), InternalError, "")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetDataset/Commit: %w", err), InternalError, "")
	}

	res := &models.GetClientDatasetResponse{
		Dataset: dataset.ToModel(),
	}

	return res, nil
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
