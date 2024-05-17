package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/repository"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/providers/dune"
	"github.com/ethereum/go-ethereum/common"
)

type DatasetsService struct {
	ctxService       context.Context
	repoUsers        repository.Users
	repoDatasets     repository.Datasets
	repoTransactions repository.Transactions

	duneProvider dune.Dune
	logging      logger.Logger
}

func NewDatasetsService(
	ctxService context.Context,
	repoUsers repository.Users,
	repoDatasets repository.Datasets,
	repoTransactions repository.Transactions,

	duneProvider dune.Dune,
	logging logger.Logger,
) Dataset {

	return &DatasetsService{
		ctxService:       ctxService,
		repoUsers:        repoUsers,
		repoDatasets:     repoDatasets,
		repoTransactions: repoTransactions,

		duneProvider: duneProvider,
		logging:      logging,
	}
}

func (s *DatasetsService) UploadDatasetsByAddress(clientId int64, address common.Address) {
	s.logging.Info("UploadDatasetsByAddress START client: ", clientId)

	wg := sync.WaitGroup{}
	syncCh := make(chan struct{})

	wg.Add(3)
	go s.getMultiChainBalances(clientId, address, &wg)
	go s.getRankLZUser(clientId, address, &wg)
	go s.getZkSyncEraUsersStatistics(clientId, address, &wg)

	s.logging.Info("Ожидание стягивания датасетов...")

	go func() {
		wg.Wait()
		syncCh <- struct{}{}
	}()
	select {
	case <-s.ctxService.Done():
		s.logging.Info("UploadDatasetsByAddress RETURN")
		return
	case <-syncCh:
	}

	s.logging.Info("Датасеты получены")

	if err := s.updateClientInfoDataset(clientId); err != nil {
		s.logging.Error("UploadDatasetsByAddress/updateClientInfoDataset: %w", err)
		return
	}

}

func (s *DatasetsService) updateClientInfoDataset(clientId int64) error {
	ctx, cancel := context.WithTimeout(s.ctxService, 10*time.Second)
	defer cancel()

	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("updateClientInfoDataset/BeginTransaction: %w", err)
	}
	defer tx.Rollback(ctx)

	client, err := s.repoUsers.GetClientById(ctx, tx, clientId)
	if err != nil {
		return fmt.Errorf("updateClientInfoDataset/GetClientById: %w", err)
	}
	client.Dataset = true
	if err := s.repoUsers.UpdateClient(ctx, tx, client); err != nil {
		return fmt.Errorf("updateClientInfoDataset/UpdateClient: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("updateClientInfoDataset/Commit: %w", err)

	}

	s.logging.Info("updateClientInfoDataset OK")

	return nil
}

/*
func (s *DatasetsService) getLayerZeroUserFees(clientId int64, address common.Address) {

		ctx, cancel := context.WithTimeout(s.ctxService, 60*10*time.Second)
		defer cancel()

		resp, err := s.duneProvider.LayerZeroUserFees(ctx, address.String())
		if err != nil {
			s.logging.Error("getLayerZeroUserFees/LayerZeroUserFees: %w", err)
			return
		}

		d := makeLayerZeroUserFeesDataset(resp)

		tx, err := s.repoTransactions.BeginTransaction(ctx)
		if err != nil {
			s.logging.Error("getLayerZeroUserFees/BeginTransaction: %w", err)
			return
		}
		defer tx.Rollback(ctx)

		if err := s.repoDatasets.AddLayerZeroUserFees(ctx, tx, clientId, d); err != nil {
			s.logging.Error("getLayerZeroUserFees/AddLayerZeroUserFees: %w", err)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			s.logging.Error("getLayerZeroUserFees/Commit: %w", err)
			return
		}

		s.logging.Info("getLayerZeroUserFees OK")
	}

	func makeLayerZeroUserFeesDataset(resp *dune.QueryResultResponse) *domain.LayerZeroUserFeesDataset {
		d := &domain.LayerZeroUserFeesDataset{}
		// норм решение для прототипа
		// for _, data := range resp.Result.Rows {
		// 	d.Rank = domain.P(data["Rank"].(int64))
		// 	d.CountTx = domain.P(data["# Transactions"].(int64))
		// 	d.TotalAmount = domain.P(data["total_amount"].(int64))
		// 	d.Arbitrum = domain.P(data["Arbitrum"].(float64))
		// 	d.Avalanche = domain.P(data["Avalanche"].(float64))
		// 	d.BSC = domain.P(data["BSC"].(float64))
		// 	d.Base = domain.P(data["Base"].(float64))
		// 	d.Celo = domain.P(data["Celo"].(float64))
		// 	d.Ethereum = domain.P(data["Ethereum"].(float64))
		// 	d.Fantom = domain.P(data["Fantom"].(float64))
		// 	d.Gnosis = domain.P(data["Gnosis"].(float64))
		// 	d.Linea = domain.P(data["Linea"].(float64))
		// 	d.Optimism = domain.P(data["Optimism"].(float64))
		// 	d.Polygon = domain.P(data["Polygon"].(float64))
		// 	d.Scroll = domain.P(data["Scroll"].(float64))
		// 	d.Zkevm = domain.P(data["Zkevm"].(float64))
		// 	d.Zksync = domain.P(data["Zksync"].(float64))
		// 	d.Zora = domain.P(data["Zora"].(float64))
		// }

		d.Rank = domain.P(int64(1))
		d.CountTx = domain.P(int64(1))
		d.TotalAmount = domain.P(int64(1))
		d.Arbitrum = domain.P(float64(2))
		d.Avalanche = domain.P(float64(2))
		d.BSC = domain.P(float64(2))
		d.Base = domain.P(float64(2))
		d.Celo = domain.P(float64(2))
		d.Ethereum = domain.P(float64(2))
		d.Fantom = domain.P(float64(2))
		d.Gnosis = domain.P(float64(2))
		d.Linea = domain.P(float64(2))
		d.Optimism = domain.P(float64(2))
		d.Polygon = domain.P(float64(2))
		d.Scroll = domain.P(float64(2))
		d.Zkevm = domain.P(float64(2))
		d.Zksync = domain.P(float64(2))
		d.Zora = domain.P(float64(2))

		return d
	}
*/

func (s *DatasetsService) getMultiChainBalances(clientId int64, address common.Address, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(s.ctxService, 60*20*time.Second)
	defer cancel()

	respResult, err := s.getExecuteResult(ctx, dune.MultiChainBalancesQueryId, address.String())
	if err != nil {
		s.logging.Error("getMultiChainBalances/getExecuteResult: %s", err.Error())
		return
	}

	d := makeMultiChainBalancesDataset(respResult)

	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		s.logging.Error("getMultiChainBalances/BeginTransaction: %w", err)
		return
	}
	defer tx.Rollback(ctx)

	if err := s.repoDatasets.AddMultiChainBalancesDataset(ctx, tx, clientId, d); err != nil {
		s.logging.Error("getMultiChainBalances/AddMultiChainBalancesDataset: %w", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		s.logging.Error("getMultiChainBalances/Commit: %w", err)
		return
	}

	s.logging.Info("getMultiChainBalances OK")
}

func makeMultiChainBalancesDataset(resp *dune.QueryResultResponse) *domain.MultiChainBalancesDataset {
	d := &domain.MultiChainBalancesDataset{}
	// норм решение для прототипа
	for _, data := range resp.Result.Rows {
		d.Chain = domain.P(strings.ReplaceAll(data["chain"].(string), " ", ""))
		d.Txs = domain.P(int64(data["TXs"].(float64)))
		d.Failed = domain.P(int64(data["Failed"].(float64)))
		d.Chain_Balance = domain.P(data["Chain_Balance"].(float64))
		d.Chain_Received = domain.P(data["Chain_Received"].(float64))
		d.Chain_Sent = domain.P(data["Chain_Sent"].(float64))
		d.Chain_TX_Fee = domain.P(data["Chain_TX_Fee"].(float64))
		d.USD_Balance = domain.P(data["USD_Balance"].(float64))
		d.USD_Received = domain.P(data["USD_Received"].(float64))
		d.USD_Sent = domain.P(data["USD_Sent"].(float64))
		d.USD_TX_Fee = domain.P(data["USD_TX_Fee"].(float64))
	}

	return d
}

func (s *DatasetsService) getRankLZUser(clientId int64, address common.Address, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(s.ctxService, 60*10*time.Second)
	defer cancel()

	respResult, err := s.getExecuteResult(ctx, dune.RankUserQueryId, address.String())
	if err != nil {
		s.logging.Error("getRankLZUser/getExecuteResult: %s", err.Error())
		return
	}

	d := makeRankLayerZeroUserDataset(respResult)

	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		s.logging.Error("getRankLZUser/BeginTransaction: %w", err)
		return
	}
	defer tx.Rollback(ctx)

	if err := s.repoDatasets.AddRankLayerZeroUser(ctx, tx, clientId, d); err != nil {
		s.logging.Error("getRankLZUser/AddRankLayerZeroUser: %w", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		s.logging.Error("getRankLZUser/Commit: %w", err)
		return
	}

	s.logging.Info("getRankLZUser OK")
}

func makeRankLayerZeroUserDataset(resp *dune.QueryResultResponse) *domain.RankLZDataset {
	d := &domain.RankLZDataset{}
	// норм решение для прототипа
	for _, data := range resp.Result.Rows {
		d.Rk = domain.P(int64(data["rk"].(float64)))
		d.Rs = domain.P(int64(data["rs"].(float64)))
		d.Tc = domain.P(int64(data["tc"].(float64)))
		d.Amt = domain.P(data["amt"].(float64))
		d.Cc = domain.P(strings.ReplaceAll(data["cc"].(string), " ", ""))
		d.Dwm = domain.P(strings.ReplaceAll(data["dwm"].(string), " ", ""))
		d.Ibt = domain.P(strings.ReplaceAll(data["ibt"].(string), " ", ""))
		d.Lzd = domain.P(int64(data["lzd"].(float64)))
	}

	return d
}

func (s *DatasetsService) getZkSyncEraUsersStatistics(clientId int64, address common.Address, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(s.ctxService, 60*10*time.Second)
	defer cancel()

	respResult, err := s.getExecuteResult(ctx, dune.ZkSyncUsersStatQueryId, address.String())
	if err != nil {
		s.logging.Error("getZkSyncUsersStatistics/getExecuteResult: %s", err.Error())
		return
	}

	d := makeZkSyncEraUsersStatDataset(respResult)

	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		s.logging.Error("getZkSyncUsersStatistics/BeginTransaction: %s", err.Error())
		return
	}
	defer tx.Rollback(ctx)

	if err := s.repoDatasets.AddZkSyncEraUserStat(ctx, tx, clientId, d); err != nil {
		s.logging.Error("getZkSyncUsersStatistics/AddZkSyncEraUserStat: %s", err.Error())
		return
	}

	if err := tx.Commit(ctx); err != nil {
		s.logging.Error("getZkSyncUsersStatistics/Commit: %w", err)
		return
	}

	s.logging.Info("getZkSyncEraUsersStatistics OK")
}

func makeZkSyncEraUsersStatDataset(resp *dune.QueryResultResponse) *domain.ZkSyncEraUserStatDataset {
	d := &domain.ZkSyncEraUserStatDataset{}
	// норм решение для прототипа
	for _, data := range resp.Result.Rows {
		d.Rk = domain.P(int64(data["rk"].(float64)))
		d.Amt = domain.P(strings.ReplaceAll(data["amt"].(string), " ", ""))
		d.Cc = domain.P(strings.ReplaceAll(data["cc"].(string), " ", ""))
		d.Dwm = domain.P(strings.ReplaceAll(data["dwm"].(string), " ", ""))
		d.Ibt = domain.P(strings.ReplaceAll(data["ibt"].(string), " ", ""))
		d.Tc = domain.P(strings.ReplaceAll(data["tc"].(string), " ", ""))
	}

	return d
}

func (s *DatasetsService) getExecuteResult(ctx context.Context, queryId, address string) (*dune.QueryResultResponse, error) {
	resp, err := s.duneProvider.ExecuteQuery(ctx, queryId, &dune.Req{
		RequsetData: &dune.RequsetData{
			Addr: address,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("getExecuteResult/ExecuteQuery: %w", err)

	}

	executionId := resp.ExecutionId
	respResult, err := s.getExecutionResult(ctx, executionId)
	if err != nil {
		return nil, fmt.Errorf("getExecuteResult/getExecutionResult: %w", err)
	}

	return respResult, nil
}

func (s *DatasetsService) getExecutionResult(ctx context.Context, executionId string) (*dune.QueryResultResponse, error) {
	var (
		respResult *dune.QueryResultResponse
		err        error
	)
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("getExecutionResult: ctx deadline")
		case <-time.After(5 * time.Second):
		}
		respResult, err = s.duneProvider.GetExecutionResult(ctx, executionId, "")
		if err != nil {
			return nil, fmt.Errorf("getExecutionResult/GetExecutionResult: %w", err)

		}
		if respResult.State == dune.COMPLETEDStateResult {
			break
		}
		if respResult.State == dune.FAILEDStateResult {
			return nil, fmt.Errorf("getExecutionResult: FAIL")
		}
	}
	return respResult, nil
}
